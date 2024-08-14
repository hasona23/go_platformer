package tilemap

import (
	"encoding/json"
)

type Object struct {
	Name string  `json:"name"`
	X    float32 `json:"x"`
	Y    float32 `json:"y"`
}
type ObjectLayerJSON struct {
	Objects []Object `json:"objects"`
	Name    string   `json:"name"`
}
type ObjectsJSON struct {
	Layers []ObjectLayerJSON `json:"layers"`
}

func NewObjectLayer(file []byte) (*ObjectsJSON, error) {
	content := file
	var ObjectLayer ObjectsJSON
	var all ObjectsJSON
	err := json.Unmarshal(content, &all)
	if err != nil {
		return nil, err
	}
	for _, i := range all.Layers {
		if len(i.Objects) != 0 {
			ObjectLayer.Layers = append(ObjectLayer.Layers, i)
		}
	}

	return &ObjectLayer, nil
}
func (o *ObjectsJSON) GetLayers() map[string][]Object {
	layerMap := make(map[string][]Object)
	for _, layer := range o.Layers {
		layerMap[layer.Name] = layer.Objects
	}

	return layerMap
}
