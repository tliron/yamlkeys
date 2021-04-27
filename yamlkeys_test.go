package yamlkeys_test

import (
	"testing"

	"github.com/tliron/yamlkeys"
)

var textGood = `
{complex1: 0, complex2: 1}: value1
{complex1: 0, complex2: 2}: value2
`

var textBad = `
{complex1: 0, complex2: 1}: value1
{complex1: 0, complex2: 1}: value2
`

var textPrimitive = `
value1
`

var textMulti = `
{complex1: 0, complex2: 1}: value1
{complex1: 0, complex2: 2}: value2
---
{complex1: 0, complex2: 1}: value3
{complex1: 0, complex2: 2}: value4
`

var key = map[interface{}]interface{}{
	"complex1": 0,
	"complex2": 2,
}

func TestDecode(t *testing.T) {
	if _, err := yamlkeys.DecodeString(textBad); err == nil {
		t.Error("did not find duplicates")
	}
}

func TestIterate(t *testing.T) {
	var data interface{}
	var map_ yamlkeys.Map
	var err error
	var ok bool

	if data, err = yamlkeys.DecodeString(textGood); err != nil {
		t.Error(err)
	}

	if map_, ok = data.(yamlkeys.Map); !ok {
		t.Error("not a map")
	}

	for k := range map_ {
		if _, ok := k.(yamlkeys.Key); !ok {
			t.Errorf("not a Key: %v", k)
		}
	}
}

func TestGet(t *testing.T) {
	var data interface{}
	var map_ yamlkeys.Map
	var err error
	var ok bool

	if data, err = yamlkeys.DecodeString(textGood); err != nil {
		t.Error(err)
	}

	if map_, ok = data.(yamlkeys.Map); !ok {
		t.Error("not a map")
	}

	if value, ok := yamlkeys.MapGet(map_, key); ok {
		if value != "value2" {
			t.Errorf("get returned: %v", value)
		}
	} else {
		t.Error("could not get")
	}
}

func TestPut(t *testing.T) {
	var data interface{}
	var map_ yamlkeys.Map
	var err error
	var ok bool

	if data, err = yamlkeys.DecodeString(textGood); err != nil {
		t.Error(err)
	}

	if map_, ok = data.(yamlkeys.Map); !ok {
		t.Error("not a map")
	}

	yamlkeys.MapPut(map_, key, "value3")

	if value, ok := yamlkeys.MapGet(map_, key); ok {
		if value != "value3" {
			t.Errorf("get returned: %v", value)
		}
	} else {
		t.Error("could not get")
	}
}

func TestPrimitive(t *testing.T) {
	var data interface{}
	var value string
	var err error
	var ok bool

	if data, err = yamlkeys.DecodeString(textPrimitive); err != nil {
		t.Error(err)
	}

	if value, ok = data.(string); !ok {
		t.Error("not a string")
	}

	if value != "value1" {
		t.Errorf("get returned: %v", value)
	}
}

func TestGetMulti(t *testing.T) {
	var data []interface{}
	var map_ yamlkeys.Map
	var err error
	var ok bool

	if data, err = yamlkeys.DecodeStringAll(textMulti); err != nil {
		t.Error(err)
	}

	if len(data) != 2 {
		t.Error("slice length not 2")
	}

	if map_, ok = data[1].(yamlkeys.Map); !ok {
		t.Error("not a map")
	}

	if value, ok := yamlkeys.MapGet(map_, key); ok {
		if value != "value4" {
			t.Errorf("get returned: %v", value)
		}
	} else {
		t.Error("could not get")
	}
}
