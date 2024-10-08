package exportCode

import (
	"MedicalLowCode-backend/util"
	"encoding/json"
	"gopkg.in/gyuho/goraph.v2"
	"strconv"
	"strings"
)

type RawCanvas struct {
	Nodes []CNode `json:"nodes"`
	Edges []CEdge `json:"edges"`
}

type CNode struct {
	Id           string `json:"id"`
	Type         string `json:"layerType"`
	Data         any    `json:"data"` //map[string]any
	Layer        Layer
	LayerName    string //python中网络层变量名
	OutputName   string //网络层输出变量名
	Predecessors []*CNode
}

type CEdge struct {
	SrcId string `json:"source"`
	TgtId string `json:"target"`
}

func (node *CNode) ID() goraph.ID {
	return goraph.StringID(node.Id)
}

func (node *CNode) String() string {
	return node.Id
}

func (node *CNode) GenerateLayer() {
	nodeType := node.Type
	switch {
	case util.SliceContains(ConvolutionLayerKinds, nodeType):
		node.Layer = GenerateConvLayer(node)
	case util.SliceContains(LinearLayerKinds, nodeType):
		node.Layer = GenerateLinearLayer(node)
	case util.SliceContains(LossFunctionKinds, nodeType):
		node.Layer = GenerateLossFunction(node)
	case util.SliceContains(NonlinearActivationKinds, nodeType):
		node.Layer = GenerateNonlinearActivation(node)
	case util.SliceContains(PoolingLayerKinds, nodeType):
		node.Layer = GeneratePoolingLayer(node)
	case util.SliceContains(OptimizerKinds, nodeType):
		node.Layer = GenerateOptimizer(node)
	case util.SliceContains(DatasetKinds, nodeType):
		node.Layer = GenerateDataset(node)
	case util.SliceContains(TrainLayerKinds, nodeType):
		node.Layer = GenerateTrainLayer(node)
	default:
		panic("unsupported exportCode type")
	}
}

func (node *CNode) GenerateDatasetCode() string {
	var code string
	switch node.Layer.(type) {
	//目前只有TrainDataset和ReasoningDataset
	case *TrainDataset:
		dataset := node.Layer.(*TrainDataset)
		code += "class CustomDataset(Dataset):\n"
		code += "    def __init__(self, data_file, label_file):\n"
		code += "        self.data = pd.read_csv(data_file)\n"
		code += "        self.label = pd.read_csv(label_file)\n"
		code += "    def __len__(self):\n"
		code += "        return min(len(self.data), len(self.label))\n"
		code += "    def __getitem__(self, idx):\n"
		code += "        data_item = torch.from_numpy(self.data.iloc[idx].values.astype(float))\n"
		code += "        label_item = torch.from_numpy(self.label.iloc[idx].values.astype(float))\n"
		code += "        return data_item, label_item\n\n"
		code += "train_dataset = CustomDataset('" + dataset.TrainDataFilePath + "', '" + dataset.TrainLabelFilePath + "')\n"
		code += "train_loader = DataLoader(train_dataset, batch_size=" + strconv.Itoa(dataset.BatchSize) + ", shuffle=" + strings.Title(strconv.FormatBool(dataset.Shuffle)) + ")\n"
		code += "test_dataset = CustomDataset('" + dataset.TestDataFilePath + "', '" + dataset.TestLabelFilePath + "')\n"
		code += "test_loader = DataLoader(test_dataset, batch_size=" + strconv.Itoa(dataset.BatchSize) + ", shuffle=" + strings.Title(strconv.FormatBool(dataset.Shuffle)) + ")\n\n"
	case *ReasoningDataset:
		code += "reasoning_data_file = '" + node.Layer.(*ReasoningDataset).ReasoningDataFilePath + "'\n"
		code += "reasoning_data = pd.read_csv(reasoning_data_file)\n"
		code += "reasoning_data = torch.from_numpy(reasoning_data.values.astype(float))\n\n"
	}
	return code
}

func (node *CNode) GenerateNetCode() (string, string) {
	// 两部分：1. 生成layer 2. 生成forward
	var layerCode, forwardCode string
	layerCode += "self." + node.LayerName + " = nn." + node.Type + "("
	layerCode += Layer2Code(node.Layer)
	layerCode += ")"
	forwardCode = node.OutputName + " = self." + node.LayerName + "("
	for _, predecessor := range node.Predecessors {
		forwardCode += predecessor.OutputName + ", "
	}
	//todo:如果没有前驱节点，默认输入x
	if len(node.Predecessors) == 0 {
		forwardCode += "x"
	}
	forwardCode += ")"
	return layerCode, forwardCode
}

func (node *CNode) GenerateTrainCode() string {
	var code string
	nodeType := node.Type
	switch {
	case util.SliceContains(OptimizerKinds, nodeType):
		code += "optimizer = optim." + node.Type + "(model.parameters(), "
		code += Layer2Code(node.Layer)
		code += ")\n\n"
	case util.SliceContains(LossFunctionKinds, nodeType):
		code += "criterion = nn." + node.Type + "("
		code += Layer2Code(node.Layer)
		code += ")\n\n"
	case nodeType == "TrainLayer":
		trainLayer := node.Layer.(*TrainLayer)
		code += "num_epochs = " + strconv.Itoa(trainLayer.NumEpochs) + "\n"
		code += "save_params_file_path = '" + trainLayer.SaveParamsFilePath + "'\n"
		code += "best_loss = float('inf')\n"
		code += "best_model_params = None\n"
		code += "for epoch in range(num_epochs):\n"
		code += "    model.train()\n"
		code += "    running_loss = 0.0\n"
		code += "    for inputs, labels in train_loader:\n"
		code += "        optimizer.zero_grad()\n"
		code += "        outputs = model(inputs)\n"
		code += "        loss = criterion(outputs, labels)\n"
		code += "        loss.backward()\n"
		code += "        optimizer.step()\n"
		code += "        running_loss += loss.item()\n"
		code += "    average_loss = running_loss / len(train_loader)\n"
		code += "    print(f'Epoch {epoch + 1}/{num_epochs}, Loss: {average_loss}')\n"
		code += "    if average_loss < best_loss:\n"
		code += "        best_loss = average_loss\n"
		code += "        best_model_params = model.state_dict()\n\n"
		code += "if best_model_params is not None:\n"
		code += "    torch.save(best_model_params, save_params_file_path)\n\n"
		code += "model.load_state_dict(best_model_params)\n"
		code += "model.eval()\n"
		code += "correct = 0\n"
		code += "total = 0\n"
		code += "with torch.no_grad():\n"
		code += "    for inputs, labels in test_loader:\n"
		code += "        outputs = model(inputs)\n"
		code += "        _, predicted = torch.max(outputs, 1)\n"
		code += "        total += labels.size(0)\n"
		code += "        correct += (predicted == labels).sum().item()\n"
		code += "accuracy = correct / total\n"
		code += "print(f'Test Accuracy: {accuracy * 100}%')\n"
	case nodeType == "ReasoningLayer":
		reasoningLayer := node.Layer.(*ReasoningLayer)
		code += "params_file_path = '" + reasoningLayer.ParamsFilePath + "'\n"
		code += "save_result_file_path = '" + reasoningLayer.SaveResultFilePath + "'\n"
		code += "model.load_state_dict(torch.load(params_file_path))\n"
		code += "model.eval()\n"
		code += "with torch.no_grad():\n"
		code += "    reasoning_result = model(reasoning_data)\n"
		//保存到csv
		code += "reasoning_result = reasoning_result.numpy()\n"
		code += "reasoning_result = pd.DataFrame(reasoning_result)\n"
		code += "reasoning_result.to_csv(save_result_file_path, index=False)\n"
	}
	return code
}

func TopologicalSort(graph goraph.Graph) ([]*CNode, bool) {
	nodeIds, ok := goraph.TopologicalSort(graph)
	if !ok {
		return nil, false
	}
	var canvasNodes []*CNode
	for _, node := range nodeIds {
		canvasNodes = append(canvasNodes, graph.GetNode(node).(*CNode))
	}
	return canvasNodes, true
}

func RecoverGraph(canvasContent string) (datasetGraph, netGraph, trainGraph goraph.Graph) {
	datasetGraph = goraph.NewGraph()
	netGraph = goraph.NewGraph()
	trainGraph = goraph.NewGraph()
	var canvas RawCanvas
	if err := json.Unmarshal([]byte(canvasContent), &canvas); err != nil {
		panic(err)
	}
	for key, node := range canvas.Nodes {
		nodeType := node.Type
		switch {
		case util.SliceContains(DatasetKinds, nodeType):
			datasetGraph.AddNode(&canvas.Nodes[key])
		case util.SliceContains(ConvolutionLayerKinds, nodeType) || util.SliceContains(LinearLayerKinds, nodeType) || util.SliceContains(NonlinearActivationKinds, nodeType) || util.SliceContains(PoolingLayerKinds, nodeType):
			netGraph.AddNode(&canvas.Nodes[key])
		case util.SliceContains(OptimizerKinds, nodeType) || util.SliceContains(LossFunctionKinds, nodeType) || util.SliceContains(TrainLayerKinds, nodeType):
			trainGraph.AddNode(&canvas.Nodes[key])
		}
	}
	for _, edge := range canvas.Edges {
		//尝试在三个图中插入边
		err1 := datasetGraph.AddEdge(goraph.StringID(edge.SrcId), goraph.StringID(edge.TgtId), 1)
		err2 := netGraph.AddEdge(goraph.StringID(edge.SrcId), goraph.StringID(edge.TgtId), 1)
		err3 := trainGraph.AddEdge(goraph.StringID(edge.SrcId), goraph.StringID(edge.TgtId), 1)
		if err1 == nil || err2 == nil || err3 == nil {
			srcNode := getNodeById(edge.SrcId, canvas.Nodes)
			tgtNode := getNodeById(edge.TgtId, canvas.Nodes)
			tgtNode.Predecessors = append(tgtNode.Predecessors, srcNode)
		}
	}
	return
}

func getNodeById(nodeId string, nodes []CNode) *CNode {
	for key, node := range nodes {
		if node.Id == nodeId {
			return &nodes[key]
		}
	}
	return nil
}
