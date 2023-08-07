package encoding

import (
	"encoding/json"

	"gopkg.in/yaml.v2"
)

type Marshaler interface {
	Marshal() any
}

type GrandMarshaler interface {
	json.Marshaler
	yaml.Marshaler
}

type grandMarshaler struct {
	Marshaler
}

func NewGrandMarshaler(m Marshaler) (gm GrandMarshaler) {
	gm = grandMarshaler{
		Marshaler: m,
	}
	return
}

func (gm grandMarshaler) MarshalJSON() ([]byte, error) {
	return json.Marshal(gm.Marshal())
}

func (gm grandMarshaler) MarshalYAML() (any, error) {
	return gm.Marshal(), nil
}
