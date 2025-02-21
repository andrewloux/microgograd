package plot

import (
	"fmt"
	"os"
	"strings"

	"microgograd/micrograd"

	"github.com/emicklei/dot"
)

// NodeLabelFunc is a function type that generates a label for a Value node
type NodeLabelFunc[K micrograd.BaseNumeric] func(v *micrograd.Value[K]) string

// PlotOption represents an option for configuring the plot
type PlotOption[K micrograd.BaseNumeric] func(*plotConfig[K])

type plotConfig[K micrograd.BaseNumeric] struct {
	labelFunc NodeLabelFunc[K]
}

// WithNodeLabelFunc returns a PlotOption that sets a custom node labeling function
func WithNodeLabelFunc[K micrograd.BaseNumeric](fn NodeLabelFunc[K]) PlotOption[K] {
	return func(cfg *plotConfig[K]) {
		cfg.labelFunc = fn
	}
}

// defaultNodeLabel is the default implementation for node labeling
func defaultNodeLabel[K micrograd.BaseNumeric](v *micrograd.Value[K]) string {
	var parts []string

	// Add name if present
	if v.Name != "" {
		parts = append(parts, v.Name)
	}

	// Add value
	parts = append(parts, fmt.Sprintf("%v", v.GetValue()))

	return fmt.Sprintf("{%s}", strings.Join(parts, " | "))
}

// nodeID generates a unique identifier for a Value node
func nodeID[K micrograd.BaseNumeric](v *micrograd.Value[K]) string {
	return fmt.Sprintf("%p", v)
}

// collectNodes collects all nodes in topological order
func collectNodes[K micrograd.BaseNumeric](v *micrograd.Value[K], visited map[string]bool, nodes *[]*micrograd.Value[K]) {
	id := nodeID(v)
	if visited[id] {
		return
	}
	visited[id] = true

	// Process children first
	children := v.GetChildren()
	for _, child := range children.Data() {
		if child == nil {
			continue
		}
		if cv, ok := child.(*micrograd.Value[K]); ok {
			collectNodes(cv, visited, nodes)
		}
	}

	// Add current node after all children
	*nodes = append(*nodes, v)
}

// dotFromValue generates a DOT representation of the computation graph
func dotFromValue[K micrograd.BaseNumeric](v *micrograd.Value[K], cfg *plotConfig[K]) string {
	// First collect nodes in topological order
	visited := make(map[string]bool)
	var nodes []*micrograd.Value[K]
	collectNodes(v, visited, &nodes)

	// Create graph
	g := dot.NewGraph(dot.Directed)

	// Set graph-level attributes
	g.Attr("rankdir", "LR")     // Left to right direction
	g.Attr("nodesep", "0.8")    // More space between nodes
	g.Attr("ranksep", "0.6")    // More space between ranks
	g.Attr("splines", "curved") // Use curved edges for better aesthetics

	// Create nodes in topological order
	nodeMap := make(map[string]dot.Node)
	for _, node := range nodes {
		nodeMap[nodeID(node)] = g.Node(nodeID(node)).
			Attr("label", cfg.labelFunc(node)).
			Attr("shape", "record").
			Attr("style", "filled").
			Attr("fillcolor", "white").
			Attr("class", "node value-node")
	}

	// Add operation nodes and edges
	for _, node := range nodes {
		if op := node.GetOperation(); op != micrograd.UNSET {
			// Create operation node
			opNodeID := fmt.Sprintf("%s_op", nodeID(node))
			opNode := g.Node(opNodeID).
				Attr("label", string(rune(op))).
				Attr("shape", "circle").
				Attr("style", "filled").
				Attr("fillcolor", "#f0f0f0").
				Attr("width", "0.5").
				Attr("height", "0.5").
				Attr("fontsize", "20").
				Attr("penwidth", "2").
				Attr("fixedsize", "true").
				Attr("margin", "0").
				Attr("pad", "0.2").
				Attr("class", "node op-node")

			// Connect operation node to parent
			g.Edge(opNode, nodeMap[nodeID(node)]).
				Attr("color", "#666666").
				Attr("penwidth", "1.5").
				Attr("class", "edge").
				Attr("id", fmt.Sprintf("edge_%s_%s", opNodeID, nodeID(node))).
				Attr("data-source", opNodeID).
				Attr("data-target", nodeID(node))

			// Connect children to operation node
			children := node.GetChildren()
			for _, child := range children.Data() {
				if child == nil {
					continue
				}
				if cv, ok := child.(*micrograd.Value[K]); ok {
					g.Edge(nodeMap[nodeID(cv)], opNode).
						Attr("color", "#666666").
						Attr("penwidth", "1.5").
						Attr("class", "edge").
						Attr("id", fmt.Sprintf("edge_%s_%s", nodeID(cv), opNodeID)).
						Attr("data-source", nodeID(cv)).
						Attr("data-target", opNodeID)
				}
			}
		}
	}

	return g.String()
}

// WriteInteractiveHTML generates an HTML file with draggable SVG graph
func WriteInteractiveHTML[K micrograd.BaseNumeric](v *micrograd.Value[K], filePath string, opts ...PlotOption[K]) error {
	// Set up configuration with defaults
	cfg := &plotConfig[K]{
		labelFunc: defaultNodeLabel[K],
	}

	// Apply options
	for _, opt := range opts {
		opt(cfg)
	}

	// Generate DOT
	dot := dotFromValue(v, cfg)

	// Read HTML template
	templateContent, err := os.ReadFile("plot/template.html")
	if err != nil {
		return fmt.Errorf("failed to read template file: %v", err)
	}

	// Replace placeholder with DOT content
	html := strings.Replace(string(templateContent), "<!-- SVG_CONTENT -->", dot, 1)

	return os.WriteFile(filePath, []byte(html), 0644)
}

// WriteGraph generates a DOT representation of the computation graph and writes it to a file
func WriteGraph[K micrograd.BaseNumeric](v *micrograd.Value[K], filePath string, opts ...PlotOption[K]) error {
	// Set up configuration with defaults
	cfg := &plotConfig[K]{
		labelFunc: defaultNodeLabel[K],
	}

	// Apply options
	for _, opt := range opts {
		opt(cfg)
	}

	// Generate DOT
	dot := dotFromValue(v, cfg)

	// Write to file
	return os.WriteFile(filePath, []byte(dot), 0644)
}
