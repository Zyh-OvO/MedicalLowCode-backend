package exportCode

var DatasetKinds = []string{"Dataset"}

type Dataset struct {
	TrainDataFilePath  string
	TrainLabelFilePath string
	BatchSize          int
	Shuffle            bool
	TestDataFilePath   string
	TestLabelFilePath  string
}

func (d *Dataset) IsLayer() {
	return
}

func GenerateDataset(node *CNode) Layer {
	switch node.Type {
	case "Dataset":
		return RawData2Layer(&Dataset{}, node.Data.(map[string]any))
	default:
		panic("unknown exportCode type")
	}
}
