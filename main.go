package main

import (
	"fmt"
	"math/big"
	"sort"
)

type Dist map[int]*big.Int

func OnAllRolls(n, sides int, callback func([]int)) {
	rv := make([]int, n)
	if n == 0 {
		callback(rv)
		return
	}
	for i := 1; i <= sides; i++ {
		rv[0] = i
		OnAllRolls(n-1, sides, func(rv2 []int) {
			copy(rv[1:], rv2)
			callback(rv)
		})
	}
}

func Dist4d6RerollOnesOnceDropLowest() Dist {
	rollValue := func(roll []int) int {
		var chosenDice []int
		for i := 0; i < 4; i++ {
			if roll[i] == 1 {
				chosenDice = append(chosenDice, roll[i+4])
			} else {
				chosenDice = append(chosenDice, roll[i])
			}
		}
		sort.Ints(chosenDice)
		var rv int
		for i := 1; i < 4; i++ {
			rv += chosenDice[i]
		}
		return rv
	}

	dist := make(Dist)
	OnAllRolls(4*2, 6, func(rolls []int) {
		val := rollValue(rolls)
		if _, p := dist[val]; !p {
			dist[val] = big.NewInt(1)
		} else {
			dist[val].Add(dist[val], big.NewInt(1))
		}
	})
	return dist
}

func aggOfTwo(a, b Dist, aggfunc func(a, b int) int) Dist {
	c := make(Dist)
	for k1, v1 := range a {
		for k2, v2 := range b {
			aggk := aggfunc(k1, k2)
			acc, p := c[aggk]
			if !p {
				acc = big.NewInt(0)
				c[aggk] = acc
			}
			c[aggk].Add(acc, big.NewInt(0).Mul(v1, v2))
		}
	}
	return c
}

func aggOf(x Dist, n int, aggfunc func(a, b int) int) Dist {
	if n == 1 {
		return x
	}
	return aggOfTwo(x, aggOf(x, n-1, aggfunc), aggfunc)
}

func showDist(prefix string, dist Dist) {
	total := big.NewInt(0)
	cumsum := big.NewInt(0)
	var keys []int
	for k, count := range dist {
		total.Add(total, count)
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, k := range keys {
		count := dist[k]
		cumsum.Add(cumsum, count)

		countF, _ := count.Float64()
		cumsumF, _ := cumsum.Float64()
		totalF, _ := total.Float64()

		probability := countF / totalF
		cumprobability := cumsumF / totalF

		fmt.Printf("%s%d\t%.6f\t%.6f\t%.6f\t%38d\t%38d\n", prefix, k, probability, cumprobability, 1-cumprobability, count, cumsum)
	}
}

func main() {
	d := Dist4d6RerollOnesOnceDropLowest()
	fmt.Println("Sum of scores")
	showDist("  ", aggOf(d, 6, func(a, b int) int { return a + b }))
	fmt.Println("Highest score")
	showDist("  ", aggOf(d, 6, func(a, b int) int {
		if b > a {
			return b
		} else {
			return a
		}
	}))
	fmt.Println("Lowest score")
	showDist("  ", aggOf(d, 6, func(a, b int) int {
		if b < a {
			return b
		} else {
			return a
		}
	}))
}
