package exportCode

var DatasetKinds = []string{"TrainDataset", "ReasoningDataset"}

type TrainDataset struct {
	TrainDataFilePath  string
	TrainLabelFilePath string
	TestDataFilePath   string
	TestLabelFilePath  string
	BatchSize          int
	Shuffle            bool
}

func (d *TrainDataset) IsLayer() {
	return
}

type ReasoningDataset struct {
	ReasoningDataFilePath string
}

func (d *ReasoningDataset) IsLayer() {
	return
}

func GenerateDataset(node *CNode) Layer {
	switch node.Type {
	case "TrainDataset":
		return RawData2Layer(&TrainDataset{}, node.Data.(map[string]any))
	case "ReasoningDataset":
		return RawData2Layer(&ReasoningDataset{}, node.Data.(map[string]any))
	default:
		panic("unknown exportCode type")
	}
}
