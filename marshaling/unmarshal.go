package marshaling

import (
	"encoding/json"

	"gopkg.in/yaml.v2"
)

type Unmarshaler interface {
	Unmarshal(unmarshal func(any) error) error
}

type GrandUnmarshaler interface {
	json.Unmarshaler
	yaml.Unmarshaler
}

type grandUnmarshaler struct {
	Unmarshaler
}

func NewGrandUnmarshaler(u Unmarshaler) (gu GrandUnmarshaler) {
	gu = &grandUnmarshaler{
		Unmarshaler: u,
	}
	return
}

func (gu *grandUnmarshaler) UnmarshalJSON(data []byte) error {
	return gu.Unmarshal(JSONUnmarshalFunc(data))
}

func (gu *grandUnmarshaler) UnmarshalYAML(unmarshal func(any) error) error {
	return gu.Unmarshal(unmarshal)
}
