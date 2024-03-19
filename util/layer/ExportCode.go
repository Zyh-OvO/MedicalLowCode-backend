package layer

import (
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
	code := "code"
	topologicalNodes, ok := TopologicalSort(graph)
	if !ok {
		panic("topological sort failed")
	}
	// 生成layer
	for _, node := range topologicalNodes {
		switch node.Type {
		case "Conv1d", "Conv2d", "Conv3d":
			GenerateConvLayer(node)
		case "Linear", "Bilinear", "LazyLinear":
			GenerateLinearLayer(node)
		case "L1Loss", "MSELoss", "CrossEntropyLoss", "BCELoss":
			GenerateLossFunction(node)
		case "ELU", "Hardshrink", "Hardsigmoid", "Hardtanh", "Hardswish", "LeakyReLU", "LogSigmoid", "PReLU", "ReLU", "ReLU6", "RReLU", "SELU", "CELU", "GELU", "Sigmoid", "SiLU", "Mish", "Softplus", "Softshrink", "Softsign", "Tanh", "Tanhshrink", "Threshold", "GLU", "Softmin", "Softmax", "Softmax2d", "LogSoftmax":
			GenerateNonlinearActivation(node)
		case "MaxPool1d", "MaxPool2d", "MaxPool3d", "AvgPool1d", "AvgPool2d", "AvgPool3d":
			GeneratePoolingLayer(node)
		}
	}

	return code, nil
}
