package encoding

import "encoding/json"

func JSONUnmarshalFunc(data []byte) func(any) error {
	return func(obj any) error {
		return json.Unmarshal(data, &obj)
	}
}
