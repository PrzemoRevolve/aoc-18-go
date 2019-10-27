package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	f, err := os.Open("./test-input.txt")
	check(err)
	defer f.Close()

	checkSum(f)

	closest := findClosest(f)

	if closest == "" {
		fmt.Println("No closest id found")
	} else {
		fmt.Println("Closest ids match with chars: ", closest)
	}
}

func findClosest(f *os.File) string {
	f.Seek(0, 0)
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	ids := make([]string, 1)

	for scanner.Scan() {
		text := scanner.Text()

		for _, id := range ids {
			diff, same := differentChars(id, text)

			if diff == 1 {
				return same
			}
		}

		ids = append(ids, text)
	}

	return ""
}

func differentChars(a, b string) (uint, string) {
	var count uint
	var same strings.Builder

	for i, run := range a {
		if []rune(b)[i] != run {
			count++
		} else {
			same.WriteRune(run)
		}
	}

	return count, same.String()
}

func checkSum(f *os.File) {
	f.Seek(0, 0)
	var doubles, triples uint

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		text := scanner.Text()
		chars := make(map[string]uint)
		var doubled, tripled bool

		for _, rune := range text {
			char := strconv.QuoteRune(rune)
			chars[char]++
		}

		for _, count := range chars {
			if count == 2 && !doubled {
				// fmt.Println("double: ", char, "in", text)
				doubles++
				doubled = true
			} else if count == 3 && !tripled {
				// fmt.Println("triple: ", char, "in", text)
				triples++
				tripled = true
			}
		}
	}

	fmt.Println("Doubles: ", doubles)
	fmt.Println("Triples: ", triples)
	fmt.Println("Checksum: ", doubles*triples)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
