package postgresql

import "encoding/json"

// KeyValsDiff get difference between two entities
func KeyValsDiff(old, new interface{}) (map[string]interface{}, error) {
	var mold map[string]interface{}
	var mnew map[string]interface{}

	b, err := json.Marshal(old)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(b, &old); err != nil {
		return nil, err
	}

	b, err = json.Marshal(new)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(b, &old); err != nil {
		return nil, err
	}

	var keyVals map[string]interface{}
	for key, val := range mnew {
		if mold[key] != mnew[key] {
			keyVals[key] = val
		}
	}

	return keyVals, nil
}
