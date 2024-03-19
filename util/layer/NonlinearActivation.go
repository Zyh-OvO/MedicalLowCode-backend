package layer

import (
	"MedicalLowCode-backend/util/exportCode"
)

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

func GenerateNonlinearActivation(node *exportCode.CNode) {
	switch node.Type {
	case "ELU":
		node.SetLayer(GenerateLayer(&ELU{}, node.Data.(map[string]any)))
	case "Hardshrink":
		node.SetLayer(GenerateLayer(&Hardshrink{}, node.Data.(map[string]any)))
	case "Hardsigmoid":
		node.SetLayer(GenerateLayer(&Hardsigmoid{}, node.Data.(map[string]any)))
	case "Hardtanh":
		node.SetLayer(GenerateLayer(&Hardtanh{}, node.Data.(map[string]any)))
	case "Hardswish":
		node.SetLayer(GenerateLayer(&Hardswish{}, node.Data.(map[string]any)))
	case "LeakyReLU":
		node.SetLayer(GenerateLayer(&LeakyReLU{}, node.Data.(map[string]any)))
	case "LogSigmoid":
		node.SetLayer(GenerateLayer(&LogSigmoid{}, node.Data.(map[string]any)))
	case "PReLU":
		node.SetLayer(GenerateLayer(&PReLU{}, node.Data.(map[string]any)))
	case "ReLU":
		node.SetLayer(GenerateLayer(&ReLU{}, node.Data.(map[string]any)))
	case "ReLU6":
		node.SetLayer(GenerateLayer(&ReLU6{}, node.Data.(map[string]any)))
	case "RReLU":
		node.SetLayer(GenerateLayer(&RReLU{}, node.Data.(map[string]any)))
	case "SELU":
		node.SetLayer(GenerateLayer(&SELU{}, node.Data.(map[string]any)))
	case "CELU":
		node.SetLayer(GenerateLayer(&CELU{}, node.Data.(map[string]any)))
	case "GELU":
		node.SetLayer(GenerateLayer(&GELU{}, node.Data.(map[string]any)))
	case "Sigmoid":
		node.SetLayer(GenerateLayer(&Sigmoid{}, node.Data.(map[string]any)))
	case "SiLU":
		node.SetLayer(GenerateLayer(&SiLU{}, node.Data.(map[string]any)))
	case "Mish":
		node.SetLayer(GenerateLayer(&Mish{}, node.Data.(map[string]any)))
	case "Softplus":
		node.SetLayer(GenerateLayer(&Softplus{}, node.Data.(map[string]any)))
	case "Softshrink":
		node.SetLayer(GenerateLayer(&Softshrink{}, node.Data.(map[string]any)))
	case "Softsign":
		node.SetLayer(GenerateLayer(&Softsign{}, node.Data.(map[string]any)))
	case "Tanh":
		node.SetLayer(GenerateLayer(&Tanh{}, node.Data.(map[string]any)))
	case "Tanhshrink":
		node.SetLayer(GenerateLayer(&Tanhshrink{}, node.Data.(map[string]any)))
	case "Threshold":
		node.SetLayer(GenerateLayer(&Threshold{}, node.Data.(map[string]any)))
	case "GLU":
		node.SetLayer(GenerateLayer(&GLU{}, node.Data.(map[string]any)))
	case "Softmin":
		node.SetLayer(GenerateLayer(&Softmin{}, node.Data.(map[string]any)))
	case "Softmax":
		node.SetLayer(GenerateLayer(&Softmax{}, node.Data.(map[string]any)))
	case "Softmax2d":
		node.SetLayer(GenerateLayer(&Softmax2d{}, node.Data.(map[string]any)))
	case "LogSoftmax":
		node.SetLayer(GenerateLayer(&LogSoftmax{}, node.Data.(map[string]any)))

	}
}
