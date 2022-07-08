package shardmanager

import (
	"fmt"
	"math"
	"testing"

	hash "github.com/theTardigrade/golang-hash"
)

func TestHash(t *testing.T) {
	const input = "b563feb7b2b84b6test7"

	hash := hash.Int8String(input)

	output := hash % 3

	fmt.Println(math.Abs(float64(output)))

}
