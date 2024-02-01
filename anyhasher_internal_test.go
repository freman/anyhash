// Copyright 2024 Shannon Wynter
//
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
package anyhasher

import (
	"bytes"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestKeyValues(t *testing.T) {
	kv := keyValues{
		keyValue{
			name: []byte{'f', 'r', 'e', 'd'},
		},
		keyValue{
			name: []byte{'b', 'l', 'o', 'g', 's'},
		},
		keyValue{
			name: []byte{'f', 'o', 'o'},
		},
		keyValue{
			name: []byte{'b', 'a', 'r'},
		},
	}

	kv.sort()

	require.Equal(t, kv[0].name, []byte{'b', 'a', 'r'})
	require.Equal(t, kv[1].name, []byte{'b', 'l', 'o', 'g', 's'})
	require.Equal(t, kv[2].name, []byte{'f', 'o', 'o'})
	require.Equal(t, kv[3].name, []byte{'f', 'r', 'e', 'd'})
}

func TestSerialise(t *testing.T) {
	obj := struct {
		A string
		D *string
		C int
		B bool
	}{
		A: "Hello",
	}

	var buf bytes.Buffer

	serialise(&buf, reflect.ValueOf(obj))
	require.Equal(t, []byte{'A', 'H', 'e', 'l', 'l', 'o'}, buf.Bytes())
	buf.Reset()

	obj.C = 32

	serialise(&buf, reflect.ValueOf(obj))
	require.Equal(t, []byte{'A', 'H', 'e', 'l', 'l', 'o', 'C', 0x20, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, buf.Bytes())
	buf.Reset()

	tmp := "World"
	obj.B = false
	obj.D = &tmp

	serialise(&buf, reflect.ValueOf(obj))
	require.Equal(t,
		[]byte{
			'A', 'H', 'e', 'l', 'l', 'o', 'C', 0x20, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			'D', 'W', 'o', 'r', 'l', 'd',
		},
		buf.Bytes(),
	)

}

func TestWithTime(t *testing.T) {
	n := time.Now()
	obj := struct {
		A string
		T time.Time
		P *time.Time
		Z string
	}{
		"It's currently",
		n,
		&n,
		"o'clock",
	}

	var buf bytes.Buffer
	serialise(&buf, reflect.ValueOf(&obj))
	res1 := buf.Bytes()

	obj.T = obj.T.Add(time.Second)

	var buf2 bytes.Buffer
	serialise(&buf2, reflect.ValueOf(&obj))
	res2 := buf2.Bytes()

	require.NotEqual(t, res1, res2)
}
