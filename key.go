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
// YamlKey
//

type YamlKey struct {
	Data interface{}
	Text string
}

func NewYamlKey(data interface{}) (*YamlKey, error) {
	var writer strings.Builder
	encoder := yaml.NewEncoder(&writer)
	encoder.SetIndent(1)
	if err := encoder.Encode(data); err == nil {
		return &YamlKey{
			Data: data,
			Text: writer.String(),
		}, nil
	} else {
		return nil, err
	}
}

// Key interface
func (self *YamlKey) GetKeyData() interface{} {
	return self.Data
}

// fmt.Stringify interface
func (self *YamlKey) String() string {
	return self.Text
}

// yaml.Marshaler interface
func (self *YamlKey) MarshalYAML() (interface{}, error) {
	return self.Data, nil
}
