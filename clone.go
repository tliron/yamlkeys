package yamlkeys

func Clone(o interface{}) interface{} {
	switch o.(type) {
	case Map:
		c := make(Map)
		for key, value := range o.(Map) {
			key = KeyData(key)
			MapPut(c, Clone(key), Clone(value))
		}
		return c

	case []interface{}:
		list := o.([]interface{})
		c := make([]interface{}, len(list))
		for index, value := range list {
			c[index] = Clone(value)
		}
		return c

	default:
		return o
	}
}
