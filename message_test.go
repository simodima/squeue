package squeue_test

import "encoding/json"

type TestMessage struct {
	Name string
}

func (e *TestMessage) UnmarshalJSON(data []byte) error {
	var raw map[string]any
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}

	e.Name = raw["name"].(string)

	return nil
}

func (e *TestMessage) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"name": e.Name,
	})
}
