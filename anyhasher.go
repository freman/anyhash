// Copyright 2024 Shannon Wynter
//
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
package anyhasher

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/binary"
	"hash"
	"io"
	"reflect"
	"sort"
	"strings"
)

// MapKeyHasher is used to hash the contents of map keys which could be structures or other interesting things...
var MapKeyHasher = md5.New

// SHA512 is a simple wrapper to run HashWith with the sha512 hashing algorithm.
func SHA512(v interface{}) []byte {
	h := sha512.New()
	serialise(h, reflect.ValueOf(v))
	return h.Sum(nil)
}

// SHA256 is a simple wrapper to run HashWith with the sha256 hashing algorithm.
func SHA256(v interface{}) []byte {
	h := sha256.New()
	serialise(h, reflect.ValueOf(v))
	return h.Sum(nil)
}

// SHA1 is a simple wrapper to run HashWith with the sha1 hashing algorithm.
func SHA1(v interface{}) []byte {
	h := sha1.New()
	serialise(h, reflect.ValueOf(v))
	return h.Sum(nil)
}

// HashWith takes the given hash object and uses it to hash the given value
func HashWith(h hash.Hash, v interface{}) {
	serialise(h, reflect.ValueOf(v))
}

func serialise(w io.Writer, val reflect.Value) {
	if val.IsZero() {
		return
	}

	switch val.Kind() {
	case reflect.String:
		w.Write([]byte(val.String()))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		binary.Write(w, binary.LittleEndian, val.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		binary.Write(w, binary.LittleEndian, val.Uint())
	case reflect.Float32, reflect.Float64:
		binary.Write(w, binary.LittleEndian, val.Float())
	case reflect.Bool:
		if val.Bool() {
			w.Write([]byte{1})
		}
	case reflect.Ptr:
		if !val.IsNil() || val.Type().Elem().Kind() == reflect.Struct {
			serialise(w, reflect.Indirect(val))
		}
	case reflect.Array, reflect.Slice:
		len := val.Len()
		for i := 0; i < len; i++ {
			serialise(w, val.Index(i))
		}
	case reflect.Map:
		serialiseMap(w, val)
	case reflect.Struct:
		serialiseStruct(w, val)
	case reflect.Interface:
		if !val.CanInterface() {
			return
		}
		serialise(w, reflect.ValueOf(val.Interface()))
	default:
		w.Write([]byte(val.String()))
	}
}

func serialiseMap(w io.Writer, val reflect.Value) {
	mk := val.MapKeys()
	kv := make(keyValues, 0, len(mk))
	for i := range mk {
		v := val.MapIndex(mk[i])
		if v.IsZero() {
			continue
		}

		kh := MapKeyHasher()
		serialise(kh, mk[i])

		kv = append(kv, keyValue{
			name:  kh.Sum(nil),
			value: mk[i],
		})
	}

	kv.sort()

	for i := range kv {
		w.Write(kv[i].name)
		serialise(w, val.MapIndex(kv[i].value))
	}
}

func serialiseStruct(w io.Writer, val reflect.Value) {
	if mbm := val.MethodByName("MarshalBinary"); mbm.IsValid() {
		vr := mbm.Call([]reflect.Value{})

		if !vr[0].IsNil() {
			w.Write(vr[0].Bytes())
			return
		}
	}

	vtype := val.Type()
	flen := vtype.NumField()
	kv := make(keyValues, 0, flen)

	// Get all fields
	for i := 0; i < flen; i++ {
		field := vtype.Field(i)
		if !field.IsExported() {
			continue
		}

		v := val.Field(i)
		if v.IsZero() {
			continue
		}

		str := strings.TrimSpace(field.Tag.Get("hash"))
		if str == "-" {
			continue
		}

		kv = append(kv, keyValue{[]byte(field.Name), v})
	}

	kv.sort()

	for i := range kv {
		w.Write(kv[i].name)
		serialise(w, kv[i].value)
	}
}

type keyValue struct {
	name  []byte
	value reflect.Value
}

type keyValues []keyValue

func (kv *keyValues) sort() {
	sort.Slice(*kv, func(i, j int) bool {
		return bytes.Compare((*kv)[i].name, (*kv)[j].name) < 0
	})
}
