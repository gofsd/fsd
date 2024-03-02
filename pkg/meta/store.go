package meta

import (
	"errors"
	"fmt"

	"github.com/gofsd/fsd/pkg/tag"
)

var FieldDuplication = errors.New("Field duplication")

const NUMBER = "float64"
const STRING = "string"
const BOOLEAN = "bool"
const ARRAY = "[]"
const OBJECT = "map"

type Data map[string]map[int]any
type Meta map[string]tag.Tag

type Store struct {
	Meta Meta
	Data Data
}

// func New() Store {
// 	return Store{
// 		Meta: make(Meta),
// 		Data: make(Data),
// 	}
// }

func (s *Store) AddField(name string, tag tag.Tag) (e error) {
	if _, ok := s.Meta[name]; ok {
		return FieldDuplication
	}
	s.Meta[name] = tag
	s.Data[name] = make(map[int]any)
	return
}

func (m *Meta) Get(name string) (field string, e error) {
	return
}

// func (m Meta) StructFromMeta(path string) interface{} {
// 	splitPath := strings.Split(path, " ")
// 	root := splitPath[0]
// 	var i []any
// 	//var test interface{}
// 	instance := dynamicstruct.NewStruct()

// 	if v, ok := m[path]; ok {
// 		if v.GetStringByTagKey("dynamic", "type") == "int" {
// 			name := v.GetStringByTagKey("dynamic", "name")
// 			//typeN := v.GetStringByTagKey("dynamic", "type")
// 			//fullTags := v.GetTagString()

// 			instance.AddField(
// 				name,
// 				0,
// 				`validate:"gte=2,lte=130"`,
// 			)
// 			// instance.AddField(
// 			// 	"NewTestVar",
// 			// 	0,
// 			// 	`validate:"gte=2,lte=130"`,
// 			// )
// 		}
// 		if v.GetStringByTagKey("dynamic", "type") == "string" {
// 			instance.AddField(
// 				v.GetStringByTagKey("dynamic", "name"),
// 				v.GetStringByTagKey("dynamic", "type"),
// 				v.GetTagString(),
// 			)
// 		}
// 		if strings.Contains(v.GetStringByTagKey("dynamic", "type"), "[]") {
// 			if _, ok := m[path]; ok {
// 				i = append(i, m.StructFromMeta(
// 					fmt.Sprintf("%s.%s", root, strings.ReplaceAll(v.GetStringByTagKey("dynamic", "type"), "[]", "")),
// 				))
// 			}
// 			instance.AddField(
// 				v.GetStringByTagKey("dynamic", "name"),
// 				i,
// 				v.GetTagString(),
// 			)
// 		}

// 	} else {

// 		for k, v := range m {
// 			if strings.Contains(k, path) {
// 				pathLength := len(strings.Split(path, "."))
// 				kLength := len(strings.Split(k, "."))
// 				if pathLength < kLength && kLength < pathLength+2 {
// 					instance.AddField(
// 						v.GetStringByTagKey("dynamic", "name"),
// 						m.StructFromMeta(k),
// 						v.GetTagString(),
// 					)
// 				}

// 			}

// 		}
// 	}

// 	return instance.Build().New()
// }

func (s *Store) Set(field string, data any) (e error) {
	if _, existMeta := s.Meta[field]; existMeta {
		idx := len(s.Data[field]) + 1

		s.Data[field][idx] = data
	}
	return
}

func (d *Data) Remove(field string, id int, data any) (e error) {
	return
}

func (d *Data) Edit(field string, id int, data any) (e error) {
	return
}

func (d *Data) GetByID(fieldName string, id int) (item any, e error) {
	return
}

func typeof(v interface{}) string {
	return fmt.Sprintf("%T", v)
}
