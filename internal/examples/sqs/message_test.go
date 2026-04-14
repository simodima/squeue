package sqs

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMyEvent_UnmarshalJSON_Success(t *testing.T) {
	var e MyEvent
	err := json.Unmarshal([]byte(`{"name":"test-event"}`), &e)
	assert.NoError(t, err)
	assert.Equal(t, "test-event", e.Name)
}

func TestMyEvent_UnmarshalJSON_MissingField(t *testing.T) {
	var e MyEvent
	err := json.Unmarshal([]byte(`{}`), &e)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing or invalid field: name")
}

func TestMyEvent_UnmarshalJSON_WrongType(t *testing.T) {
	var e MyEvent
	err := json.Unmarshal([]byte(`{"name": 42}`), &e)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing or invalid field: name")
}

func TestMyEvent_UnmarshalJSON_InvalidJSON(t *testing.T) {
	var e MyEvent
	err := json.Unmarshal([]byte(`not-json`), &e)
	assert.Error(t, err)
}

func TestMyEvent_MarshalJSON(t *testing.T) {
	e := MyEvent{Name: "test-event"}
	data, err := json.Marshal(&e)
	assert.NoError(t, err)
	assert.JSONEq(t, `{"name":"test-event"}`, string(data))
}
