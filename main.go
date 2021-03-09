package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"sync"
	"time"

	big "github.com/ncw/gmp"
	// "github.com/cheggaaa/pb/v3"
)

var (
	zero *big.Int = big.NewInt(0)
	one  *big.Int = big.NewInt(1)
	two  *big.Int = big.NewInt(2)

	mu         sync.Mutex
	numWorkers = runtime.NumCPU()
	// numWorkers = 1

	currentP = uint(1)
)

const (
	iterationsPerJob = 2048
)

func main() {
	fmt.Println("Starting to look for primes")
	fmt.Println("numWorkers:", numWorkers)
	fmt.Println(time.Now())

	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(file)

	wg := &sync.WaitGroup{}

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go spawnWorker()
	}

	wg.Wait()
}

func spawnWorker() {
	for {
		mu.Lock()
		initialP := currentP
		currentP += iterationsPerJob
		mu.Unlock()

		// bar := pb.StartNew(iterationsPerJob)

		for i := uint(0); i < iterationsPerJob; i++ {
			myP := initialP + i
			if big.NewInt(int64(myP)).ProbablyPrime(4) && LucasLehmer(myP) {
				fmt.Printf("2^%d-1 is prime\n", myP)
				log.Println("2^", myP, "-1 is prime")
				// fmt.Println(time.Now())
			}
			// bar.Increment()
		}
	}
}

func LucasLehmer(p uint) (isPrime bool) {
	var dummy1, dummy2 big.Int
	s := big.NewInt(4)
	m := big.NewInt(0)
	m = m.Sub(m.Lsh(one, p), one) // = (1 << p) - 1

	for i := 0; i < int(p)-2; i++ {
		s = s.Sub(s.Mul(s, s), two)
		for s.Cmp(m) == 1 {
			s.Add(dummy1.And(s, m), dummy2.Rsh(s, p))
		}
		if s.Cmp(m) == 0 {
			s = zero
		}
	}
	return s.Cmp(zero) == 0
}
