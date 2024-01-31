// Copyright 2024 Shannon Wynter
//
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
package anyhasher_test

import (
	"crypto/md5"
	"encoding/hex"
	"hash/crc32"
	"math"
	"testing"

	"github.com/freman/anyhasher"
	"github.com/stretchr/testify/require"
)

type Version1 struct {
	Name  string
	Value string
}

type Version2 struct {
	Name          string
	SomethingElse string
	Value         string
}

type Version3 struct {
	Name          string
	SomethingElse string
	ANumber       int
	ABool         bool
	Value         string
}

type Recursive1 struct {
	Name  string
	Value string
	Other Version1
	Ptr   *Version2
}

type Recursive2 struct {
	Name  string
	Value string
	Other Version1
	Unset Version1
	Ptr   *Version2
}

type Recursive3 struct {
	Name      string
	Value     string
	Other     Version1
	Unset     Version1
	Ptr       *Version2
	UnsertPtr *Version3
}

type WhatAreYou struct {
	Name              string
	Value             string
	IgnoreMe          string `hash:"-"`
	GettingTricky     map[string]string
	GettingRealTricky map[interface{}]interface{}
	WhyDoThis         []interface{}
	privateField      string
}

func TestHasher(t *testing.T) {
	dat := &Version1{
		Name:  "fred",
		Value: "blogs",
	}

	dat2 := &Version2{
		Name:  "fred",
		Value: "blogs",
	}

	dat3 := &Version3{
		Name:  "fred",
		Value: "blogs",
	}

	dat4 := &Recursive1{
		Name:  "fred",
		Value: "blogs",
		Other: *dat,
		Ptr:   dat2,
	}

	dat5 := &Recursive2{
		Name:  "fred",
		Value: "blogs",
		Other: *dat,
		Ptr:   dat2,
	}

	dat6 := &Recursive3{
		Name:  "fred",
		Value: "blogs",
		Other: *dat,
		Ptr:   dat2,
	}

	dat7 := &WhatAreYou{}

	startHash := anyhasher.SHA512(dat)
	require.Equal(t, startHash, anyhasher.SHA512(dat2), "New zero fields should have no impact on the hash")
	require.Equal(t, startHash, anyhasher.SHA512(dat3), "New zero fields should have no impact on the hash")

	dat3.ABool = true
	require.NotEqual(t, startHash, anyhasher.SHA512(dat3), "Changing a new field should definitely impact the hash")

	complexStartHash := anyhasher.SHA512(dat4)
	require.Equal(t, complexStartHash, anyhasher.SHA512(dat5), "New zero fields should have no impact on even a complicated hash")
	require.Equal(t, complexStartHash, anyhasher.SHA512(dat6), "New zero fields should have no impact on even a complicated hash")

	finalBoss := hex.EncodeToString(anyhasher.SHA512(dat7))

	dat7.GettingTricky = map[string]string{
		"hello":  "world",
		"how":    "goes it",
		"empty!": "",
	}

	require.NotEqual(t, finalBoss, anyhasher.SHA512(dat7), "All changes should be reflected")

	dat7.GettingTricky["hello"] = "Venus"
	require.NotEqual(t, finalBoss, anyhasher.SHA512(dat7), "Even the smallest changes should be reflected")

	dat7.GettingRealTricky = map[interface{}]interface{}{
		"hello":     "world",
		dat:         dat2,
		dat3:        dat4,
		dat5:        dat6,
		"nil!":      nil,
		uint32(123): uint64(10),
	}

	require.NotEqual(t, finalBoss, anyhasher.SHA512(dat7), "All changes should be reflected, even using ptrs as map keys")

	dat7.WhyDoThis = append(dat7.WhyDoThis, dat, dat2, dat3, dat4, dat5)
	lastCheck := anyhasher.SHA512(dat7)
	require.NotEqual(t, finalBoss, lastCheck, "All changes should be reflected, including absurd slices of interface")

	dat7.IgnoreMe = "Hello, Frank Walker from National Tiles"
	require.Equal(t, lastCheck, anyhasher.SHA512(dat7), "Ignored keys should have no impact")

	dat7.privateField = "42: The answer to life, the universe and everything"
	require.Equal(t, lastCheck, anyhasher.SHA512(dat7), "Private keys should have no impact")
}

func TestHasherSimpleTypes(t *testing.T) {
	require.Equal(t,
		[]byte{0x7b, 0x50, 0x2c, 0x3a, 0x1f, 0x48, 0xc8, 0x60, 0x9a, 0xe2, 0x12, 0xcd, 0xfb, 0x63, 0x9d, 0xee, 0x39, 0x67, 0x3f, 0x5e},
		anyhasher.SHA1("Hello world"),
	)

	require.Equal(t,
		[]byte{0xd2, 0xc5, 0x9, 0x49, 0xa6, 0x76, 0x3c, 0xb1, 0x3f, 0x4f, 0xab, 0x89, 0x6f, 0xe6, 0x78, 0x65, 0x48, 0x7, 0x6, 0x7d},
		anyhasher.SHA1(42),
	)

	require.Equal(t,
		[]byte{0x2e, 0xe6, 0xf4, 0x3, 0x2b, 0x1b, 0x6, 0x28, 0xe1, 0xe7, 0x7f, 0x38, 0xa5, 0x1, 0x1e, 0x6d, 0xcb, 0xf0, 0xb2, 0x58},
		anyhasher.SHA1(math.Pi),
	)

	require.Equal(t,
		[]byte{0x76, 0xbb, 0xc2, 0xd8, 0x41, 0x32, 0x85, 0x8e, 0x51, 0xd3, 0x6a, 0xba, 0x3c, 0xc7, 0x51, 0x18, 0x5e, 0x24, 0x64, 0x15, 0xe6, 0x6e, 0xed, 0x8, 0xa9, 0x20, 0xe7, 0xbe, 0x83, 0xf6, 0x34, 0xbd},
		anyhasher.SHA256(map[string]string{"Hello": "world"}),
	)
}

func TestBYOHash(t *testing.T) {
	md5hasher := md5.New()
	anyhasher.HashWith(md5hasher, "Hi there")

	require.Equal(t,
		[]byte{0xd9, 0x38, 0x54, 0x62, 0xd3, 0xde, 0xff, 0x78, 0xc3, 0x52, 0xeb, 0xb3, 0xf9, 0x41, 0xce, 0x12},
		md5hasher.Sum(nil),
	)

	crc32Hasher := crc32.New(crc32.IEEETable)
	anyhasher.HashWith(crc32Hasher, "You never know, it might be useful")

	require.Equal(t,
		[]byte{0xa4, 0x87, 0xf8, 0x1e},
		crc32Hasher.Sum(nil),
	)

}
