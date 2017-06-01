// Copyright 2017 <CompanyName>, Inc. All Rights Reserved.

package hashgen

import (
	"crypto/rand"
	"crypto/sha512"
	"log"
)

// Recommended randomly generated salt is at least as long as digest
const DefaultDigestSize = 64

// Generates a random salt to be added to the password before hashing
func getSalt(length int) []byte {

	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	return b
}

// When digest size is zero, does not add salt
func getCryptoHash(byteBuffer []byte, digestSize int) []byte {

	hasher := sha512.New()
	hasher.Write(byteBuffer)

	if digestSize > 0 {
		hasher.Write(getSalt(digestSize))
	}
	result := hasher.Sum(nil)
	return result
}
