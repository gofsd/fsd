package kv

import (
	"bytes"
	"encoding/ascii85"
	"encoding/binary"
	"strconv"
)

type Key []byte

type Value []byte

type Pair struct {
	K Key
	V Value
}

type Pairs []Pair

type IInt interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64
}

func GetKeyFromSlice(s []byte) Key {
	k := make([]byte, 8)
	copy(k, s)
	return k
}

func GetKeyFromString(s string) Key {
	k := Key{}
	k.SetParsedUint64(s)
	return k
}

func GetKeyFromInt[I IInt](key I) Key {
	k := make([]byte, 8, 8)

	binary.BigEndian.PutUint64(k, uint64(key))
	return k
}

func GetByteFromUint8(key uint8) Key {
	k := []byte{}

	binary.BigEndian.PutUint64(k, uint64(key))
	return k
}

func (k *Key) SetString(s string) {
	copy(*k, s)
}

func (k *Key) SetUint(key uint64) {
	binary.BigEndian.PutUint64(*k, uint64(key))
}

func (k *Key) SetInt(key int) {
	binary.BigEndian.PutUint64(*k, uint64(key))
}

func (k *Key) SetSlice(key []byte) {
	copy(*k, key)
}

func (k *Key) SetArray(key []byte) {
	copy(*k, key[:])
}

func (k Key) String() string {
	dst := make([]byte, 8, 8)

	ascii85.Decode(dst, k[:], false)
	return string(dst)
}

func (k Key) Uint64() uint64 {
	return binary.BigEndian.Uint64(k)
}

func (k Key) Bytes() []byte {
	return bytes.TrimLeft(k[:], "\x00")
}

func (k Key) Slice() []byte {
	return k[:]
}

func (k *Key) SetParsedUint64(s string) {
	i, e := strconv.ParseUint(s, 10, 64)
	if e == nil {
		k.SetUint(i)
	}
}

func GetValueFromString(s string) (v Value) {
	copy(v[:], s)
	return
}

func GetValueFromSlice(s []byte) Value {
	v := make([]byte, len(s))
	copy(v, s)
	return v
}

func (v Value) FromSlice(s []byte) Value {
	copy(v[:], s)
	return v
}

func (v Value) FromString(s []byte) Value {
	copy(v[:], s)
	return v
}

func (v Value) String(s []byte) string {
	dst := make([]byte, 8, 8)

	ascii85.Decode(dst, v[:], false)
	return string(dst)
}

func (k Value) Bytes() []byte {
	return bytes.TrimRight(k, "\x00")
}
