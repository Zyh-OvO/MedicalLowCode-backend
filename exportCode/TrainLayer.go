package exportCode

var TrainLayerKinds = []string{"TrainLayer", "ReasoningLayer"}

type TrainLayer struct {
	NumEpochs         int
	SaveParamsDirPath string
}

func (t *TrainLayer) IsLayer() {
	return
}

type ReasoningLayer struct {
	ParamsFilePath    string
	SaveResultDirPath string
}

func (r *ReasoningLayer) IsLayer() {
	return
}

func GenerateTrainLayer(node *CNode) Layer {
	switch node.Type {
	case "TrainLayer":
		return RawData2Layer(&TrainLayer{}, node.Data.(map[string]any))
	case "ReasoningLayer":
		return RawData2Layer(&ReasoningLayer{}, node.Data.(map[string]any))
	default:
		panic("unknown exportCode type")
	}
}
