package main

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMyMessage_UnmarshalJSON_Success(t *testing.T) {
	var m myMessage
	err := json.Unmarshal([]byte(`{"name":"hello"}`), &m)
	assert.NoError(t, err)
	assert.Equal(t, "hello", m.name)
}

func TestMyMessage_UnmarshalJSON_MissingField(t *testing.T) {
	var m myMessage
	err := json.Unmarshal([]byte(`{}`), &m)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing or invalid field: name")
}

func TestMyMessage_UnmarshalJSON_WrongType(t *testing.T) {
	var m myMessage
	err := json.Unmarshal([]byte(`{"name": 42}`), &m)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing or invalid field: name")
}

func TestMyMessage_UnmarshalJSON_InvalidJSON(t *testing.T) {
	var m myMessage
	err := json.Unmarshal([]byte(`not-json`), &m)
	assert.Error(t, err)
}

func TestMyMessage_MarshalJSON(t *testing.T) {
	m := myMessage{name: "hello"}
	data, err := json.Marshal(&m)
	assert.NoError(t, err)
	assert.JSONEq(t, `{"name":"hello"}`, string(data))
}
