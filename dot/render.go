package dot

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"task-lineage-diagram/schema"

	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
)

func strSliceToStrSet(strSlice []string) []string {
	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func makeDepSet(task schema.Task) []string {
	var depArr []string
	// find Dependency
	tgtDepArr := task.Dependency
	// extract its dependency array
	for _, dep := range tgtDepArr {
		depArr = append(depArr, dep.TaskID)
	}
	// remove duplicates
	return strSliceToStrSet(depArr)
}

func getColor(layer string, colorMap map[string]string) string {
	// Try to match the layer using strings.HasPrefix
	for key, color := range colorMap {
		if strings.HasPrefix(layer, key) {
			return color
		}
	}

	// Fallback to default color if no match is found
	if defaultColor, exists := colorMap["default"]; exists {
		return defaultColor
	}

	return "#FFFFFF" // Backup fallback in case default is missing in YAML
}

type nodeReachabilityTable map[string]int

func newNodeReachabilityTable(nodes map[string]*PNode) nodeReachabilityTable {
	route := make(nodeReachabilityTable)
	for key := range nodes {
		route[key] = 0
	}
	return route
}

func (a *nodeReachabilityTable) Add(b nodeReachabilityTable) nodeReachabilityTable {
	route := make(nodeReachabilityTable)
	for key, val := range *a {
		route[key] = val + b[key]
	}
	return route
}

type PNode struct {
	node             *cgraph.Node
	children         []string
	parents          []string
	visited          bool
	nodeReachability nodeReachabilityTable // a lookup map for reachable nodes
	reachableEdges   []string              // lookup set for all reachable edges
	incomingEdges    []string
	outgoingEdges    []string
}

type PNodeOut struct {
	Children         []string              `json:"children"`
	Parents          []string              `json:"parents"`
	NodeReachability nodeReachabilityTable `json:"nodeReachability"` // a lookup map for reachable nodes
	ReachableEdges   []string              `json:"reachableEdges"`   // lookup set for all reachable edges
	IncomingEdges    []string              `json:"incomingEdges"`
	OutgoingEdges    []string              `json:"outgoingEdges"`
}

func getReachability(key string, nodesMapPtr *map[string]*PNode) (nodeReachabilityTable, []string) {
	node := (*nodesMapPtr)[key]
	if !node.visited {
		node.visited = true
		// new a lookup table
		nodeReachability := newNodeReachabilityTable(*nodesMapPtr)
		reachableEdges := node.outgoingEdges
		// iterate all nodes, if the node is in the children, it is 1, otherwise 0
		for _, child_key := range node.children {
			childNodeReachability, childEdgeReachability := getReachability(child_key, nodesMapPtr)
			nodeReachability[child_key] = 1
			nodeReachability = nodeReachability.Add(childNodeReachability)
			reachableEdges = append(reachableEdges, childEdgeReachability...)
		}
		node.nodeReachability = nodeReachability
		node.reachableEdges = strSliceToStrSet(reachableEdges)
	}
	return node.nodeReachability, node.reachableEdges
}

func Render(tasks map[string]schema.Task, filename string, config *schema.Config, format graphviz.Format, layout graphviz.Layout, groupCluster bool, color bool, size string, reachFile string, noReach bool) {
	g := graphviz.New()
	g.SetLayout(layout)
	graph, err := g.Graph()
	processNodes := make(map[string]*PNode)
	subGraphs := make(map[string]*cgraph.Graph)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := graph.Close(); err != nil {
			log.Fatal(err)
		}
		g.Close()
	}()

	// create nodes
	for key, task := range tasks {
		if task.Level != "" {
			cluster := task.Level
			// check if SubGraph exists, if not create one.
			_, subGraphExists := subGraphs[cluster]
			if !subGraphExists {
				subGraph := graph.SubGraph("cluster_"+cluster, 1)
				subGraph.SetID(cluster)
				subGraph.SafeSet("peripheries", "0", "")
				subGraphs[cluster] = subGraph
			}
			// call again to make sure it is there
			subGraph, _ := subGraphs[cluster]
			// create node in the subgraph
			node, err := subGraph.CreateNode(key)
			if err != nil {
				log.Fatal(err)
			}
			node.SetLayer(task.Level)
			node.SetLabel(task.Task)
			node.SetXLabel(task.Level)
			node.SetShape("box")
			node.SetStyle("rounded")
			node.SetFontSize(20)
			node.SetID(key)
			if color {
				color := getColor(task.Level, config.Colors)
				node.SetColor(color)
				node.SetFontColor(color)
			}
			node.SafeSet("fontname", "Arial", "")
			// TODO can use this with css to style
			// node.SafeSet("class", key, "")

			var children []string
			// store nodes
			processNodes[key] = &PNode{
				node:     node,
				children: children,
				visited:  false,
			}

		}
	}

	// close subgraphs
	// for _, subgraph := range subGraphs {
	// 	subgraph.Close()
	// }

	// create edges
	edgeCount := 0
	for key, task := range tasks {
		deps := makeDepSet(task)
		for _, dep := range deps {
			_, node1Exists := processNodes[dep]
			_, node2Exists := processNodes[key]
			if node1Exists && node2Exists && processNodes[dep].node != processNodes[key].node {
				edgeID := fmt.Sprintf("edge%d", edgeCount)
				// set edge
				edgeName := dep + "->" + key
				e, err := graph.CreateEdge(edgeName, processNodes[dep].node, processNodes[key].node)
				if err != nil {
					log.Fatal(err)
				}
				if color {
					e.SetColor("#737D82")
				}
				e.SetTooltip(edgeName)
				e.SetID(edgeID)
				// TODO can use this with css to style
				// e.SafeSet("class", key, "")

				processNodes[key].parents = append(processNodes[key].parents, dep)
				processNodes[key].incomingEdges = append(processNodes[key].incomingEdges, edgeID)

				// set node children field
				processNodes[dep].children = append(processNodes[dep].children, key)
				processNodes[dep].outgoingEdges = append(processNodes[dep].outgoingEdges, edgeID)

				edgeCount++
			}
		}
	}

	// set graph
	graph.SetID("RootGraph")
	graph.SetSmoothing("graph_dist")
	graph.SetOutputOrder("edgesfirst")
	graph.SetNoTranslate(true)
	graph.SetRankDir("LR")
	graph.SetRankSeparator(4)
	graph.SetPad(3)
	graph.SetOverlap(false)
	graph.SetCenter(true)
	if color {
		graph.SetBackgroundColor("#323237")
		graph.SetFontColor("#F3F3F3")
	}
	graph.SafeSet("fontname", "Arial", "")
	graph.SetLabel("Task Lineage Diagram")
	graph.SetFontSize(120)
	graph.SetLabelLocation("t")
	graph.SetLabelJust("r")
	if groupCluster {
		graph.SetClusterRank("local")
	} else {
		graph.SetClusterRank("global")
	}

	printSize := strings.ToLower(size)
	if printSize == "a3" { // A3
		graph.SetSize(11.7, 16.5)
	} else if printSize == "fhd" { // full hd
		graph.SetSize(20, 11.25)
	} else { // full hd
		graph.SetSize(20, 11.25)
	}

	fmt.Printf("Found %d nodes(processes) and %d edges(dependencies).\n", len(processNodes), edgeCount)

	// write to file directly
	fmt.Printf("Rendering to file \"%s\" ...\n", filename)
	if err := g.RenderFilename(graph, format, filename); err != nil {
		log.Fatal(err)
	}

	if !noReach {
		// do analysis
		fmt.Println("Analyzing graph reachability...")
		totalReachability := make(map[string]PNodeOut)
		for key := range processNodes {
			parents, children, incomingE, outgoingE := processNodes[key].parents, processNodes[key].children, processNodes[key].incomingEdges, processNodes[key].outgoingEdges
			nodeReachability, reachableEdges := getReachability(key, &processNodes)
			totalReachability[key] = PNodeOut{
				Parents:          parents,
				Children:         children,
				IncomingEdges:    incomingE,
				OutgoingEdges:    outgoingE,
				NodeReachability: nodeReachability,
				ReachableEdges:   reachableEdges,
			}
		}
		jsonData, err := json.Marshal(totalReachability)
		if err != nil {
			log.Printf("Error: %s", err.Error())
		}

		// write to JSON file
		fmt.Println("Writing reachability report...")
		jsonFile, err := os.Create(reachFile)

		if err != nil {
			panic(err)
		}
		defer jsonFile.Close()

		jsonFile.Write(jsonData)
		jsonFile.Close()
		fmt.Printf("Reachability data written to \"%s\"\n", reachFile)
	} else {
		fmt.Println("Skipping reachability analysis.")
	}
}
