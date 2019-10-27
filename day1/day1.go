package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	f, err := os.Open("./test-input.txt")
	check(err)
	defer f.Close()

	getSum(f)

	findDouble(f)
}

func findDouble(f *os.File) {
	freqs := make(map[int]bool)
	var found bool
	var duplicate, counter, freq int

	for !found {
		counter++
		fmt.Println("Round ", counter)
		found, duplicate = pushIntoMap(&freq, freqs, f)
	}

	fmt.Println("Found duplicate: ", duplicate, found)
}

func pushIntoMap(freq *int, m map[int]bool, f *os.File) (found bool, duplicate int) {
	f.Seek(0, 0)
	s := bufio.NewScanner(f)
	s.Split(bufio.ScanLines)

	for s.Scan() {
		text := s.Text()
		number, err := strconv.Atoi(text)

		check(err)

		*freq += number

		if _, ok := m[*freq]; ok {
			return ok, *freq
		}

		m[*freq] = true
	}

	return false, 0
}

func getSum(f *os.File) {
	f.Seek(0, 0)
	var sum int
	s := bufio.NewScanner(f)
	s.Split(bufio.ScanLines)

	for s.Scan() {
		text := s.Text()
		number, err := strconv.Atoi(text)
		check(err)

		sum += number
	}

	fmt.Println("Sum freq: ", sum)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
