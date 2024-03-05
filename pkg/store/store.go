package store

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/boltdb/bolt"
	typs "github.com/gofsd/fsd-types"
	"github.com/gofsd/fsd/pkg/kv"
)

const (
	DEFAULT_DB_PATH                 = "/tmp"
	DEFAULT_DB_NAME                 = "test.bolt"
	DEFAULT_BUCKET_NAME             = "DEFAUTL"
	RW_FOR_ALL          os.FileMode = 0666
)

type Stores map[string]*Store

func (s *Stores) CloseAll() {
	for _, v := range *s {
		v.Close()
	}
}

func (s Stores) SaveAll() (e error) {
	var tx *bolt.Tx
	var i int
	for _, store := range s {
		if i == 0 {
			tx, e = store.db.Begin(true)
			if e != nil {
				return e
			}
			defer tx.Rollback()
		}
		i++
		store.SaveTx(tx)
	}
	if e == nil {
		tx.Commit()
	}
	return e
}

// Get all buckets and values by prefixes
func (stores *Stores) Read(qs *Queries) (Stores, error) {
	*stores = make(Stores)
	var (
		tx *bolt.Tx
		e  error
		db *bolt.DB
	)
	for i, q := range *qs {
		var opt []Option

		if db != nil {
			opt = append(opt, SetDB(db))
		}
		for _, n := range q.BucketName {
			opt = append(opt, SetBucketName(n))
		}
		opt = append(opt, SetPrefix(q.Prefix))
		if q.FullDatabaseName != "" {
			opt = append(opt, SetFullDbName(q.FullDatabaseName))
		}
		if st, ok := (*stores)[q.FullDatabaseName]; ok {
			opt = append(opt, SetDB(st.db))
		}

		store := New(opt...)
		if store.db != nil {
			if _, ok := (*stores)[store.GetFullDBName()]; !ok {
				tx, e = store.db.Begin(false)
				if e != nil {
					return (*stores), e
				}
				defer tx.Rollback()
				store.GetAllTx(tx, store.Rows)
				(*qs)[i].Rows = store.Rows
				(*stores)[store.GetFullDBName()] = store
			}

		}

	}

	if e == nil {
		tx.Commit()
	}
	return (*stores), e
}

func (stores *Stores) Write(qs *Queries) (Stores, error) {
	var (
		e  error
		tx *bolt.Tx
		db *bolt.DB
	)

	*stores = make(Stores)
	for _, q := range *qs {
		var opt []Option
		if db != nil {
			opt = append(opt, SetDB(db))
		}
		for _, b := range q.BucketName {
			opt = append(opt, SetBucketName(b))
		}
		opt = append(opt, SetPrefix(q.Prefix))
		opt = append(opt, SetFullDbName(q.FullDatabaseName))
		store := New(opt...)
		if _, ok := (*stores)[store.GetFullDBName()]; !ok {
			tx, e = store.db.Begin(true)
			if e != nil {
				return (*stores), e
			}
			defer tx.Rollback()
		}

	}
	if e == nil {
		tx.Commit()
	}
	return *stores, nil
}

type KV struct {
	K []byte
	V []byte
}

type Pair struct {
	ID uint16 `json:"i" validate:"min=1,max=65535"`
	V  string `json:"v" validate:"min=1,max=255"`
}

func (Pair *Pair) ToPretifiedJson() (b []byte, e error) {
	b, e = json.MarshalIndent(Pair.V, "", "  ")

	return b, e
}

func (Pair *Pair) ToJson() (b []byte, e error) {
	b, e = json.Marshal(Pair)
	return b, e
}

func (Pair *Pair) ToString() (s string) {
	return fmt.Sprintf("%d", Pair.ID)
}

func (Pair *Pair) JustUpdate() error {
	return nil
}

func (Pair *Pair) GetValue() any {
	return Pair
}

func (Pair *Pair) FromSlice(k []byte) (e error) {
	e = json.Unmarshal(k, Pair.V)

	return
}

func (Pair *Pair) SetKey(k uint16) {
	Pair.ID = k
}

func (Pair *Pair) SetValue(v string) {
	Pair.V = v
}

func (Pair *Pair) SetID(k uint64) {
	Pair.ID = uint16(k)
}

func (Pair *Pair) Update() error {
	return nil
}

func (Pair *Pair) Read() error {
	return nil
}

func (Pair *Pair) Delete() error {
	return nil
}

func (Pair *Pair) GetKey() []byte {
	return kv.GetKeyFromInt(Pair.ID)
}

type Query struct {
	FullDatabaseName string
	BucketName       [][]byte
	Prefix           []byte
	Rows             *kv.Pairs
}

type Queries []Query

func (qs *Queries) Set(prefix []byte, bucketName [][]byte, dbName string) {
	q := Query{
		Prefix:           prefix,
		BucketName:       bucketName,
		FullDatabaseName: dbName,
	}
	*qs = append(*qs, q)
}

type Store struct {
	fullDbName string
	Rows       *kv.Pairs
	db         *bolt.DB
	keyPrefix  []byte
	bucketName [][]byte
}

type Option func(*Store)

func New(setters ...Option) (store *Store) {
	var e error
	store = &Store{}
	store.Rows = &kv.Pairs{}
	store.SetDefaultOptions()

	for _, setter := range setters {
		setter(store)
	}

	if store.db == nil {
		store.db, e = bolt.Open(store.GetFullDBName(), RW_FOR_ALL, &bolt.Options{
			Timeout: time.Second,
		})
		if e != nil {
			panic(e)
		}
	}

	store.Init()
	return store
}

func (store *Store) GetFullDBName() string {

	return fmt.Sprintf("%s", store.fullDbName)
}

func (store *Store) SetDefaultOptions() {
	store.fullDbName = fmt.Sprintf("%s/%s", DEFAULT_DB_PATH, DEFAULT_DB_NAME)

}
func (db *Store) Init() {
	if len(db.bucketName) == 0 {
		db.bucketName = append(db.bucketName, []byte(DEFAULT_BUCKET_NAME))
	}

	var bucket *bolt.Bucket
	tx, _ := db.db.Begin(true)
	defer tx.Rollback()

	for _, bucketName := range db.bucketName {
		if bucket != nil {
			bucket, _ = bucket.CreateBucketIfNotExists(bucketName)
		} else {
			bucket, _ = tx.CreateBucketIfNotExists(bucketName)
		}

		// bucket.FillPercent = 1.0
	}

	tx.Commit()
}

// name is absolute path to file with file name
func SetFullDbName(name string) Option {
	return func(db *Store) {
		if name != "" {
			db.fullDbName = name
		}
	}
}

func SetBucketName(name []byte) Option {
	return func(db *Store) {
		db.bucketName = append(db.bucketName, name)
	}
}

func SetPrefix(prefix []byte) Option {
	return func(db *Store) {
		db.keyPrefix = prefix
	}
}
func SetDB(db *bolt.DB) Option {
	return func(store *Store) {
		store.db = db
	}
}
func (s *Store) Set(key, value any) *Store {
	s.setKV(key, value)
	return s
}

func (s *Store) GetKeyWithPrefixIfExist(key any) []byte {
	return append(s.keyPrefix, GetBytesFromType(key)...)
}

func (s *Store) Get(k []byte) (value []byte, ok bool) {
	var b *bolt.Bucket
	s.db.View(func(tx *bolt.Tx) error {
		b = s.GetBucket(tx)
		value = b.Get(k)
		return nil
	})
	if value != nil {
		ok = true
	}
	return value, ok

}
func (s *Store) GetBucket(tx *bolt.Tx) (b *bolt.Bucket) {
	for _, bucketName := range s.bucketName {
		if b != nil {
			b = b.Bucket(bucketName)
		} else {
			b = tx.Bucket(bucketName)
		}
		//b.FillPercent = 1.0
	}
	return b
}

func (s *Store) GetString(key string) string {
	var b *bolt.Bucket
	var buf []byte
	k := s.GetKeyWithPrefixIfExist(key)
	s.db.View(func(tx *bolt.Tx) error {
		b = s.GetBucket(tx)

		buf = b.Get(k)
		return nil
	})

	return string(buf)

}

func (s *Store) GetAll(k []byte) (values []KV, ok bool) {
	var c *bolt.Cursor
	k = append(s.keyPrefix, GetBytesFromType(k)...)
	s.db.View(func(tx *bolt.Tx) error {
		c = s.GetBucket(tx).Cursor()

		for key, val := c.Seek(k); key != nil && bytes.HasPrefix(key, k); key, val = c.Next() {
			values = append(values, KV{
				K: key,
				V: val,
			})
		}
		return nil
	})
	if values != nil {
		ok = true
	}
	return values, ok
}

func (s *Store) GetAllTx(tx *bolt.Tx, values *kv.Pairs) (ok bool) {
	var c *bolt.Cursor
	c = s.GetBucket(tx).Cursor()
	var idx int

	for key, val := c.Seek(s.keyPrefix); key != nil && bytes.HasPrefix(key, s.keyPrefix); key, val = c.Next() {
		*values = append(*values, kv.Pair{
			K: kv.GetKeyFromSlice(key),
			V: kv.GetValueFromSlice(val),
		})
		idx++
		if idx == len(*values) {
			break
		}
	}
	if idx > 0 {
		ok = true
	}
	return ok
}

func (s *Store) GetAllByPrefix(k []byte) (values []KV, ok bool) {
	var c *bolt.Cursor
	s.db.View(func(tx *bolt.Tx) error {
		c = s.GetBucket(tx).Cursor()

		for key, val := c.Seek(k); key != nil && bytes.HasPrefix(key, k); key, val = c.Next() {
			values = append(values, KV{
				K: key,
				V: val,
			})
		}
		return nil
	})
	if values != nil {
		ok = true
	}
	return values, ok
}

func (s *Store) GetRange(min, max []byte) (values []KV, ok bool) {
	var c *bolt.Cursor
	s.db.View(func(tx *bolt.Tx) error {
		c = s.GetBucket(tx).Cursor()

		for key, val := c.Seek(min); key != nil && bytes.Compare(key, max) <= 0; key, val = c.Next() {
			values = append(values, KV{
				K: key,
				V: val,
			})
		}
		return nil
	})
	if values != nil {
		ok = true
	}
	return values, ok
}

func (s *Store) Save() (e error) {
	var b *bolt.Bucket
	e = s.db.Batch(func(tx *bolt.Tx) error {
		for _, pair := range *s.Rows {
			b = s.GetBucket(tx)
			e = b.Put(pair.K, pair.V)
		}
		return e
	})
	return e
}

func (s *Store) SaveTx(tx *bolt.Tx) (e error) {
	var b *bolt.Bucket
	for _, pair := range *s.Rows {
		b = s.GetBucket(tx)
		e = b.Put(pair.K[:], pair.V[:])
	}
	return e
}

func (db *Store) SaveTo(bucket any) any {

	return ""
}

func (db *Store) Close() {
	db.db.Close()
}

func (s *Store) setKV(k, v any) {
	var kv kv.Pair
	if len(s.keyPrefix) > 0 {
		kv.K = s.keyPrefix
	}
	kv.K = append(kv.K, GetBytesFromType(k)...)
	kv.V = GetBytesFromType(v)
	s.setRow(kv)
}

func (s *Store) setRow(kv kv.Pair) {
	if i := len(*s.Rows); i < 128 {
		*s.Rows = append(*s.Rows, kv)
	}
}

func GetBytesFromType(a any) []byte {
	switch v := a.(type) {
	case string:
		return []byte(v)
	case []byte:
		return v
	default:
		return []byte("")
	}
}

func (store *Store) JustCreate(s typs.ICrud) (e error) {
	var (
		id uint64
		d  []byte
	)
	e = store.db.Update(func(tx *bolt.Tx) error {
		b := store.GetBucket(tx)
		id, e = b.NextSequence()
		s.SetID(id)
		d, e = s.Json()
		e = b.Put(s.GetKey(), d)
		return e
	})
	return
}

func (store *Store) JustRead(s typs.ICrud) (e error) {
	e = store.db.View(func(tx *bolt.Tx) error {
		b := store.GetBucket(tx)
		val := b.Get(s.GetKey())
		s.FromJson(val)
		return e
	})
	return
}

func (store *Store) JustUpdate(s typs.ICrud) (e error) {
	e = store.db.Update(func(tx *bolt.Tx) error {
		b := store.GetBucket(tx)
		d, e := s.Json()
		e = b.Put(s.GetKey(), d)
		return e
	})
	return
}

func (store *Store) JustDelete(s typs.ICrud) (e error) {
	e = store.db.Update(func(tx *bolt.Tx) error {
		b := store.GetBucket(tx)
		e := b.Delete(s.GetKey())
		return e
	})
	return
}

func (store *Store) JustSet(s typs.ICrud) {
	store.db.Update(func(tx *bolt.Tx) error {
		b := store.GetBucket(tx)
		id, _ := b.NextSequence()
		s.SetID(id)
		d, e := s.Json()
		b.Put(s.GetKey(), d)
		return e
	})
}

var data = map[string]any{}
