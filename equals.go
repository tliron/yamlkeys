package yamlkeys

func Equals(a interface{}, b interface{}) bool {
	switch a.(type) {
	case Map:
		if bMap, ok := b.(Map); ok {
			aMap := a.(Map)

			// Does A have all the keys that are in B?
			for key := range bMap {
				if _, ok := MapGet(aMap, key); !ok {
					return false
				}
			}

			// Are all values in A equal to those in B?
			for key, aValue := range aMap {
				if bValue, ok := MapGet(bMap, key); ok {
					if !Equals(aValue, bValue) {
						return false
					}
				} else {
					return false
				}
			}

			return true
		} else {
			return false
		}

	case []interface{}:
		if bList, ok := b.([]interface{}); ok {
			aList := a.([]interface{})

			// Must have same lengths
			if len(aList) != len(bList) {
				return false
			}

			for index, aValue := range aList {
				bValue := bList[index]
				if !Equals(aValue, bValue) {
					return false
				}
			}

			return true
		} else {
			return false
		}

	default:
		return a == b
	}
}
