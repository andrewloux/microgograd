package plot

import (
	"embed"
	"fmt"
	"os"
	"strings"

	"microgograd/micrograd"

	"github.com/emicklei/dot"
)

//go:embed template.html
var templateFS embed.FS

// NodeLabelFunc is a function type that generates a label for a Value node
type NodeLabelFunc[K micrograd.BaseNumeric] func(node micrograd.Numeric[K]) string

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
func defaultNodeLabel[K micrograd.BaseNumeric](node micrograd.Numeric[K]) string {
	var parts []string

	name := node.GetName()
	if name == "" {
		name = "?"
	}

	parts = append(parts, fmt.Sprintf("Name: %s", name))
	parts = append(parts, fmt.Sprintf("Value: %.4f", node.GetValue()))
	parts = append(parts, fmt.Sprintf("Grad: %.4f", node.GetGradient()))

	return strings.Join(parts, "\\n")
}

// nodeID generates a unique identifier for a Value node
func nodeID[K micrograd.BaseNumeric](node micrograd.Numeric[K]) string {
	return fmt.Sprintf("%p", node)
}

func pairToSlice[K micrograd.BaseNumeric](p micrograd.Pair[micrograd.Numeric[K]]) []micrograd.Numeric[K] {
	return []micrograd.Numeric[K]{p.X(), p.Y()}
}

// collectNodes collects all nodes in topological order
func collectNodes[K micrograd.BaseNumeric](node micrograd.Numeric[K], visited map[string]bool, nodes *[]micrograd.Numeric[K]) {
	id := nodeID(node)
	if visited[id] {
		return
	}
	visited[id] = true

	for _, child := range pairToSlice(node.GetChildren()) {
		if child != nil {
			collectNodes(child, visited, nodes)
		}
	}

	*nodes = append(*nodes, node)
}

// dotFromValue generates a DOT representation of the computation graph
func dotFromValue[K micrograd.BaseNumeric](node micrograd.Numeric[K], cfg *plotConfig[K]) string {
	// First collect nodes in topological order
	visited := make(map[string]bool)
	var allNodes []micrograd.Numeric[K]
	collectNodes(node, visited, &allNodes)

	// Create graph
	g := dot.NewGraph(dot.Directed)

	// Set graph-level attributes
	g.Attr("rankdir", "LR")     // Left to right direction
	g.Attr("nodesep", "2.0")    // More space between nodes
	g.Attr("ranksep", "2.0")    // More space between ranks
	g.Attr("margin", "1.0")     // Graph margin
	g.Attr("splines", "curved") // Use curved edges for better aesthetics

	// Create nodes in topological order
	nodeMap := make(map[string]dot.Node)
	for _, n := range allNodes {
		nodeMap[nodeID(n)] = g.Node(nodeID(n)).
			Attr("label", cfg.labelFunc(n)).
			Attr("shape", "record").
			Attr("style", "filled").
			Attr("fillcolor", "white").
			Attr("class", "node value-node")
	}

	// Add operation nodes and edges
	for _, n := range allNodes {
		if op := n.GetOperation(); op != micrograd.UNSET {
			// Create operation node
			opNodeID := fmt.Sprintf("%s_op", nodeID(n))
			opNode := g.Node(opNodeID).
				Attr("label", string(rune(op))).
				Attr("shape", "ellipse").
				Attr("style", "filled").
				Attr("fillcolor", "#f0f0f0").
				Attr("width", "0.8").
				Attr("height", "0.8").
				Attr("fontsize", "24").
				Attr("penwidth", "2").
				Attr("fixedsize", "true").
				Attr("margin", "0.2").
				Attr("pad", "0.3").
				Attr("class", "node op-node")

			// Connect operation node to parent
			g.Edge(opNode, nodeMap[nodeID(n)]).
				Attr("color", "#666666").
				Attr("penwidth", "1.5").
				Attr("class", "edge").
				Attr("id", fmt.Sprintf("edge_%s_%s", opNodeID, nodeID(n))).
				Attr("data-source", opNodeID).
				Attr("data-target", nodeID(n))

			// Connect children to operation node
			for _, child := range pairToSlice(n.GetChildren()) {
				if child != nil {
					g.Edge(nodeMap[nodeID(child)], opNode).
						Attr("color", "#666666").
						Attr("penwidth", "1.5").
						Attr("class", "edge").
						Attr("id", fmt.Sprintf("edge_%s_%s", nodeID(child), opNodeID)).
						Attr("data-source", nodeID(child)).
						Attr("data-target", opNodeID)
				}
			}
		}
	}

	return g.String()
}

// WriteInteractiveHTML generates an HTML file with draggable SVG graph
func WriteInteractiveHTML[K micrograd.BaseNumeric](node micrograd.Numeric[K], filePath string, opts ...PlotOption[K]) error {
	// Set up configuration with defaults
	cfg := &plotConfig[K]{
		labelFunc: defaultNodeLabel[K],
	}

	// Apply options
	for _, opt := range opts {
		opt(cfg)
	}

	// Generate DOT
	dot := dotFromValue(node, cfg)

	// Read HTML template from embedded file
	templateContent, err := templateFS.ReadFile("template.html")
	if err != nil {
		return fmt.Errorf("failed to read embedded template file: %v", err)
	}

	// Replace placeholder with DOT content
	html := strings.Replace(string(templateContent), "<!-- DOT_CONTENT -->", dot, 1)

	// Write the output file
	return os.WriteFile(filePath, []byte(html), 0644)
}
