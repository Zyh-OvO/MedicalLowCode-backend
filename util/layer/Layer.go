package layer

import "encoding/json"

type Layer interface {
	IsLayer()
}

func GenerateLayer(layer Layer, rawData map[string]any) Layer {
	tmpJson, err := json.Marshal(rawData)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(tmpJson, layer)
	if err != nil {
		panic(err)
	}
	return layer
}
