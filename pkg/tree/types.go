package tree

import (
	"context"
	"fmt"

	"github.com/gofsd/fsd/pkg/kv"
	"github.com/gofsd/fsd/pkg/store"
	"github.com/gofsd/fsd/pkg/tag"
)

const (
	ROOT uint8 = iota + 1
	TRUNK
	BRANCH
	TWIG
	LEAF
	CROWN
)

type PartPrefix []byte

func (prefix *PartPrefix) Set(part uint8) {
	b := byte(part)
	if len(*prefix) >= 1 && (*prefix)[len(*prefix)-1] != b {
		*prefix = append(*prefix, b)
	} else if len(*prefix) == 0 {
		*prefix = append(*prefix, b)
	}
}

type KV[Key any, Value any] interface {
	Key() Key
	Value() Value
}

type Tag struct {
	ID     uint64
	Parent uint64
	tag.Tag
	Siblings *[]kv.Key
	Childs   *[]kv.Key
	Level    uint
}

type TreeNodes map[uint64]Tag

type Tree struct {
	handler   ForEachFor
	ctx       context.Context
	Tags      TreeNodes
	Args      tag.Tags
	DefKey    uint64
	Acc       kv.Pairs
	init      func(*Tree) error
	destroy   func(*Tree) error
	extract   func(*Tree) error
	transform func(*Tree) error
	load      func(*Tree) error
	store.Stores
	Query store.Queries
	e     error
}

type ForEachFor func(old *Tree, nTags map[uint64]Tag, outerIdx, innerIdx uint64, outer, inner Tag) (ID, parentID kv.Key)
type Handler func(*Tree) error

type Option func(*Tree)

type Get interface {
	Get([]byte) ([]byte, bool)
	GetAll([]byte) ([]store.KV, bool)
	GetAllByPrefix([]byte) ([]store.KV, bool)
	GetRange([]byte) ([]store.KV, bool)
}

type Set interface {
	Set(key, value []byte) any
	Save() error
}

func SetCtx(ctx context.Context) Option {
	return func(tree *Tree) {
		tree.ctx = ctx
	}
}

func SetHandler[Key comparable, Data tag.ITag, Accumulator any](handler ForEachFor) Option {
	return func(tree *Tree) {
		tree.handler = handler
	}
}

func SetArgs(tags tag.Tags) Option {
	return func(tree *Tree) {
		tree.Args = tags
	}
}

func SetInit(handler Handler) Option {
	return func(tree *Tree) {
		tree.init = handler
	}
}

func SetExtract(extract Handler) Option {
	return func(tree *Tree) {
		tree.extract = extract
	}
}

func SetLoad(load Handler) Option {
	return func(tree *Tree) {
		tree.load = load
	}
}

func SetTransform(transform Handler) Option {
	return func(tree *Tree) {
		tree.transform = transform
	}
}

func SetAccumulator(acc kv.Pairs) Option {
	return func(tree *Tree) {
		tree.Acc = acc
	}
}

func SetSlice[Key comparable, Data, Root any](items []Data, prepare func(Data) (uint, Data)) Option {
	return func(tree *Tree) {
		for range items {

		}
	}
}

func Init(t *Tree) (e error) {
	qs := store.Queries{}
	t.Args.ForEachFor("meta", "name", func(tag *tag.Tag, s string) {
		q := store.Query{}
		var bucketPrefix PartPrefix
		var isPersist bool
		tag.ForEach("tree", "root", func(s string) {
			q.BucketName = append(q.BucketName, []byte(s))
			bucketPrefix.Set(ROOT)
			isPersist = true
		})

		tag.ForEach("tree", "trunk", func(s string) {
			q.BucketName = append(q.BucketName, []byte(s))

		})

		tag.ForEach("tree", "branch", func(s string) {
			q.BucketName = append(q.BucketName, []byte(s))

		})

		tag.ForEach("tree", "twig", func(s string) {
			q.BucketName = append(q.BucketName, []byte(s))

		})

		tag.ForEach("tree", "crown", func(s string) {
			q.BucketName = append(q.BucketName, []byte(s))

		})

		tag.ForEach("tree", "name", func(s string) {
			q.FullDatabaseName = s
		})

		if isPersist && len(q.BucketName) == 0 {
			e = fmt.Errorf("Bucket name required for field: %s", s)
		} else {
			if len(bucketPrefix) != 0 {
				q.BucketName = append(q.BucketName, bucketPrefix)
				if v, ok := tag.Get("meta", "def"); ok {
					if typ, ok := tag.Get("meta", "typ"); ok {
						q.Prefix = GetBytesFromStringWithType(v, typ)
					}
				}

				qs = append(qs, q)
			}
		}

	})

	t.Query = qs

	//t.Stores.PrepareStores(store.Queries{})
	return e
}

func GetBytesFromStringWithType(v, typ string) (t []byte) {
	switch typ {
	case "string":
		return []byte(v)
	case "int":
		return kv.GetKeyFromString(v)
	}
	return t
}

func New(setters ...Option) (tree Tree, e error) {
	tree = Tree{}

	for _, setter := range setters {
		setter(&tree)
	}

	if tree.init != nil {
		e = tree.init(&tree)
	}

	if e == nil && tree.extract != nil {
		e = tree.extract(&tree)
	}

	if e == nil && tree.transform != nil {
		e = tree.transform(&tree)
	}

	if e == nil && tree.load != nil {
		e = tree.load(&tree)
	}

	if e == nil && tree.destroy != nil {
		e = tree.destroy(&tree)
	}

	return tree, e
}

func (t *Tree) Exec() (e error) {

	return e
}

func (t *Tree) forEachFor(handler ForEachFor) {
	nTags := TreeNodes{}
	var ID, parentID kv.Key
	for k, v := range t.Tags {
		for key, value := range t.Tags {
			if nTags[key].Parent != t.DefKey {
				break
			}
			ID, parentID = handler(t, nTags, k, key, v, value)
			if ID.Uint64() != t.DefKey && parentID.Uint64() != t.DefKey {
				nTags.UpdateTree(ID, parentID)
				break
			}
		}
	}
	t.Tags = nTags
}

func (t TreeNodes) UpdateTree(ID, parentID kv.Key) {

}

func (t *Tree) Add(data kv.Value, parent kv.Key, level uint) {

}

func (t *Tree) Rm(id kv.Key) {

}

func Merge[t *Tree](first, second t) (r *Tree, err error) {
	return r, err
}

func Reduce[t *Tree](first, second t) (r *Tree, err error) {
	return r, err
}
