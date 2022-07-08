package sqs

import (
	"encoding/json"
)

type MyEvent struct {
	Name string
}

func (e *MyEvent) UnmarshalJSON(data []byte) error {
	var raw map[string]any
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}

	e.Name = raw["name"].(string)

	return nil
}

func (e *MyEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"name": e.Name,
	})
}
