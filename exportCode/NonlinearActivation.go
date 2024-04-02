package exportCode

var NonlinearActivationKinds = []string{"ELU", "Hardshrink", "Hardsigmoid", "Hardtanh", "Hardswish", "LeakyReLU", "LogSigmoid", "PReLU", "ReLU", "ReLU6", "RReLU", "SELU", "CELU", "GELU", "Sigmoid", "SiLU", "Mish", "Softplus", "Softshrink", "Softsign", "Tanh", "Tanhshrink", "Threshold", "GLU", "Softmin", "Softmax", "Softmax2d", "LogSoftmax"}

type ELU struct {
	Alpha   *float64 `default:"1.0"`
	Inplace *bool    `default:"false"`
}

func (e *ELU) IsLayer() {
	return
}

type Hardshrink struct {
	Lambda *float64 `default:"0.5"`
}

func (h *Hardshrink) IsLayer() {
	return
}

type Hardsigmoid struct {
	Inplace *bool `default:"false"`
}

func (h *Hardsigmoid) IsLayer() {
	return
}

type Hardtanh struct {
	MinVal  *float64 `default:"-1"`
	MaxVal  *float64 `default:"1"`
	Inplace *bool    `default:"false"`
}

func (h *Hardtanh) IsLayer() {
	return
}

type Hardswish struct {
	Inplace *bool `default:"false"`
}

func (h *Hardswish) IsLayer() {
	return
}

type LeakyReLU struct {
	NegativeSlope *float64 `default:"0.01"`
	Inplace       *bool    `default:"false"`
}

func (l *LeakyReLU) IsLayer() {
	return
}

type LogSigmoid struct {
}

func (l *LogSigmoid) IsLayer() {
	return
}

type PReLU struct {
	NumParameters *int     `default:"1"`
	Init          *float64 `default:"0.25"`
}

func (p *PReLU) IsLayer() {
	return
}

type ReLU struct {
	Inplace *bool `default:"false"`
}

func (r *ReLU) IsLayer() {
	return
}

type ReLU6 struct {
	Inplace *bool `default:"false"`
}

func (r *ReLU6) IsLayer() {
	return
}

type RReLU struct {
	Lower   *float64 `default:"0.125"`
	Upper   *float64 `default:"0.3333333333333333"`
	Inplace *bool    `default:"false"`
}

func (r *RReLU) IsLayer() {
	return
}

type SELU struct {
	Inplace *bool `default:"false"`
}

func (s *SELU) IsLayer() {
	return
}

type CELU struct {
	Alpha   *float64 `default:"1.0"`
	Inplace *bool    `default:"false"`
}

func (c *CELU) IsLayer() {
	return
}

type GELU struct {
	Approximate *string `default:"none"`
}

func (g *GELU) IsLayer() {
	return
}

type Sigmoid struct {
}

func (s *Sigmoid) IsLayer() {
	return
}

type SiLU struct {
	Inplace *bool `default:"false"`
}

func (s *SiLU) IsLayer() {
	return
}

type Mish struct {
	Inplace *bool `default:"false"`
}

func (m *Mish) IsLayer() {
	return
}

type Softplus struct {
	Beta      *int `default:"1"`
	Threshold *int `default:"20"`
}

func (s *Softplus) IsLayer() {
	return
}

type Softshrink struct {
	Lambda *float64 `default:"0.5"`
}

func (s *Softshrink) IsLayer() {
	return
}

type Softsign struct {
}

func (s *Softsign) IsLayer() {
	return
}

type Tanh struct {
}

func (t *Tanh) IsLayer() {
	return
}

type Tanhshrink struct {
}

func (t *Tanhshrink) IsLayer() {
	return
}

type Threshold struct {
	Threshold float64
	Value     float64
	Inplace   *bool `default:"false"`
}

func (t *Threshold) IsLayer() {
	return
}

type GLU struct {
	Dim *int `default:"1"`
}

func (g *GLU) IsLayer() {
	return
}

type Softmin struct {
	Dim *int `default:"nil"`
}

func (s *Softmin) IsLayer() {
	return
}

type Softmax struct {
	Dim *int `default:"nil"`
}

func (s *Softmax) IsLayer() {
	return
}

type Softmax2d struct {
}

func (s *Softmax2d) IsLayer() {
	return
}

type LogSoftmax struct {
	Dim *int `default:"nil"`
}

func (l *LogSoftmax) IsLayer() {
	return
}

func GenerateNonlinearActivation(node *CNode) Layer {
	switch node.Type {
	case "ELU":
		return RawData2Layer(&ELU{}, node.Data.(map[string]any))
	case "Hardshrink":
		return RawData2Layer(&Hardshrink{}, node.Data.(map[string]any))
	case "Hardsigmoid":
		return RawData2Layer(&Hardsigmoid{}, node.Data.(map[string]any))
	case "Hardtanh":
		return RawData2Layer(&Hardtanh{}, node.Data.(map[string]any))
	case "Hardswish":
		return RawData2Layer(&Hardswish{}, node.Data.(map[string]any))
	case "LeakyReLU":
		return RawData2Layer(&LeakyReLU{}, node.Data.(map[string]any))
	case "LogSigmoid":
		return RawData2Layer(&LogSigmoid{}, node.Data.(map[string]any))
	case "PReLU":
		return RawData2Layer(&PReLU{}, node.Data.(map[string]any))
	case "ReLU":
		return RawData2Layer(&ReLU{}, node.Data.(map[string]any))
	case "ReLU6":
		return RawData2Layer(&ReLU6{}, node.Data.(map[string]any))
	case "RReLU":
		return RawData2Layer(&RReLU{}, node.Data.(map[string]any))
	case "SELU":
		return RawData2Layer(&SELU{}, node.Data.(map[string]any))
	case "CELU":
		return RawData2Layer(&CELU{}, node.Data.(map[string]any))
	case "GELU":
		return RawData2Layer(&GELU{}, node.Data.(map[string]any))
	case "Sigmoid":
		return RawData2Layer(&Sigmoid{}, node.Data.(map[string]any))
	case "SiLU":
		return RawData2Layer(&SiLU{}, node.Data.(map[string]any))
	case "Mish":
		return RawData2Layer(&Mish{}, node.Data.(map[string]any))
	case "Softplus":
		return RawData2Layer(&Softplus{}, node.Data.(map[string]any))
	case "Softshrink":
		return RawData2Layer(&Softshrink{}, node.Data.(map[string]any))
	case "Softsign":
		return RawData2Layer(&Softsign{}, node.Data.(map[string]any))
	case "Tanh":
		return RawData2Layer(&Tanh{}, node.Data.(map[string]any))
	case "Tanhshrink":
		return RawData2Layer(&Tanhshrink{}, node.Data.(map[string]any))
	case "Threshold":
		return RawData2Layer(&Threshold{}, node.Data.(map[string]any))
	case "GLU":
		return RawData2Layer(&GLU{}, node.Data.(map[string]any))
	case "Softmin":
		return RawData2Layer(&Softmin{}, node.Data.(map[string]any))
	case "Softmax":
		return RawData2Layer(&Softmax{}, node.Data.(map[string]any))
	case "Softmax2d":
		return RawData2Layer(&Softmax2d{}, node.Data.(map[string]any))
	case "LogSoftmax":
		return RawData2Layer(&LogSoftmax{}, node.Data.(map[string]any))
	default:
		panic("unknown exportCode type")
	}
}
