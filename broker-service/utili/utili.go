package utili

import "encoding/json"

func MarshalIndentTab(obj any) ([]byte, error) {
	buf, err := json.MarshalIndent(obj, "", "\t")
	if err != nil {
		return nil, err
	}
	return buf, nil
}
