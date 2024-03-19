package layer

type Conv1d struct {
	InChannels  int
	OutChannels int
	KernelSize  int
	Stride      *int    `default:"1"`
	Padding     any     `default:"0"` //int or str
	PaddingMode *string `default:"zeros"`
	Dilation    *int    `default:"1"`
	Groups      *int    `default:"1"`
	Bias        *bool   `default:"true"`
}

func (c *Conv1d) IsLayer() {
	return
}

type Conv2d struct {
	InChannels  int
	OutChannels int
	KernelSize  []int
	Stride      []int   `default:"1"`
	Padding     any     `default:"0"` //int[] or str
	PaddingMode *string `default:"zeros"`
	Dilation    []int   `default:"1"`
	Groups      *int    `default:"1"`
	Bias        *bool   `default:"true"`
}

func (c *Conv2d) IsLayer() {
	return
}

type Conv3d struct {
	InChannels  int
	OutChannels int
	KernelSize  []int
	Stride      []int   `default:"1"`
	Padding     any     `default:"0"` //int[] or str
	PaddingMode *string `default:"zeros"`
	Dilation    []int   `default:"1"`
	Groups      *int    `default:"1"`
	Bias        *bool   `default:"true"`
}

func (c *Conv3d) IsLayer() {
	return
}

func GenerateConvLayer(node *CNode) {
	switch node.Type {
	case "Conv1d":
		node.SetLayer(GenerateLayer(&Conv1d{}, node.Data.(map[string]any)))
	case "Conv2d":
		node.SetLayer(GenerateLayer(&Conv2d{}, node.Data.(map[string]any)))
	case "Conv3d":
		node.SetLayer(GenerateLayer(&Conv3d{}, node.Data.(map[string]any)))
	}
}
