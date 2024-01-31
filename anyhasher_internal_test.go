// Copyright 2024 Shannon Wynter
//
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
package anyhasher

import (
	"testing"

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
