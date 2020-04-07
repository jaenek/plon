package main

import (
	"math/rand"
	"time"
)

const dict = "abcdefghijklmnopqrstuwxyzABCDEFGHIJKLMNOPQRSTUWXYZ0123456789"

// Create a new unique id.
func NewUID() string {
	rand.Seed(time.Now().UnixNano())

	var b []byte
	for i := 0; i < 6; i++ {
		b = append(b, dict[rand.Intn(len(dict))])
	}

	return string(b)
}
