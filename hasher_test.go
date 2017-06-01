// Copyright 2017 <CompanyName>, Inc. All Rights Reserved.

package hashgen

import (
    "fmt"
    "encoding/base64"
    "testing"
)

func TestHasherNoSalt(t *testing.T) {

    var testCases = []struct {
        input string
        encodedCryptoHash string
    }{
        {"a", "H0D8ktokFpR1CXnubPWC8tXX0o4YM13gWrxU0FYOD1MChgxlK_CNVgJSql50IQVG82n7u86MEs_HlXsmUv6adQ=="},
        {"benvolio", "4yPu0KYYGSeLgRp0vk3NRplY1Z2vvyAy5zMkl3ZacIXNzRINxhzo5mAd1RHAhrmNtsO8zVDN4bzfbjDzLJ9oAQ=="},
    }

    for _, testcase := range testCases {

        byteBuffer := []byte(testcase.input)
        result := getCryptoHash(byteBuffer, 0)
        encodedResult := base64.URLEncoding.EncodeToString(result)
        fmt.Printf("password: %q encodedCryptoHash: %q", testcase.input, encodedResult)
        if encodedResult != testcase.encodedCryptoHash {
            t.Errorf("Expected %q, got %q", testcase.encodedCryptoHash, encodedResult)
        }
    }
}

/*
func TestHasher(t *testing.T) {

    password := "a"

    byteBuffer := []byte(password)


    //you never store a SHA as a string in a database, but as raw bytes
    //when you want to display a SHA to a user, a common way is Hexadecimal
    //when you want a string representation because it must fit in an URL or in a filename, the usual solution is Base64, which is more compact

    result := getCryptoHash(byteBuffer, 64)

    encodedResult := base64.URLEncoding.EncodeToString(result)

    fmt.Printf("encoded result = %s\n", encodedResult)

    encodedResult2 := base64.URLEncoding.EncodeToString(getCryptoHash(byteBuffer, 64))

    if encodedResult == encodedResult2 {
        t.Fatal()
    }
}
*/
