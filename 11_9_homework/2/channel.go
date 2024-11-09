package main

import (
	"fmt"
	"math/rand"
	"sort"
	"sync"
)

type Answer struct {
	ID     int
	Random int
}

var Answers []Answer

func main() {
	ch := make(chan Answer, 20)
	var wg sync.WaitGroup

	wg.Add(20)
	for i := 0; i < 20; i++ {
		random := rand.Intn(1145)
		go func(id, random int) {
			defer wg.Done()
			ch <- Answer{ID: i, Random: random}
		}(i, random)
	}
	wg.Wait()
	close(ch)

	for r := range ch {
		Answers = append(Answers, r)
		fmt.Printf("Unsorted Goroutine Index: %d, Random: %d\n", r.ID, r.Random)
	}

	fmt.Println("————————————————————————————————————————————————————————————")

	sort.SliceStable(Answers, func(i, j int) bool {
		return Answers[i].ID < Answers[j].ID
	})

	for _, r := range Answers {
		fmt.Printf("Sorted Goroutine Index: %d, Random: %d\n", r.ID, r.Random)
	}
}
