package yamlkeys

type Map = map[interface{}]interface{}

func MapGet(map_ Map, key interface{}) (interface{}, bool) {
	if isBasicType(key) {
		value, ok := map_[key]
		return value, ok
	} else {
		if key_, ok := key.(Key); ok {
			key = key_.GetKeyData()
		}

		for k, value := range map_ {
			if Equals(key, KeyData(k)) {
				return value, true
			}
		}
	}

	return nil, false
}

func MapPut(map_ Map, key interface{}, value interface{}) (interface{}, bool) {
	if isBasicType(key) {
		if existing, ok := map_[key]; ok {
			map_[key] = value
			return existing, true
		}
		map_[key] = value
		return nil, false
	} else {
		if key_, ok := key.(Key); ok {
			keyData := key_.GetKeyData()

			for k, existing := range map_ {
				if Equals(keyData, KeyData(k)) {
					map_[k] = value
					return existing, true
				}
			}

			map_[key] = value
		} else {
			for k, existing := range map_ {
				if Equals(key, KeyData(k)) {
					map_[k] = value
					return existing, true
				}
			}

			var key_ Key
			var err error
			if key_, err = NewYamlKey(key); err == nil {
				map_[key_] = value
			} else {
				panic(err)
			}
		}

		return nil, false
	}
}

func MapDelete(map_ Map, key interface{}) (interface{}, bool) {
	if isBasicType(key) {
		if existing, ok := map_[key]; ok {
			delete(map_, key)
			return existing, true
		}
	} else {
		if key_, ok := key.(Key); ok {
			key = key_.GetKeyData()
		}

		for k, existing := range map_ {
			if Equals(key, KeyData(k)) {
				delete(map_, k)
				return existing, true
			}
		}
	}
	return nil, false
}

func MapMerge(to Map, from Map, override bool) {
	if override {
		for key, value := range from {
			MapPut(to, key, value)
		}
	} else {
		for key, value := range from {
			if key_, ok := key.(Key); ok {
				keyData := key_.GetKeyData()

				exists := false
				for k := range to {
					if Equals(keyData, KeyData(k)) {
						exists = true
						break
					}
				}

				if exists {
					continue
				}

				to[key] = value
			} else {
				if _, ok := to[key]; ok {
					continue
				}

				to[key] = value
			}
		}
	}
}
