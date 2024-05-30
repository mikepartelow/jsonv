package jsonv_test

import (
	"encoding/json"
	"mp/jsonv/pkg/jsonv"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Thing struct {
	Kind string `json:"kind"`
}
type ThingV1 struct {
	Thing
	V1Key string `json:"v1key"`
}
type ThingV2 struct {
	ThingV1
	V2Key string `json:"v2key"`
}

// TestTheProblem doesn't test any of our code, it just illustrates the problem we solve.
func TestTheProblem(t *testing.T) {
	want := `
{
	"kind": "thing/v2",
	"v1key": "value1",
	"v2key": "value2"
}`

	// unmarshal into v1, which doesn't have v2key
	var v1 ThingV1
	err := json.Unmarshal([]byte(want), &v1)
	assert.NoError(t, err)
	assert.Equal(t, v1.Kind, "thing/v2")
	assert.Equal(t, v1.V1Key, "value1")

	// marshal from v1
	got, err := json.Marshal(v1)
	assert.NoError(t, err)
	// lost information: new JSON doesn't have v2key anymore
	assert.NotContains(t, got, "v2key")
}

func TestTheSolution(t *testing.T) {
	want := `
{
	"kind": "thing/v2",
	"v1key": "value1",
	"v2key": "value2"
}`

	// unmarshal into v1, which doesn't have v2key
	var v1 ThingV1
	err := jsonv.Unmarshal([]byte(want), &v1)
	assert.NoError(t, err)
	assert.Equal(t, v1.Kind, "thing/v2")
	assert.Equal(t, v1.V1Key, "value1")

	// marshal from v1
	got, err := jsonv.Marshal(v1)
	assert.NoError(t, err)
	// no lost information: JSON retains v2key despite no such field in ThingV1
	assert.Contains(t, got, "v2key")
}
