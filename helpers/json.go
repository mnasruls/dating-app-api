package helpers

import (
	"bytes"
	"encoding/json"
)

func Unmarshal(data any, res interface{}) error {
	byted, err := json.Marshal(&data)
	if err != nil {
		return err
	}
	return json.Unmarshal(byted, &res)
}

func JsonMarsalIndent(data interface{}) string {
	jsons, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return ""
	}

	return string(jsons)
}

func JsonMinify(obj interface{}) string {
	jsons, _ := json.Marshal(obj)
	dst := &bytes.Buffer{}
	err := json.Compact(dst, []byte(jsons))
	if err != nil {
		return ""
	}
	return dst.String()
}
