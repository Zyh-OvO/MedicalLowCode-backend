package exportCode

import (
	"MedicalLowCode-backend/util"
	"errors"
	"gopkg.in/gyuho/goraph.v2"
	"strconv"
)

func ExportCode(canvasContent string) string {
	datasetGraph, netGraph, trainGraph := RecoverGraph(canvasContent)
	//数据集代码
	dataCode, err := genDataCode(datasetGraph)
	if err != nil {
		panic(err)
	}
	//网络代码
	netCode, err := genNetCode(netGraph)
	if err != nil {
		panic(err)
	}
	//训练代码
	trainCode, err := genTrainCode(trainGraph)
	if err != nil {
		panic(err)
	}
	return dataCode + netCode + trainCode
}

func genDataCode(graph goraph.Graph) (string, error) {
	var code string
	topologicalNodes, ok := TopologicalSort(graph)
	if !ok {
		return "", errors.New("topological sort failed")
	}
	for _, node := range topologicalNodes {
		node.GenerateLayer()
		code += node.GenerateDatasetCode()
	}
	return code, nil
}

func genNetCode(graph goraph.Graph) (string, error) {
	var code string
	topologicalNodes, ok := TopologicalSort(graph)
	if !ok {
		return "", errors.New("topological sort failed")
	}
	//生成两部分代码
	var initLayerCodes, forwardCodes []string
	layerCounter := make(map[string]int)
	for _, node := range topologicalNodes {
		layerCounter[node.Type]++
		//先生成layer
		node.GenerateLayer()
		node.LayerName = node.Type + "_" + strconv.Itoa(layerCounter[node.Type])
		node.OutputName = node.LayerName + "_output"
		node.OutputName = util.CamelCaseToSnakeCase(node.OutputName)
		//再生成code
		layerCode, forwardCode := node.GenerateNetCode()
		initLayerCodes = append(initLayerCodes, layerCode)
		forwardCodes = append(forwardCodes, forwardCode)
	}
	//生成整个网络的代码
	code += "import torch.nn as nn\n\n"
	code += "class Net(nn.Module):\n"
	code += "    def __init__(self):\n"
	code += "        super().__init__()\n"
	for _, layerCode := range initLayerCodes {
		code += "        " + layerCode + "\n"
	}
	code += "    def forward(self, x):\n"
	for _, forwardCode := range forwardCodes {
		code += "        " + forwardCode + "\n"
	}
	code += "        return " + topologicalNodes[len(topologicalNodes)-1].OutputName + "\n"
	code += "model = Net()\n"
	return code, nil
}

func genTrainCode(graph goraph.Graph) (string, error) {
	var code string
	topologicalNodes, ok := TopologicalSort(graph)
	if !ok {
		return "", errors.New("topological sort failed")
	}
	for _, node := range topologicalNodes {
		node.GenerateLayer()
		code += node.GenerateTrainCode()
	}
	return code, nil
}
