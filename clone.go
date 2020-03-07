package yamlkeys

func Clone(value interface{}) interface{} {
	switch value_ := value.(type) {
	case Map:
		clone := make(Map)
		for key, value := range value_ {
			key = KeyData(key)
			MapPut(clone, Clone(key), Clone(value))
		}
		return clone

	case Sequence:
		clone := make([]interface{}, len(value_))
		for index, value := range value_ {
			clone[index] = Clone(value)
		}
		return clone

	default:
		return value
	}
}
