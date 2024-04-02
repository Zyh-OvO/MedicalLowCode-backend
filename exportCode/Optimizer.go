package exportCode

var OptimizerKinds = []string{"Adadelta", "Adagrad", "Adam", "AdamW", "SparseAdam", "Adamax", "ASGD", "LBFGS", "NAdam", "RAdam", "RMSprop", "Rprop", "SGD"}

type Adadelta struct {
	Rho            *float64
	Eps            *float64
	Lr             *float64
	WeightDecay    *float64
	Foreach        *bool
	Maximize       *bool
	Differentiable *bool
}

func (a *Adadelta) IsLayer() {
	return
}

type Adagrad struct {
	Lr             *float64
	LrDecay        *float64
	WeightDecay    *float64
	InitialAcc     *float64
	Eps            *float64
	Foreach        *bool
	Maximize       *bool
	Differentiable *bool
}

func (a *Adagrad) IsLayer() {
	return
}

type Adam struct {
	Lr             *float64
	Betas          []float64
	Eps            *float64
	WeightDecay    *float64
	Amsgrad        *bool
	Foreach        *bool
	Maximize       *bool
	Capturable     *bool
	Differentiable *bool
	Fused          *bool
}

func (a *Adam) IsLayer() {
	return
}

type AdamW struct {
	Lr             *float64
	Betas          []float64
	Eps            *float64
	WeightDecay    *float64
	Amsgrad        *bool
	Foreach        *bool
	Maximize       *bool
	Capturable     *bool
	Differentiable *bool
	Fused          *bool
}

func (a *AdamW) IsLayer() {
	return
}

type SparseAdam struct {
	Lr       *float64
	Betas    []float64
	Eps      *float64
	Maximize *bool
}

func (s *SparseAdam) IsLayer() {
	return
}

type Adamax struct {
	Lr             *float64
	Betas          []float64
	Eps            *float64
	WeightDecay    *float64
	Foreach        *bool
	Maximize       *bool
	Differentiable *bool
}

func (a *Adamax) IsLayer() {
	return
}

type ASGD struct {
	Lr             *float64
	Lambd          *float64
	Alpha          *float64
	T0             *float64
	WeightDecay    *float64
	Foreach        *bool
	Maximize       *bool
	Differentiable *bool
	Capturable     *bool
}

func (a *ASGD) IsLayer() {
	return
}

type LBFGS struct {
	Lr              *float64
	MaxIter         *int
	MaxEval         *int
	ToleranceGrad   *float64
	ToleranceChange *float64
	HistorySize     *int
	LineSearchFn    *string
}

func (l *LBFGS) IsLayer() {
	return
}

type NAdam struct {
	Lr                   *float64
	Betas                []float64
	Eps                  *float64
	WeightDecay          *float64
	MomentumDecay        *float64
	DecoupledWeightDecay *bool
	Foreach              *bool
	Capturable           *bool
	Differentiable       *bool
}

func (n *NAdam) IsLayer() {
	return
}

type RAdam struct {
	Lr                   *float64
	Betas                []float64
	Eps                  *float64
	WeightDecay          *float64
	DecoupledWeightDecay *bool
	Foreach              *bool
	Differentiable       *bool
}

func (r *RAdam) IsLayer() {
	return
}

type RMSprop struct {
	Lr             *float64
	Momentum       *float64
	Alpha          *float64
	Eps            *float64
	Centered       *bool
	WeightDecay    *float64
	Foreach        *bool
	Maximize       *bool
	Differentiable *bool
}

func (r *RMSprop) IsLayer() {
	return
}

type Rprop struct {
	Lr             *float64
	Etas           []float64
	StepSizes      []float64
	Foreach        *bool
	Maximize       *bool
	Differentiable *bool
}

func (r *Rprop) IsLayer() {
	return
}

type SGD struct {
	Lr             *float64
	Momentum       *float64
	WeightDecay    *float64
	Dampening      *float64
	Nesterov       *bool
	Maximize       *bool
	Foreach        *bool
	Differentiable *bool
}

func (s *SGD) IsLayer() {
	return
}

func GenerateOptimizer(node *CNode) Layer {
	switch node.Type {
	case "Adadelta":
		return RawData2Layer(&Adadelta{}, node.Data.(map[string]any))
	case "Adagrad":
		return RawData2Layer(&Adagrad{}, node.Data.(map[string]any))
	case "Adam":
		return RawData2Layer(&Adam{}, node.Data.(map[string]any))
	case "AdamW":
		return RawData2Layer(&AdamW{}, node.Data.(map[string]any))
	case "SparseAdam":
		return RawData2Layer(&SparseAdam{}, node.Data.(map[string]any))
	case "Adamax":
		return RawData2Layer(&Adamax{}, node.Data.(map[string]any))
	case "ASGD":
		return RawData2Layer(&ASGD{}, node.Data.(map[string]any))
	case "LBFGS":
		return RawData2Layer(&LBFGS{}, node.Data.(map[string]any))
	case "NAdam":
		return RawData2Layer(&NAdam{}, node.Data.(map[string]any))
	case "RAdam":
		return RawData2Layer(&RAdam{}, node.Data.(map[string]any))
	case "RMSprop":
		return RawData2Layer(&RMSprop{}, node.Data.(map[string]any))
	case "Rprop":
		return RawData2Layer(&Rprop{}, node.Data.(map[string]any))
	case "SGD":
		return RawData2Layer(&SGD{}, node.Data.(map[string]any))
	default:
		panic("unknown exportCode type")
	}
}
