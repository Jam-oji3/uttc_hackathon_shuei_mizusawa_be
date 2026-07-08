package util

import (
	"log"
	"math/rand"
	"time"

	"github.com/oklog/ulid/v2"
)

var entropy *ulid.MonotonicEntropy

func init() {
	entropy = ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)
}

func GenerateULID() string {
	id, err := ulid.New(ulid.Timestamp(time.Now()), entropy)
	if err != nil {
		log.Fatalf("failed to generate ULID: %v", err)
	}
	return id.String()
}
