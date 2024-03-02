package meta

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gofsd/fsd/pkg/store"
	"github.com/gofsd/fsd/pkg/tag"
	"github.com/gofsd/fsd/pkg/tree"
	dynamicstruct "github.com/ompluscator/dynamic-struct"
)

type Options struct {
	tags  tag.Tags
	tree  tree.Tree
	typ   reflect.Type
	val   reflect.Value
	value any
	err   error
}
type Option func(*Options)
type MetaData struct {
	args    any
	returns any
	e       error
}

func getProp(d interface{}, label string) (interface{}, bool) {
	switch reflect.TypeOf(d).Kind() {
	case reflect.Struct:
		v := reflect.ValueOf(d).FieldByName(label)
		return v.Interface(), true
	case reflect.Ptr:
		var v reflect.Value = reflect.ValueOf(d)
		canSet := v.CanSet()
		isNil := v.IsNil()
		if isNil && canSet {
			v.Set(reflect.New(v.Type().Elem()))

		}
		v = reflect.Indirect(v).FieldByName(label)

		return v.Interface(), true
	}

	return nil, false
}

func (r MetaData) ToJson() (b []byte, e error) {
	return
}

func getField(v any, field string) string {
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByName(field)
	return f.String()
}

func (r Options) GetTypByPath(path string) (prop any, ok bool) {
	fullPath := strings.Split(path, ".")
	prop = r.value
	for _, p := range fullPath {
		prop, ok = getProp(prop, p)
		if !ok {
			return prop, ok
		}
	}
	return prop, ok
}

func (r Options) GetString(path string) (string, bool) {
	if value, ok := r.GetTypByPath(path); ok {
		if v, ok := value.(string); ok {
			return v, ok
		}
	}
	return "", false
}
func Field(name, typ, tg string) Option {
	t := tag.FromString(fmt.Sprintf(`meta:"name=%s,typ=%s" %s`, name, typ, tg))
	return func(o *Options) {
		o.tags = append(o.tags, t)
	}
}

func Set(t tag.Tag) Option {
	return func(o *Options) {
		o.tags = append(o.tags, t)
	}
}

func Settags(tags ...tag.Tag) Option {
	return func(o *Options) {
		o.tags = append(o.tags, tags...)
	}
}

func GetStringFromAny(v any) (s, t string) {
	switch v.(type) {
	case int:
		s, t = fmt.Sprintf("%d", v), "int"
		return
	case float64:
		s, t = fmt.Sprintf("%f", v), "float"
		return
	case string:
		s, t = v.(string), "string"
		return
	case bool:
		s, t = fmt.Sprintf("%t", v), "bool"
		return
	case []any:
		s = "array"
		return
	case struct{}:
		s = "struct"
		return

	}
	return
}

func GetDefaultValueByTypeName(typ string) any {
	if typ == "int" {
		return 0
	} else if typ == "string" {
		return ""
	} else if typ == "bool" {
		return false
	} else if typ == "float" {
		return 0.0
	} else if typ == "array" {
		return []any{}
	} else if typ == "struct" {
		return "struct"
	}
	return nil
}

func OptionsToStruct(o *Options, parent string) any {
	instance := dynamicstruct.NewStruct()
	parents := strings.Split(parent, ".")
	for _, v := range o.tags {
		if name, ok := v.Get("meta", "name"); ok {
			names := strings.Split(name, ".")
			var canContinue bool
			if len(parents)+1 == len(names) {
				for i, _ := range parents {
					if parents[i] == names[i] {
						canContinue = true
					} else {
						canContinue = false
						break
					}
				}

			} else if parent == name || len(parents) == len(names) && parent != "" {
				continue
			} else if len(parents) > len(names) || len(parents)+1 < len(names) {
				continue
			} else {
				canContinue = true
			}
			if !canContinue {
				continue
			}
			if typ, ok := v.Get("meta", "typ"); ok {
				if nilVal := GetDefaultValueByTypeName(typ); nilVal == nil {
					continue
				}
				n := strings.Title(names[len(names)-1])
				if typ == "struct" {
					obj := OptionsToStruct(o, name)
					instance.AddField(
						n,
						obj,
						v.String(),
					)
				} else {
					instance.AddField(
						n,
						GetDefaultValueByTypeName(typ),
						v.String(),
					)
				}

			}
		}

	}

	return instance.Build().New()

}

func OptionsToStructWIthMapsAndArrays(s any) any {
	return s
}

func Neww(setters ...Option) (meta MetaData, err error) {
	args := &Options{}

	for _, setter := range setters {
		setter(args)
	}
	dynamicStruct := OptionsToStruct(args, "")

	OptionsToStructWIthMapsAndArrays(dynamicStruct)
	return meta, nil
}

func InitializeStruct(t reflect.Type, v reflect.Value) {

	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		ft := t.Field(i)
		switch ft.Type.Kind() {
		case reflect.Map:
			f.Set(reflect.MakeMap(ft.Type))
		case reflect.Slice:
			f.Set(reflect.MakeSlice(ft.Type, 0, 0))
		case reflect.Chan:
			f.Set(reflect.MakeChan(ft.Type, 0))
		case reflect.Struct:
			InitializeStruct(ft.Type, f)
		case reflect.Ptr:
			fv := reflect.New(ft.Type.Elem())
			InitializeStruct(ft.Type.Elem(), fv.Elem())
			f.Set(fv)
		case reflect.String:
			ft.Tag.Get("meta")
			f.SetString("Hello")
		case reflect.Int:
			//f.SetInt()
		default:
		}
	}
}

func Handler[Key comparable, Data any, Root any](old *tree.Tree, ntags map[Key]tree.Tag, outerIdx, innerIdx Key, outer, inner tree.Tag) (ID, parentID Key) {
	return ID, parentID
}

func Extract(t *tree.Tree) (e error) {
	t.Stores.ExecQueries(t.Query)
	return e
}

func Load(t *tree.Tree) (e error) {
	t.Stores.ExecQueries(store.Queries{})
	return e
}

func New(setters ...Option) (meta MetaData) {

	options := &Options{}

	for _, setter := range setters {
		setter(options)
	}

	if options.err != nil {
		meta.e = options.err
		return meta
	}

	meta.args, meta.e = tree.New(
		tree.SetArgs(options.tags),
		tree.SetInit(tree.Init),
		tree.SetExtract(Extract),
		tree.SetLoad(Load),
	)

	options.typ = OptionsToNewStruct(options, "")
	options.val = reflect.New(options.typ).Elem()

	InitializeStruct(options.typ, options.val)
	options.value = options.val.Addr().Interface()
	meta.returns = options.value
	return meta
}

func (m *Options) Value() any {
	return m.value
}

func (m *Options) Validate() any {
	return m.value
}

func Op() *Options {
	var opt Options
	return &opt
}

func (opt *Options) Set() Option {
	return func(o *Options) {
		o.tags = append(o.tags, opt.tags...)
	}
}

// Set name and default value
func (opt *Options) ReturnAs(name string) *Options {
	t := opt.tags[len(opt.tags)-1]
	t.Set("return", "value", name)
	return opt
}

// Set name and default value
func (opt *Options) ReturnKeyAs(name string) *Options {
	t := opt.tags[len(opt.tags)-1]
	t.Set("return", "key", name)
	return opt
}

// Set name and default value
func (opt *Options) Leaf(name ...any) *Options {
	opt.tags = append(opt.tags, tag.New(`tree:""`))

	return opt
}

// Set name and default value
func (opt *Options) Crown(names ...string) *Options {
	for _, name := range names {
		if t, ok := opt.tags.FindValue("meta", "name", name); ok {
			t.Set("tree", "crown", name)
		}
	}

	return opt
}

func (opt *Options) AddAction(action string, names []string) *Options {
	if opt.err != nil {
		return opt
	}
	for _, n := range names {
		if v, ok := opt.tags.FindValue("meta", "name", n); ok {
			v.Set("action", "type", action)
		} else {
			opt.err = fmt.Errorf("Field with name: %s does'nt exist", n)
		}
	}

	return opt
}

// Set name and default value
func (opt *Options) Create(names ...string) *Options {
	opt.AddAction("create", names)
	return opt
}

// Set name and default value
func (opt *Options) Read(names ...string) *Options {
	opt.AddAction("read", names)
	return opt
}

// Set name and default value
func (opt *Options) Update(names ...string) *Options {
	opt.AddAction("update", names)
	return opt
}

// Set name and default value
func (opt *Options) Delete(names ...string) *Options {
	opt.AddAction("delete", names)
	return opt
}

// Set name and default value
func (opt *Options) Root(root string, names ...string) *Options {
	for _, name := range names {
		if t, ok := opt.tags.FindValue("meta", "name", name); ok {
			t.Set("tree", "root", root)
		}
	}

	return opt
}

// Get leaf
func (opt *Options) Get(name string) *Options {
	opt.tags = append(opt.tags, tag.New(""))

	return opt
}

// Set name and default value
func (opt *Options) ND(name string, def any) *Options {
	opt.tags = append(opt.tags, tag.New(""))

	return opt
}

// Set - field name, default value, validation rules
func (opt *Options) NDV(name string, def any, validate string) *Options {
	t := tag.New("")
	t.Set("meta", "name", name)
	v, typ := GetStringFromAny(def)

	t.Set("meta", "def", v)
	t.Set("meta", "typ", typ)
	t.Set("meta", "validate", validate)

	opt.tags = append(opt.tags, t)

	return opt
}

// Set - field name, default value, parent field name, validation rules
func (opt *Options) NDPV(name string, def any, parent, validate string) {

}

func OptionsToNewStruct(o *Options, parent string) reflect.Type {
	instance := []reflect.StructField{}
	parents := strings.Split(parent, ".")
	for _, tag := range o.tags {
		if name, ok := tag.Get("meta", "name"); ok {
			names := strings.Split(name, ".")
			var canContinue bool
			if len(parents)+1 == len(names) {
				for i, _ := range parents {
					if parents[i] == names[i] {
						canContinue = true
					} else {
						canContinue = false
						break
					}
				}

			} else if parent == name || len(parents) == len(names) && parent != "" {
				continue
			} else if len(parents) > len(names) || len(parents)+1 < len(names) {
				continue
			} else {
				canContinue = true
			}
			if !canContinue {
				continue
			}
			if typ, ok := tag.Get("meta", "typ"); ok {
				if nilVal := GetDefaultValueByTypeName(typ); nilVal == nil {
					continue
				}
				n := strings.Title(names[len(names)-1])
				if typ == "struct" {
					obj := OptionsToNewStruct(o, name)
					instance = append(instance, reflect.StructField{
						Name: n,
						Type: obj,
						Tag:  reflect.StructTag(tag.String()),
					})
				} else {
					switch v := GetDefaultValueByTypeName(typ).(type) {
					case string:
						instance = append(instance, reflect.StructField{
							Name: n,
							Type: reflect.TypeOf(v),
							Tag:  reflect.StructTag(tag.String()),
						})
					}
				}

			}
		}

	}

	return reflect.StructOf(instance)

}

func If[T any](b bool, first, second T) T {
	if b == true {
		return first
	} else {
		return second
	}
}
