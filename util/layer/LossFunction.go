package layer

type L1Loss struct {
	SizeAverage *bool   `default:"true"`
	Reduce      *bool   `default:"true"`
	Reduction   *string `default:"mean"` //none, mean, sum
}

func (l *L1Loss) IsLayer() {
	return
}

type MSELoss struct {
	SizeAverage *bool   `default:"true"`
	Reduce      *bool   `default:"true"`
	Reduction   *string `default:"mean"` //none, mean, sum
}

func (m *MSELoss) IsLayer() {
	return
}

type CrossEntropyLoss struct {
	Weight         []float64 `default:"nil"`
	SizeAverage    *bool     `default:"true"`
	IgnoreIndex    *int      `default:"-100"`
	Reduce         *bool     `default:"true"`
	Reduction      *string   `default:"mean"` //none, mean, sum
	LabelSmoothing *float64  `default:"0"`
}

func (c *CrossEntropyLoss) IsLayer() {
	return
}

type BCELoss struct {
	Weight      []float64 `default:"nil"`
	SizeAverage *bool     `default:"true"`
	Reduce      *bool     `default:"true"`
	Reduction   *string   `default:"mean"` //none, mean, sum
}

func (b *BCELoss) IsLayer() {
	return
}

func GenerateLossFunction(node *CNode) Layer {
	switch node.Type {
	case "L1Loss":
		return RawData2Layer(&L1Loss{}, node.Data.(map[string]any))
	case "MSELoss":
		return RawData2Layer(&MSELoss{}, node.Data.(map[string]any))
	case "CrossEntropyLoss":
		return RawData2Layer(&CrossEntropyLoss{}, node.Data.(map[string]any))
	case "BCELoss":
		return RawData2Layer(&BCELoss{}, node.Data.(map[string]any))
	default:
		panic("unknown layer type")
	}
}
