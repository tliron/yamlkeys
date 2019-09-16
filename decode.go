package yamlkeys

import (
	"fmt"
	"io"
	"strings"

	"gopkg.in/yaml.v3"
)

func Decode(reader io.Reader) (Map, error) {
	var node yaml.Node
	decoder := yaml.NewDecoder(reader)
	if err := decoder.Decode(&node); err == nil {
		if data, err := DecodeNode(&node); err == nil {
			if map_, ok := data.(Map); ok {
				return map_, nil
			} else {
				return nil, fmt.Errorf("malformed YAML: not a map: %T", data)
			}
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

func DecodeString(s string) (Map, error) {
	return Decode(strings.NewReader(s))
}

func DecodeNode(node *yaml.Node) (interface{}, error) {
	switch node.Kind {
	case yaml.AliasNode:
		return DecodeNode(node.Alias)

	case yaml.DocumentNode:
		var slice []interface{}

		for _, childNode := range node.Content {
			if value, err := DecodeNode(childNode); err == nil {
				slice = append(slice, value)
			} else {
				return nil, err
			}
		}

		switch len(slice) {
		case 1:
			// Single document
			return slice[0], nil
		case 0:
			// Empty
			return make(Map), nil
		default:
			// Multiple documents
			return slice, nil
		}

	case yaml.MappingNode:
		map_ := make(Map)

		// Content is a slice of pairs of key followed by value
		length := len(node.Content)
		if length%2 != 0 {
			panic("malformed YAML map: not a list of key-value pairs")
		}

		for i := 0; i < length; i += 2 {
			keyNode := node.Content[i]
			valueNode := node.Content[i+1]

			if value, err := DecodeNode(valueNode); err == nil {
				if (keyNode.Kind == yaml.ScalarNode) && (keyNode.Tag == "!!merge") {
					// See: https://yaml.org/type/merge.html
					switch value.(type) {
					case Map:
						MapMerge(map_, value.(Map), false)

					case []interface{}:
						for _, v := range value.([]interface{}) {
							if m, ok := v.(Map); ok {
								MapMerge(map_, m, false)
							} else {
								panic(fmt.Sprintf("malformed YAML @%d,%d: merge", node.Line, node.Column))
							}
						}

					default:
						panic(fmt.Sprintf("malformed YAML @%d,%d: merge", node.Line, node.Column))
					}
				} else {
					if key, keyData, err := DecodeKeyNode(keyNode); err == nil {
						// Check for duplicate keys
						if keyData == nil {
							if _, ok := map_[key]; ok {
								return nil, errorDuplicateKey(keyNode, key)
							}
						} else {
							for k := range map_ {
								if Equals(keyData, KeyData(k)) {
									return nil, errorDuplicateKey(keyNode, key)
								}
							}
						}

						map_[key] = value
					} else {
						return nil, err
					}
				}
			} else {
				return nil, err
			}
		}

		return map_, nil

	case yaml.SequenceNode:
		slice := make([]interface{}, 0)
		for _, childNode := range node.Content {
			if value, err := DecodeNode(childNode); err == nil {
				slice = append(slice, value)
			} else {
				return nil, err
			}
		}

		return slice, nil

	case yaml.ScalarNode:
		var value interface{}
		if err := node.Decode(&value); err == nil {
			return value, nil
		} else {
			return nil, err
		}
	}

	panic("malformed YAML node")
}

func DecodeKeyNode(node *yaml.Node) (interface{}, interface{}, error) {
	if data, err := DecodeNode(node); err == nil {
		if isBasicType(data) {
			return data, nil, nil
		} else {
			key, err := NewYamlKey(data)
			return key, data, err
		}
	} else {
		return nil, nil, err
	}
}

// Utils

func isBasicType(data interface{}) bool {
	switch data.(type) {
	case bool, string, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, complex64, complex128:
		return true
	}
	return false
}

func errorDuplicateKey(node *yaml.Node, key interface{}) error {
	return fmt.Errorf("malformed YAML @%d,%d: duplicate map key: %s", node.Line, node.Column, key)
}