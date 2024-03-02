package meta

import (
	"encoding/json"
	"fmt"

	"github.com/gofsd/fsd/pkg/tag"
	"github.com/tidwall/gjson"
)

type Json string

func IsJSON(str string) (js json.RawMessage, isJson bool) {
	isJson = json.Unmarshal([]byte(str), &js) == nil
	return js, isJson
}

func JsonToMeta(json string, meta Meta, root string) {
	gjson.Parse(json).ForEach(func(key, value gjson.Result) bool {
		newRoot := fmt.Sprintf("%s.%s", root, key.Str)
		if value.Type == gjson.String ||
			value.Type == gjson.Number ||
			value.Type == gjson.False ||
			value.Type == gjson.True ||
			value.Type == gjson.Null {

			meta[newRoot] = tag.FromString(value.Str)
		} else if value.Type == gjson.JSON {
			if value.IsObject() {
				JsonToMeta(value.Raw, meta, newRoot)
			} else if value.IsArray() {
				//JsonToMeta(value.Raw, meta, newRoot)
			}
		}
		return true
	})
}
