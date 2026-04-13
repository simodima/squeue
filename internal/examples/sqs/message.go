package sqs

import (
	"encoding/json"
	"fmt"
)

type MyEvent struct {
	Name string
}

func (e *MyEvent) UnmarshalJSON(data []byte) error {
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	name, ok := raw["name"].(string)
	if !ok {
		return fmt.Errorf("missing or invalid field: name")
	}

	e.Name = name

	return nil
}

func (e *MyEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"name": e.Name,
	})
}
