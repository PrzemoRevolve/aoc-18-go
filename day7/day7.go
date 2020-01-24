package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type step struct {
	id        byte
	nextSteps []*step
	prevSteps []*step
}

func (s *step) addNext(n *step) {
	s.nextSteps = append(s.nextSteps, n)
}

func (s *step) addPrev(n *step) {
	s.prevSteps = append(s.prevSteps, n)
}

func (s *step) isUnlocked(finished []byte) bool {
	if len(s.prevSteps) == 0 {
		return true
	}

	for _, n := range s.prevSteps {
		var unlocked bool

		for _, b := range finished {
			if n.id == b {
				unlocked = true
			}
		}

		if !unlocked {
			return false
		}
	}

	return true
}

func (s *step) writeSteps() (string, string) {
	var ns strings.Builder
	var ps strings.Builder

	for _, n := range s.nextSteps {
		ns.WriteByte(n.id)
	}
	for _, p := range s.prevSteps {
		ps.WriteByte(p.id)
	}

	return ns.String(), ps.String()
}
func (s *step) duration() int {
	return int(s.id-'A') + 61
}

func newStep(id byte) *step {
	next := make([]*step, 0)
	prev := make([]*step, 0)

	return &step{id, next, prev}
}

type graph map[byte]*step

func (g graph) getFirst() (first *step) {
	for _, s := range g {
		if len(s.prevSteps) == 0 {
			first = s
		}
	}

	return
}

func (g graph) getLowest(nexts []*step, finished []byte) *step {
	var lowest byte = 255
	var next *step

	for _, s := range nexts {
		if s.id < lowest && !contains(finished, s.id) {
			lowest, next = s.id, s
		}
	}

	return next
}

func (g graph) print() {
	var last *step

	for _, s := range g {
		n, p := s.writeSteps()
		fmt.Printf("Step: %s, next: %#v, prev: %#v\n", string(s.id), n, p)

		if len(n) == 0 {
			fmt.Printf("Last Step: %s, next: %#v, prev: %#v\n", string(s.id), n, p)
			if last != nil {
				panic(fmt.Sprint("More that 1 last step found!", string(s.id), string(last.id)))
			}

			last = s
		}
	}
}

func (g graph) getInstructions() string {
	bytes := make([]byte, 0)
	nexts := make([]*step, 0)
	next := g.getFirst()
	nexts = append(nexts, next.nextSteps...)
	bytes = append(bytes, next.id)

	if next == nil {
		panic("No first step!")
	}

	for {
		next = g.getLowest(nexts, bytes)

		if next == nil {
			break
		}

		bytes = append(bytes, next.id)

		for _, s := range next.nextSteps {
			if s.isUnlocked(bytes) && !containsStep(nexts, s) {
				nexts = append(nexts, s)
			}
		}
	}

	if bytes[0] == '0' {
		bytes = bytes[1:]
	}

	return string(bytes)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	f, err := os.Open("./input.txt")
	defer f.Close()
	check(err)

	steps := createGraph(f)
	// steps.print()
	chain := steps.getInstructions()

	fmt.Println("chain:", len(chain), " ", chain)
}

func createGraph(f *os.File) graph {
	f.Seek(0, 0)
	steps := make(graph, 0)
	firsts := make([]*step, 0)

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		text := scanner.Text()

		words := strings.Split(text, " ")
		id := []byte(words[1])[0]
		next := []byte(words[7])[0]

		if id == next {
			panic("Step should not reference itself!")
		}

		currStep := steps[id]

		if currStep == nil {
			currStep = newStep(id)
			steps[id] = currStep
		}

		nextStep := steps[next]

		if nextStep == nil {
			nextStep = newStep(next)
			steps[next] = nextStep
		}

		currStep.addNext(nextStep)
		nextStep.addPrev(currStep)
	}

	for _, s := range steps {
		if len(s.prevSteps) == 0 {
			firsts = append(firsts, s)
		}
	}

	if len(firsts) == 0 {
		panic("No first step!")
	}

	if len(firsts) != 1 {
		zeroStep := newStep('0')

		for _, s := range firsts {
			zeroStep.addNext(s)
			s.addPrev(zeroStep)
		}

		steps[zeroStep.id] = zeroStep
	}

	return steps
}

func contains(arr []byte, c byte) bool {
	for _, item := range arr {
		if item == c {
			return true
		}
	}

	return false
}

func containsStep(arr []*step, c *step) bool {
	for _, item := range arr {
		if item == c {
			return true
		}
	}

	return false
}
