package colly

import "encoding/json"

// ParseJSON is the function to auto check and parse Json data
func ParseJSON(resp []byte) ([]map[string]interface{}, error) {
	var retJSONL []map[string]interface{}
	if resp[0] == 123 {
		tmpT := new(map[string]interface{})
		err := json.Unmarshal(resp, tmpT)
		if err != nil {
			return nil, err
		}
		retJSONL = append(retJSONL, *tmpT)
	} else if resp[0] == 91 {
		tmpT := new([]map[string]interface{})
		err := json.Unmarshal(resp, tmpT)
		if err != nil {
			return nil, err
		}
		retJSONL = *tmpT
	} else {
		return nil, ErrUnknownJSONStart
	}
	return retJSONL, nil
}
