package layer

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

func GenerateLinearLayer(node *CNode) {
	switch node.Type {
	case "Linear":
		node.SetLayer(GenerateLayer(&Linear{}, node.Data.(map[string]any)))
	case "Bilinear":
		node.SetLayer(GenerateLayer(&Bilinear{}, node.Data.(map[string]any)))
	case "LazyLinear":
		node.SetLayer(GenerateLayer(&LazyLinear{}, node.Data.(map[string]any)))
	}
}
