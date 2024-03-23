package layer

import (
	"errors"
	"gopkg.in/gyuho/goraph.v2"
)

func ExportCode(canvasContent string) string {
	graph := RecoverGraph(canvasContent)
	code, err := generateCode(graph)
	if err != nil {
		panic(err)
	}
	return code
}

func generateCode(graph goraph.Graph) (string, error) {
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
