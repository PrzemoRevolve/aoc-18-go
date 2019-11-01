package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

type state int

const (
	awake  = iota
	asleep = iota
)

type guard struct {
	id        int
	minutes   [60]int
	sumAsleep int
	maxMinute int
}

type timeEntry struct {
	time    time.Time
	guardID int
	state   state
	text    string
}

// byTime implements sort.Interface for []timeEntry based on
// the time field.
type byTime []*timeEntry

func (a byTime) Len() int           { return len(a) }
func (a byTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byTime) Less(i, j int) bool { return a[i].time.Unix() < a[j].time.Unix() }

func newTimeEntry(text string) *timeEntry {
	var dateString, timeString, stateString string
	var guardID int

	split := strings.Split(text, " ")
	stateString = split[len(split)-1]

	fmt.Sscanf(text, "[%10s %5s] Guard #%d", &dateString, &timeString, &guardID)

	timeData, err := time.Parse("2006-01-02 15:04", fmt.Sprintf("%s %s", dateString, timeString))

	// we only care about time after midnight, this piece moves the startpoint for to midnight if is earlier
	if hour, minute, _ := timeData.Clock(); hour == 23 {
		minuteString := fmt.Sprint(60-minute, "m")
		duration, err := time.ParseDuration(minuteString)
		check(err)
		timeData = timeData.Add(duration)
	}

	check(err)
	t := timeEntry{text: text, guardID: guardID, time: timeData}

	switch stateString {
	case "shift":
		t.state = awake
	case "asleep":
		t.state = asleep
	case "up":
		t.state = awake
	default:
		panic(fmt.Sprint("Unrecognized guard state: ", stateString))
	}

	return &t
}

func main() {
	f, err := os.Open("./input.txt")
	check(err)
	defer f.Close()

	timeEntries := getTimeEntries(f)
	sort.Sort(byTime(timeEntries))
	guards := markGuards(timeEntries)

	markMinutes(timeEntries, guards)

	sleepiest := findSleepiestGuard(guards)
	maxMinute := findMaxMinute(sleepiest)
	mostFreqID, mostFreqMin := findMostFrequentMinute(guards)

	fmt.Println("Sleepiest guardID:", sleepiest.id, " Max sleep minute:", maxMinute)
	fmt.Println("Checksum:", sleepiest.id*maxMinute)
	fmt.Println("Most Frequent guardID:", mostFreqID, " most freq minute:", mostFreqMin)
	fmt.Println("Checksum:", mostFreqID*mostFreqMin)
}

func findMostFrequentMinute(guards map[int]*guard) (mostFreqID int, mostFreqMin int) {
	var max int

	for id, g := range guards {
		for m, n := range g.minutes {
			if max < n {
				mostFreqMin, max, mostFreqID = m, n, id
			}
		}
	}

	return
}

func findMaxMinute(g *guard) int {
	var maxI, maxM int

	for i, m := range g.minutes {
		if maxM < m {
			maxI, maxM = i, m
		}
	}

	return maxI
}

func findSleepiestGuard(guards map[int]*guard) *guard {
	var max *guard

	for _, g := range guards {
		if max == nil {
			max = g
		} else if max.sumAsleep < g.sumAsleep {
			max = g
		}
	}

	return max
}

func markMinutes(timeEntries []*timeEntry, guards map[int]*guard) {
	var previous *timeEntry

	for _, current := range timeEntries {
		if previous == nil || previous.state != asleep {
			previous = current
			continue
		}

		currentGuard := guards[previous.guardID]
		startMinute := previous.time.Minute()

		var endMinute int

		if previous.guardID == current.guardID {
			endMinute = current.time.Minute()
		} else {
			endMinute = 60
		}

		for i := startMinute; i < endMinute; i++ {
			currentGuard.minutes[i]++
			currentGuard.sumAsleep++
		}

		previous = current
	}
}

func markGuards(timeEntries []*timeEntry) map[int]*guard {
	var currentGuardID int
	guards := make(map[int]*guard, 0)

	for _, t := range timeEntries {
		if t.guardID != 0 {
			currentGuardID = t.guardID
		} else if currentGuardID != 0 {
			t.guardID = currentGuardID
		} else {
			panic("First time entry without guard ID")
		}

		if guards[t.guardID] == nil {
			var minutes [60]int
			guards[t.guardID] = &guard{id: t.guardID, minutes: minutes}
		}
	}
	return guards
}

func getTimeEntries(f *os.File) []*timeEntry {
	f.Seek(0, 0)
	timeEntries := make([]*timeEntry, 0)

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		text := scanner.Text()
		timeEntries = append(timeEntries, newTimeEntry(text))
	}

	return timeEntries
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
