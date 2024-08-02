package data_source

type DataSource interface {
	GetInterface(keys ...string) interface{}
}

type Stringer interface {
	String() string
}

type MapDataSource map[string]interface{}

func NewMapDataSource(m map[string]interface{}) MapDataSource {
	return m
}

func (s MapDataSource) GetInterface(keys ...string) interface{} {
	// common search
	if len(keys) == 1 {
		for k, v := range s {
			if k == keys[0] {
				return v
			}
		}
	} else if len(keys) > 1 {
		// recursive search
		for k, v := range s {
			if k == keys[0] {
				switch v.(type) {
				case map[string]interface{}:
					return MapDataSource(v.(map[string]interface{})).GetInterface(keys[1:]...)

				default:
				}
			}
		}
	}

	return nil
}
