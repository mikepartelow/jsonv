package jsonv_test

import (
	"encoding/json"
	"mp/jsonv/pkg/jsonv"
	"strings"
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
type ThingV3 struct {
	ThingV2
	V3Key struct {
		V3Sub1Key string `json:"v3sub1key"`
	} `json:"v3key"`
}

const TestThingV2JSON = `{"kind":"thing/v2","v1key":"value1","v2key":"value2"}`
const TestThingV3JSON = `{"kind":"thing/v2","v1key":"value1","v2key":"value2","v3key":{"v3sub1key":"value3"}}`

// TestTheProblem doesn't test any of our code, it just illustrates the problem we solve.
func TestTheProblem(t *testing.T) {
	want := TestThingV2JSON
	// unmarshal into v1, which doesn't have v2key
	var v1 ThingV1
	err := json.Unmarshal([]byte(want), &v1)
	assert.NoError(t, err)
	assert.Equal(t, v1.Kind, "thing/v2")
	assert.Equal(t, v1.V1Key, "value1")

	// marshal from v1
	got, err := json.Marshal(v1)
	assert.NoError(t, err)
	// lost information: new JSON doesn't have v2key anymore, but we want it to
	assert.NotContains(t, got, "v2key")
}

func TestBasic(t *testing.T) {
	given := TestThingV2JSON

	// unmarshal into v1, which doesn't have v2key
	var v1 ThingV1
	var x jsonv.Struct

	err := jsonv.Unmarshal([]byte(given), &v1, &x)
	assert.NoError(t, err)
	assert.Equal(t, v1.Kind, "thing/v2")
	assert.Equal(t, v1.V1Key, "value1")

	wantX := jsonv.Struct{"kind": "thing/v2", "v1key": "value1", "v2key": "value2"}
	assert.Equal(t, wantX, x)

	v1.V1Key = "babbabooie"

	// marshal from v1
	got, err := jsonv.Marshal(v1, x)
	assert.NoError(t, err)
	// no lost information: JSON retains v2key despite no such field in ThingV1
	// it also contains our modification to v1.V1Key
	want := strings.Replace(given, "value1", "babbabooie", 1)
	assert.Equal(t, want, string(got))
}

// func TestFancy(t *testing.T) {
// 	want := TestThingV3JSON

// 	// unmarshal into v1, which doesn't have v2key or v3key
// 	var v1 ThingV1
// 	var x jsonv.Struct

// 	err := jsonv.Unmarshal([]byte(want), &v1, &x)
// 	assert.NoError(t, err)
// 	assert.Equal(t, v1.Kind, "thing/v2")
// 	assert.Equal(t, v1.V1Key, "value1")

// 	wantX := jsonv.Struct{"kind": "thing/v2", "v1key": "value1", "v2key": "value2", "v3key": map[string]interface{}{"v3sub1key": "value3"}}
// 	assert.Equal(t, wantX, x)

// 	// marshal from v1
// 	got, err := jsonv.Marshal(v1, x)
// 	assert.NoError(t, err)
// 	// no lost information: JSON retains v2key despite no such field in ThingV1
// 	assert.Equal(t, want, string(got))

// 	t.Fail() // FIXME: Unmarshal v3 into v3, alter V3Sub1Key, Marshal

// 	// FIXME: arrays, nested arrays, nested arrays with nested structs
// 	// FIXME: ints, floats, etc
// 	// FIXME: unexported fields

// }
