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
	var map_ yamlkeys.Map
	var err error
	if map_, err = yamlkeys.DecodeString(textGood); err != nil {
		t.Error(err)
	}

	for k, _ := range map_ {
		if _, ok := k.(yamlkeys.Key); !ok {
			t.Errorf("not a Key: %v", k)
		}
	}
}

func TestGet(t *testing.T) {
	var map_ yamlkeys.Map
	var err error
	if map_, err = yamlkeys.DecodeString(textGood); err != nil {
		t.Error(err)
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
	var map_ yamlkeys.Map
	var err error
	if map_, err = yamlkeys.DecodeString(textGood); err != nil {
		t.Error(err)
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
