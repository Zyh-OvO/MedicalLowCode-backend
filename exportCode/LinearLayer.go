package exportCode

var LinearLayerKinds = []string{"Linear", "Bilinear", "LazyLinear"}

type Linear struct {
	InFeatures  int
	OutFeatures int
	Bias        *bool `default:"true"`
}

func (l *Linear) IsLayer() {
	return
}

type Bilinear struct {
	In1Features int
	In2Features int
	OutFeatures int
	Bias        *bool `default:"true"`
}

func (b *Bilinear) IsLayer() {
	return
}

type LazyLinear struct {
	OutFeatures int
	Bias        *bool `default:"true"`
}

func (l *LazyLinear) IsLayer() {
	return
}

func GenerateLinearLayer(node *CNode) Layer {
	switch node.Type {
	case "Linear":
		return RawData2Layer(&Linear{}, node.Data.(map[string]any))
	case "Bilinear":
		return RawData2Layer(&Bilinear{}, node.Data.(map[string]any))
	case "LazyLinear":
		return RawData2Layer(&LazyLinear{}, node.Data.(map[string]any))
	default:
		panic("unknown exportCode type")
	}
}
