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

func GenerateLossFunction(node *CNode) {
	switch node.Type {
	case "L1Loss":
		node.SetLayer(GenerateLayer(&L1Loss{}, node.Data.(map[string]any)))
	case "MSELoss":
		node.SetLayer(GenerateLayer(&MSELoss{}, node.Data.(map[string]any)))
	case "CrossEntropyLoss":
		node.SetLayer(GenerateLayer(&CrossEntropyLoss{}, node.Data.(map[string]any)))
	case "BCELoss":
		node.SetLayer(GenerateLayer(&BCELoss{}, node.Data.(map[string]any)))
	}
}
