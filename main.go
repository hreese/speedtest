package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"

	"github.com/seehuhn/mt19937"
)

var (
	localPRNG *rand.Rand
)

// NullReader returns \0
type NullReader bool

func (nr NullReader) Read(p []byte) (n int, err error) {
	numbytes := len(p)
	for idx := 0; idx < numbytes; idx++ {
		p[idx] = 0x00
	}
	return numbytes, nil
}

func init() {
	localPRNG = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func main() {
	// use math/rand non-crypto PRNG
	http.HandleFunc("/1G", func(w http.ResponseWriter, r *http.Request) {
		randreader := io.LimitReader(localPRNG, 1024*1024*1024)
		io.Copy(w, randreader)
	})
	// use Mersenne Twister implementation
	http.HandleFunc("/1GMT", func(w http.ResponseWriter, r *http.Request) {
		rng := rand.New(mt19937.New())
		rng.Seed(time.Now().UnixNano())
		randreader := io.LimitReader(rng, 1024*1024*1024)
		io.Copy(w, randreader)
	})
	// Return \0
	http.HandleFunc("/1GNull", func(w http.ResponseWriter, r *http.Request) {
		randreader := io.LimitReader(NullReader(true), 1024*1024*1024)
		io.Copy(w, randreader)
	})
	fmt.Println("Starting http server")
	http.ListenAndServe(":3333", nil)
}
