package shuffle

import (
	"math/rand"
	"time"
)

func ShuffleArrUint32(arr []uint32) []uint32 {
	shuffled := make([]uint32, len(arr))
	copy(shuffled, arr)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := len(arr) - 1; i > 0; i-- {
		j := r.Intn(i + 1)
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	}
	return shuffled
}
