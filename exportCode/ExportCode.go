package exportCode

import (
	"errors"
	"gopkg.in/gyuho/goraph.v2"
)

func ExportCode(canvasContent string) string {
	//数据处理代码
	dataCode, err := genDataCode()
	if err != nil {
		panic(err)
	}
	//网络代码
	graph := RecoverGraph(canvasContent)
	netCode, err := genNetCode(graph)
	if err != nil {
		panic(err)
	}
	//训练代码
	trainCode, err := genTrainCode()
	if err != nil {
		panic(err)
	}
	return dataCode + netCode + trainCode
}

func genDataCode() (string, error) {
	//todo
	return "", nil
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
		node.GenerateLayer(layerCounter[node.Type])
		layerCode, forwardCode := node.GenerateCode()
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
	return code, nil
}

func genTrainCode() (string, error) {
	//todo
	return "", nil
}
