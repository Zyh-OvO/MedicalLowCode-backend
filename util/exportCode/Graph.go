package exportCode

import (
	"MedicalLowCode-backend/util/layer"
	"encoding/json"
	"gopkg.in/gyuho/goraph.v2"
)

type RawCanvas struct {
	Nodes []CNode `json:"nodes"`
	Edges []CEdge `json:"edges"`
}

type CNode struct {
	Id    string `json:"id"`
	Type  string `json:"type"`
	Data  any    `json:"data"` //todo
	Layer layer.Layer
}

type CEdge struct {
	SrcId string `json:"source"`
	TgtId string `json:"target"`
}

func (node *CNode) ID() goraph.ID {
	return goraph.StringID(node.Id)
}

func (node *CNode) String() string {
	return node.Id
}

func (node *CNode) SetLayer(l layer.Layer) {
	node.Layer = l
}

func TopologicalSort(graph goraph.Graph) ([]*CNode, bool) {
	nodeIds, ok := goraph.TopologicalSort(graph)
	if !ok {
		return nil, false
	}
	var canvasNodes []*CNode
	for _, node := range nodeIds {
		canvasNodes = append(canvasNodes, graph.GetNode(node).(*CNode))
	}
	return canvasNodes, true
}

func RecoverGraph(canvasContent string) goraph.Graph {
	graph := goraph.NewGraph()
	var canvas RawCanvas
	if err := json.Unmarshal([]byte(canvasContent), &canvas); err != nil {
		panic(err)
	}
	for key, _ := range canvas.Nodes {
		graph.AddNode(&canvas.Nodes[key])
	}
	for _, edge := range canvas.Edges {
		err := graph.AddEdge(goraph.StringID(edge.SrcId), goraph.StringID(edge.TgtId), 1)
		if err != nil {
			panic(err)
		}
	}
	return graph
}
