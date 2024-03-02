package tag

import (
	"fmt"
	"strings"
)

type KV struct {
	Key   string
	Value string
}

type Tags []Tag

type Tag map[string][]KV

type ITag interface {
	Get(t, k string) (string, bool)
	Set(t, k, v string)
}

func FromString(str string) Tag {
	t := Tag{}
	str = strings.Trim(str, " ")
	tags := strings.Split(str, " ")
	for _, tag := range tags {
		KeyValuePair := strings.Split(tag, ":")
		options := strings.Split(strings.ReplaceAll(KeyValuePair[1], "\"", ""), ",")
		for _, option := range options {
			optionKV := strings.Split(option, "=")

			kv := strings.ReplaceAll(KeyValuePair[0], "\"", "")
			kv = strings.ReplaceAll(KeyValuePair[0], "\\", "")

			optionkv := strings.ReplaceAll(optionKV[0], "\"", "")
			optionkv = strings.ReplaceAll(optionKV[0], "\\", "")
			if t[kv] == nil {
				t[kv] = []KV{}
			}
			if len(optionKV) > 1 {
				t[kv] = append(t[kv], KV{
					Key:   optionkv,
					Value: optionKV[1],
				})
			} else {
				t[kv] = append(t[kv], KV{
					Key:   optionkv,
					Value: "",
				})
			}
		}
	}
	print(str)
	return t
}

func (t Tag) Get(tag, key string) (value string, ok bool) {
	for idx, v := range t[tag] {
		if v.Key == key {
			return t[tag][idx].Value, true
		}
	}

	return "", false
}

func (t Tag) GetAll(tag, key string) (s []string, ok bool) {
	for idx, v := range t[tag] {
		if v.Key == key {
			s = append(s, t[tag][idx].Value)
		}
	}
	if len(s) > 0 {
		return s, true
	} else {
		return nil, false
	}
}

func (t Tag) ForEach(k, v string, callback func(string)) Tag {
	if values, ok := t.GetAll(k, v); ok {
		for _, val := range values {
			callback(val)
		}
	}
	return t
}

func (t *Tag) Set(tag, key, value string) {
	var kv KV
	kv.Key = key
	kv.Value = value
	var exist bool
	for idx, v := range (*t)[tag] {
		if v.Key == kv.Key {
			exist = true
			(*t)[tag][idx] = kv
		}
	}
	if !exist {
		(*t)[tag] = append((*t)[tag], kv)
	}
}

func (tag *Tag) String() string {
	var t, v string
	var tags []string
	for tag, keys := range *tag {
		v = ""
		t = ""
		keyLength := 0
		for _, value := range keys {
			keyLength++
			if value.Value == "" {
				v += fmt.Sprintf("%s", value.Key)
			} else {
				v += fmt.Sprintf("%s=%s", value.Key, value.Value)
			}
			if keyLength < len(keys) {
				v = fmt.Sprintf("%s,", v)
			}
		}
		t = fmt.Sprintf("%s:\"%s\"", tag, v)
		tags = append(tags, t)
	}
	return strings.Join(tags, " ")
}

func New(tag string) Tag {
	return Tag{}
}

func (t *Tags) Append(tag string) {
	newTag := FromString(tag)
	*t = append(*t, newTag)
}

func (t *Tags) FindValue(tag, key, value string) (tg *Tag, ok bool) {
	for _, t := range *t {
		if v, ok := t.Get(tag, key); ok && v == value {
			return &t, true
		}
	}
	return tg, false
}

func (t *Tags) ForEach(callback func(tag *Tag)) {
	for _, v := range *t {
		callback(&v)
	}
}

func (t *Tags) ForEachFor(k, v string, callback func(*Tag, string)) {
	for _, tag := range *t {
		tag.ForEach(k, v, func(s string) {
			callback(&tag, s)
		})
	}
}

func (t *Tags) FindFromList(tag, key string, values []string) {
	for _, v := range values {
		for _, tag := range *t {
			if val, ok := tag.Get("tree", "root"); ok {
				if val == v {

				}
			}
		}
	}
}
