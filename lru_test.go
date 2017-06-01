// Copyright 2017 <CompanyName>, Inc. All Rights Reserved.

package hashgen

import (
	"fmt"
	"testing"
)

var testCases = []struct {
	desc     string
	addUuid  string
	getUuid  string
	expectOk bool
}{
	{"Add then Get", "abcd", "abcd", true},
	{"No Add then Get", "abcd", "xyz", false},
	{"Add 2", "a8c37fd2-ade0-4b7b-b62c-528016d73e1b", "a8c37fd2-ade0-4b7b-b62c-528016d73e1b", true},
}

func TestAddAndGet(t *testing.T) {

	cache := NewCache(0, "cachetest.json")

	if cache.Capacity != 1000 {
		t.Fatalf("unexpected capacity")
	}

	for _, testcase := range testCases {
		fmt.Printf("testcase: %s\n", testcase.desc)
		cache.Add(testcase.addUuid, []byte(""), []byte("12345"))
		value, ok := cache.Get(testcase.getUuid)
		// Check type
		/*
		   _, typeok := value(*UserRecord)
		   if !typeok {
		       t.Fatal("Incorrect type")
		   }
		*/
		if ok != testcase.expectOk {
			t.Fatalf("%s expected Get to return %q, got %q\n", testcase.desc, testcase.expectOk, ok)
		}
		if ok && value == nil {
			t.Fatalf("%s failed\n", testcase.desc)
		} else if value != nil {
			fmt.Printf("type value is %T", value)
			fmt.Printf("uuid: %s salt: %s hashbytes: %s\n", value.Uuid, string(value.Salt), string(value.HashBytes))
		}
	}
}
