package utils

import (
	"encoding/json"
	"strconv"
)

func JsonString2map(str string) (jmap map[string]interface{}) {
	if err := json.Unmarshal([]byte(str), &jmap); err == nil {

	}
	return
}

func Interface2Json(data interface{}, hasIndent bool, coding bool)([]byte,error){
	var content []byte
	var err error
	if hasIndent {
		content, err = json.MarshalIndent(data, "", "  ")
	} else {
		content, err = json.Marshal(data)
	}
	if err != nil {
		return content,err
	}
	if coding {
		content = []byte(stringsToJson(string(content)))
	}
	return content,err
}

func stringsToJson(str string) string {
	rs := []rune(str)
	jsons := ""
	for _, r := range rs {
		rint := int(r)
		if rint < 128 {
			jsons += string(r)
		} else {
			jsons += "\\u" + strconv.FormatInt(int64(rint), 16) // json
		}
	}
	return jsons
}