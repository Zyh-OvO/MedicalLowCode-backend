package layer

type MaxPool1d struct {
	KernelSize    int
	Stride        int
	Padding       *int  `default:"0"`
	Dilation      *int  `default:"1"`
	ReturnIndices *bool `default:"false"`
	CeilMode      *bool `default:"false"`
}

func (m *MaxPool1d) IsLayer() {
	return
}

type MaxPool2d struct {
	KernelSize    []int
	Stride        []int
	Padding       []int `default:"0"`
	Dilation      []int `default:"1"`
	ReturnIndices *bool `default:"false"`
	CeilMode      *bool `default:"false"`
}

func (m *MaxPool2d) IsLayer() {
	return
}

type MaxPool3d struct {
	KernelSize    []int
	Stride        []int
	Padding       []int `default:"0"`
	Dilation      []int `default:"1"`
	ReturnIndices *bool `default:"false"`
	CeilMode      *bool `default:"false"`
}

func (m *MaxPool3d) IsLayer() {
	return
}

type AvgPool1d struct {
	KernelSize      int
	Stride          int
	Padding         *int  `default:"0"`
	CeilMode        *bool `default:"false"`
	CountIncludePad *bool `default:"true"`
}

func (a *AvgPool1d) IsLayer() {
	return
}

type AvgPool2d struct {
	KernelSize      []int
	Stride          []int
	Padding         []int `default:"0"`
	CeilMode        *bool `default:"false"`
	CountIncludePad *bool `default:"true"`
	DivisorOverride *int  `default:"nil"`
}

func (a *AvgPool2d) IsLayer() {
	return
}

type AvgPool3d struct {
	KernelSize      []int
	Stride          []int
	Padding         []int `default:"0"`
	CeilMode        *bool `default:"false"`
	CountIncludePad *bool `default:"true"`
	DivisorOverride *int  `default:"nil"`
}

func (a *AvgPool3d) IsLayer() {
	return
}

func GeneratePoolingLayer(node *CNode) {
	switch node.Type {
	case "MaxPool1d":
		node.SetLayer(GenerateLayer(&MaxPool1d{}, node.Data.(map[string]any)))
	case "MaxPool2d":
		node.SetLayer(GenerateLayer(&MaxPool2d{}, node.Data.(map[string]any)))
	case "MaxPool3d":
		node.SetLayer(GenerateLayer(&MaxPool3d{}, node.Data.(map[string]any)))
	case "AvgPool1d":
		node.SetLayer(GenerateLayer(&AvgPool1d{}, node.Data.(map[string]any)))
	case "AvgPool2d":
		node.SetLayer(GenerateLayer(&AvgPool2d{}, node.Data.(map[string]any)))
	case "AvgPool3d":
		node.SetLayer(GenerateLayer(&AvgPool3d{}, node.Data.(map[string]any)))
	}
}
