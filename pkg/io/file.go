package io

import (
	"encoding/json"
	"os"
	"strings"

	"encoding/binary"

	"github.com/boltdb/bolt"
	"github.com/gofsd/fsd/pkg/util"
	"github.com/samber/lo"
)

type io interface {
	GetID() int
	Unmarshal([]byte) error
	Marshal() ([]byte, error)
}

type S[T io, R io] struct {
	T    T
	d    []T
	b    []byte
	s    string
	ss   []string
	R    R
	sr   []R
	e    error
	DB   *bolt.DB
	Path string
}
type A []string

func (s *S[T, R]) ReadFile(path string) *S[T, R] {
	s.b, s.e = os.ReadFile(path)
	if s.e != nil {
		util.HandleError(s.e)
	}
	s.s = string(s.b)
	return s
}

func (s *S[T, R]) WriteFile(path string) *S[T, R] {
	s.e = os.WriteFile(path, s.SliceToByteJSON().b, 0644)
	if s.e != nil {
		util.HandleError(s.e)
	}
	return s
}

func (s *S[T, R]) SliceToByteJSON() *S[T, R] {
	s.e = json.Unmarshal(s.b, s.d)
	return s
}

func (s *S[T, R]) Split(seperator string) *S[T, R] {
	s.ss = strings.Split(s.s, seperator)
	return s
}

func (s *S[T, R]) Join(seperator string) *S[T, R] {
	s.s = strings.Join(s.ss, string(seperator))
	return s
}

func (s *S[T, R]) Reduce(reducer func(agg R, item T, idx int) R) *S[T, R] {
	var a R
	s.R = lo.Reduce(s.d, reducer, a)
	return s
}

func (s *S[T, R]) FilterMap(filterMap func(item T, idx int) (R, bool)) *S[T, R] {
	s.sr = lo.FilterMap(s.d, filterMap)
	return s
}

func (s *S[T, R]) ForEach(forEach func(i R, idx int)) *S[T, R] {
	lo.ForEach(s.sr, forEach)
	return s
}

func (s *S[T, R]) ForEachFromSplit(forEach func(i string, idx int) T) *S[T, R] {
	s.d = lo.Map(s.ss, forEach)
	return s
}
func (s *S[T, R]) Put() *S[T, R] {
	err := s.DB.Batch(func(tx *bolt.Tx) error {
		path := strings.Split(s.Path, "/")
		pathLength := len(path)
		var bk, r *bolt.Bucket
		if pathLength == 1 {
			bk, s.e = tx.CreateBucketIfNotExists([]byte(path[0]))
		} else {
			r, s.e = tx.CreateBucketIfNotExists([]byte(path[0]))
			bk, s.e = r.CreateBucketIfNotExists([]byte(path[1]))
		}
		if s.e != nil {
			util.HandleError(s.e)
		}
		s.ForEach(func(i R, idx int) {
			key := convertIntToBytes(i.GetID())
			v := bk.Get(key)
			if v == nil {
				b, _ := i.Marshal()
				bk.Put(key, b)
			}
		})
		return nil
	})
	if err != nil {
		util.HandleError(err)
	}
	return s
}

func (s *S[T, R]) Set() *S[T, R] {
	err := s.DB.Batch(func(tx *bolt.Tx) error {
		path := strings.Split(s.Path, "/")
		pathLength := len(path)
		var bk, r *bolt.Bucket
		if pathLength == 1 {
			bk, s.e = tx.CreateBucketIfNotExists([]byte(path[0]))
		} else {
			r, s.e = tx.CreateBucketIfNotExists([]byte(path[0]))
			bk, s.e = r.CreateBucketIfNotExists([]byte(path[1]))
		}
		s.ForEach(func(i R, idx int) {
			key := convertIntToBytes(i.GetID())
			b, _ := i.Marshal()
			bk.Put(key, b)
		})
		return nil
	})
	if err != nil {
		util.HandleError(err)
	}
	return s
}

func (s *S[T, R]) Select(fillAndFilter func(item T, idx int) (R, bool)) *S[T, R] {
	path := strings.Split(s.Path, "/")
	pathLength := len(path)
	var bk, r *bolt.Bucket

	s.e = s.DB.View(func(tx *bolt.Tx) error {

		if pathLength == 1 {
			bk = tx.Bucket([]byte(path[0]))
		} else {
			r = tx.Bucket([]byte(path[0]))
			bk = r.Bucket([]byte(path[1]))
		}
		if s.e != nil {
			util.HandleError(s.e)
		}
		bk.ForEach(func(k, v []byte) error {
			idx := convertBytesToInt(k)
			t := s.T
			s.e = t.Unmarshal(v)
			if s.e != nil {
				util.HandleError(s.e)
			}
			r, ok := fillAndFilter(t, idx)
			if ok == true {
				s.d = append(s.d, t)
				s.sr = append(s.sr, r)
			}
			return nil
		})
		return nil
	})
	if s.e != nil {
		util.HandleError(s.e)
	}
	return s
}

func (s *S[T, R]) Get(id int) R {
	var bk, r *bolt.Bucket
	path := strings.Split(s.Path, "/")
	pathLength := len(path)
	var item R
	s.DB.View(func(tx *bolt.Tx) error {
		if pathLength == 1 {
			bk, s.e = tx.CreateBucketIfNotExists([]byte(path[0]))
		} else {
			r, s.e = tx.CreateBucketIfNotExists([]byte(path[0]))
			bk, s.e = r.CreateBucketIfNotExists([]byte(path[1]))
		}
		v := bk.Get(convertIntToBytes(id))
		if v != nil {
			item.Unmarshal(v)
		}
		return nil
	})
	return item
}

func (s *S[T, R]) Sort() *S[T, R] {
	return s
}

func (s *S[T, R]) Merge() *S[T, R] {
	return s
}

// itob returns an 8-byte big endian representation of v.
func convertIntToBytes(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

// itob returns an 8-byte big endian representation of v.
func convertBytesToInt(v []byte) int {
	i := binary.BigEndian.Uint64(v)
	return int(i)
}
