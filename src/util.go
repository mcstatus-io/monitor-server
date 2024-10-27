package main

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"os"
	"strconv"
)

// GetInstanceID returns the INSTANCE_ID environment variable parsed as an unsigned 16-bit integer.
func GetInstanceID() uint16 {
	if instanceID := os.Getenv("INSTANCE_ID"); len(instanceID) > 0 {
		value, err := strconv.ParseUint(instanceID, 10, 16)

		if err != nil {
			log.Fatal(err)
		}

		return uint16(value)
	}

	return 0
}

// GetInstanceCount returns the INSTANCE_COUNT environment variable parsed as an unsigned 16-bit integer.
func GetInstanceCount() uint16 {
	if instanceID := os.Getenv("INSTANCE_COUNT"); len(instanceID) > 0 {
		value, err := strconv.ParseUint(instanceID, 10, 16)

		if err != nil {
			log.Fatal(err)
		}

		return uint16(value)
	}

	return 1
}

// RandomHexString generates a random hexadecimal string using the specified byte length.
func RandomHexString(byteLength int) string {
	data := make([]byte, byteLength)

	if _, err := rand.Read(data); err != nil {
		panic(err)
	}

	return hex.EncodeToString(data)
}

// PointerOf returns the pointer of the value.
func PointerOf[T any](v T) *T {
	return &v
}
