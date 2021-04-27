package yamlkeys

import (
	"fmt"
	"io"
	"strings"

	"gopkg.in/yaml.v3"
)

func Decode(reader io.Reader) (interface{}, error) {
	var node yaml.Node
	decoder := yaml.NewDecoder(reader)
	if err := decoder.Decode(&node); err == nil {
		return DecodeNode(&node)
	} else {
		return nil, WrapWithDecodeError(err)
	}
}

func DecodeAll(reader io.Reader) ([]interface{}, error) {
	decoder := yaml.NewDecoder(reader)
	var values []interface{}
	for {
		var node yaml.Node
		if err := decoder.Decode(&node); err == nil {
			if value, err := DecodeNode(&node); err == nil {
				values = append(values, value)
			} else {
				return nil, err
			}
		} else if err == io.EOF {
			return values, nil
		} else {
			return nil, WrapWithDecodeError(err)
		}
	}
}

func DecodeString(s string) (interface{}, error) {
	return Decode(strings.NewReader(s))
}

func DecodeStringAll(s string) ([]interface{}, error) {
	return DecodeAll(strings.NewReader(s))
}

func DecodeNode(node *yaml.Node) (interface{}, error) {
	switch node.Kind {
	case yaml.AliasNode:
		return DecodeNode(node.Alias)

	case yaml.DocumentNode:
		if len(node.Content) != 1 {
			panic(fmt.Sprintf("malformed YAML @%d,%d: document content count is %d", node.Line, node.Column, len(node.Content)))
		}

		return DecodeNode(node.Content[0])

	case yaml.MappingNode:
		map_ := make(Map)
		mergeMap := make(Map)

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
					switch value_ := value.(type) {
					case Map:
						MapMerge(mergeMap, value_, false)

					case Sequence:
						for _, v := range value_ {
							if m, ok := v.(Map); ok {
								MapMerge(mergeMap, m, false)
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
								return nil, NewDuplicateKeyErrorFor(key, keyNode)
							}
						} else {
							for k := range map_ {
								if Equals(keyData, KeyData(k)) {
									return nil, NewDuplicateKeyErrorFor(key, keyNode)
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

		MapMerge(map_, mergeMap, false)

		return map_, nil

	case yaml.SequenceNode:
		slice := make(Sequence, 0)
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
			return nil, WrapWithDecodeError(err)
		}
	}

	panic(fmt.Sprintf("malformed YAML node: %T", node))
}

func DecodeKeyNode(node *yaml.Node) (interface{}, interface{}, error) {
	if data, err := DecodeNode(node); err == nil {
		if IsSimpleKey(data) {
			return data, nil, nil
		} else {
			key, err := NewYAMLKey(data)
			return key, data, err
		}
	} else {
		return nil, nil, err
	}
}
