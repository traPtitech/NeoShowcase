package util

import "encoding/json"

func IntoField(v any) (map[string]any, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	var f map[string]any
	err = json.Unmarshal(b, &f)
	if err != nil {
		return nil, err
	}
	return f, nil
}
