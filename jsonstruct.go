package colly

import "encoding/json"

func ParseJson(resp []byte) ([]map[string]interface{}, error) {
	var retJsonL []map[string]interface{}
	if resp[0] == 123 {
		tmpT := new(map[string]interface{})
		err := json.Unmarshal(resp, tmpT)
		if err != nil {
			return nil, err
		}
		retJsonL = append(retJsonL, *tmpT)
	} else if resp[0] == 91 {
		tmpT := new([]map[string]interface{})
		err := json.Unmarshal(resp, tmpT)
		if err != nil {
			return nil, err
		}
		retJsonL = *tmpT
	} else {
		return nil, ErrUnknownJsonStart
	}
	return retJsonL, nil
}
