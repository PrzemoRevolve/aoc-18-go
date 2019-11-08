package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"strings"
)

func main() {
	f, err := ioutil.ReadFile("./input.txt")
	check(err)

	min := stripMers(f)

	fmt.Println("After processing the polymer length is: ", min)
}

func stripMers(input []byte) int {
	min := len(input)
	diff := byte('a') - byte('A')

	for c := byte('A'); c <= byte('Z'); c++ {
		stripped := make([]byte, 0)

		for _, i := range input {
			if i != c && i != c+diff {
				stripped = append(stripped, i)
			}
		}

		res := fullyReactPolymer(stripped)
		newMin := len(strings.TrimSpace(string(res)))

		if newMin < min {
			min = newMin
			fmt.Println(string(c))
		}
	}

	return min
}

func fullyReactPolymer(input []byte) string {
	changed := true

	for changed {
		input, changed = reactPolymer(input)
	}

	return string(input)
}

func reactPolymer(input []byte) ([]byte, bool) {
	var prev byte
	var changed bool
	out := make([]byte, 0)

	for i := 0; i < len(input); i++ {
		if prev == 0 {
			prev = input[i]

			if i == len(input)-1 {
				out = append(out, prev)
			}

			continue
		}

		curr := input[i]

		if math.Abs(float64(curr)-float64(prev)) != 32 {
			out = append(out, prev)
			prev = curr

			if i == len(input)-1 {
				out = append(out, curr)
			}
		} else {
			prev = 0
			changed = true
		}
	}

	return out, changed
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
