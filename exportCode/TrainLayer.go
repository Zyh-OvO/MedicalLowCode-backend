package exportCode

var TrainLayerKinds = []string{"TrainLayer"}

type TrainLayer struct {
	NumEpochs         int
	SaveParamsDirPath string
}

func (t *TrainLayer) IsLayer() {
	return
}

func GenerateTrainLayer(node *CNode) Layer {
	switch node.Type {
	case "TrainLayer":
		return RawData2Layer(&TrainLayer{}, node.Data.(map[string]any))
	default:
		panic("unknown exportCode type")
	}
}
