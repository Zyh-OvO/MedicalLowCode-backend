package layer

import (
	"MedicalLowCode-backend/util"
	"encoding/json"
	"gopkg.in/gyuho/goraph.v2"
	"strconv"
)

type RawCanvas struct {
	Nodes []CNode `json:"nodes"`
	Edges []CEdge `json:"edges"`
}

type CNode struct {
	Id           string `json:"id"`
	Type         string `json:"layerType"`
	Data         any    `json:"data"` //map[string]any
	Layer        Layer
	LayerName    string //python中网络层变量名
	OutputName   string //网络层输出变量名
	Predecessors []*CNode
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

func (node *CNode) GenerateLayer(index int) {
	switch node.Type {
	case "Conv1d", "Conv2d", "Conv3d":
		node.Layer = GenerateConvLayer(node)
	case "Linear", "Bilinear", "LazyLinear":
		node.Layer = GenerateLinearLayer(node)
	case "L1Loss", "MSELoss", "CrossEntropyLoss", "BCELoss":
		node.Layer = GenerateLossFunction(node)
	case "ELU", "Hardshrink", "Hardsigmoid", "Hardtanh", "Hardswish", "LeakyReLU", "LogSigmoid", "PReLU", "ReLU", "ReLU6", "RReLU", "SELU", "CELU", "GELU", "Sigmoid", "SiLU", "Mish", "Softplus", "Softshrink", "Softsign", "Tanh", "Tanhshrink", "Threshold", "GLU", "Softmin", "Softmax", "Softmax2d", "LogSoftmax":
		node.Layer = GenerateNonlinearActivation(node)
	case "MaxPool1d", "MaxPool2d", "MaxPool3d", "AvgPool1d", "AvgPool2d", "AvgPool3d":
		node.Layer = GeneratePoolingLayer(node)
	default:
		panic("unsupported layer type")
	}
	node.LayerName = node.Type + "_" + strconv.Itoa(index)
	node.OutputName = node.LayerName + "_output"
	node.OutputName = util.CamelCaseToSnakeCase(node.OutputName)
}

func (node *CNode) GenerateCode() (string, string) {
	// 两部分：1. 生成layer 2. 生成forward
	var layerCode, forwardCode string
	layerCode = Layer2Code(node.Layer, node.LayerName)
	forwardCode = node.OutputName + " = self." + node.LayerName + "("
	for _, predecessor := range node.Predecessors {
		forwardCode += predecessor.OutputName + ", "
	}
	//todo:如果没有前驱节点，默认输入x
	if len(node.Predecessors) == 0 {
		forwardCode += "x"
	}
	forwardCode += ")"
	return layerCode, forwardCode
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
		//todo: 维护前驱节点
		err := graph.AddEdge(goraph.StringID(edge.SrcId), goraph.StringID(edge.TgtId), 1)
		if err != nil {
			panic(err)
		}
	}
	return graph
}
