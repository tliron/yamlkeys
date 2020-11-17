package yamlkeys

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

func KeyData(data interface{}) interface{} {
	if key, ok := data.(Key); ok {
		return key.GetKeyData()
	} else {
		return data
	}
}

func KeyString(data interface{}) string {
	if string_, ok := data.(string); ok {
		return string_
	} else if stringer, ok := data.(fmt.Stringer); ok {
		return stringer.String()
	} else {
		return fmt.Sprintf("%v", data)
	}
}

//
// Key
//

type Key interface {
	GetKeyData() interface{}
}

//
// YAMLKey
//

type YAMLKey struct {
	Data interface{}
	Text string
}

func NewYAMLKey(data interface{}) (*YAMLKey, error) {
	var writer strings.Builder
	encoder := yaml.NewEncoder(&writer)
	encoder.SetIndent(1) // as compact as possible
	if err := encoder.Encode(data); err == nil {
		return &YAMLKey{
			Data: data,
			Text: strings.TrimSuffix(writer.String(), "\n"),
		}, nil
	} else {
		return nil, err
	}
}

// Key interface
func (self *YAMLKey) GetKeyData() interface{} {
	return self.Data
}

// fmt.Stringify interface
func (self *YAMLKey) String() string {
	return self.Text
}

// yaml.Marshaler interface
func (self *YAMLKey) MarshalYAML() (interface{}, error) {
	return self.Data, nil
}

// Utils

func isSimpleKey(data interface{}) bool {
	switch data.(type) {
	case Map, Sequence:
		return false
	}
	return true
}
