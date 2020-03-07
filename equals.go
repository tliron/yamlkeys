package yamlkeys

func Equals(a interface{}, b interface{}) bool {
	switch a_ := a.(type) {
	case Map:
		if bMap, ok := b.(Map); ok {
			// Does A have all the keys that are in B?
			for key := range bMap {
				if _, ok := MapGet(a_, key); !ok {
					return false
				}
			}

			// Are all values in A equal to those in B?
			for key, aValue := range a_ {
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

	case Sequence:
		if bList, ok := b.([]interface{}); ok {
			// Must have same lengths
			if len(a_) != len(bList) {
				return false
			}

			for index, aValue := range a_ {
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
