package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"
)

type point struct {
	x, y int
	id   string
}

type boardData struct {
	points                 boardArea
	minX, maxX, minY, maxY int
}

type boardArea [][]*point
type boardRow []*point
type result map[string]int

func (b *boardData) coord(x, y int) (int, int) {
	return b.minX + x, b.minY + y
}

func (b *boardData) printBoard() {
	for j, row := range b.points {
		for i, p := range row {
			bx, by := b.coord(i, j)

			if p != nil {
				if p.x == bx && p.y == by {
					fmt.Print(p.id)
				} else {
					fmt.Print(strings.ToLower(p.id))
				}
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println("")
	}
}

func (p *point) distToPoint(o *point) int {
	return p.dist(o.x, o.y)
}
func (p *point) dist(x, y int) int {
	return int(math.Abs(float64(p.x-x)) + math.Abs(float64(p.y-y)))
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

	points := createPoints(f)
	board := createBoard(points)
	res := fillDistances(board, points)
	// fmt.Printf("%#v\n", res)
	// board.printBoard()
	truncInfinites(board, res)

	fmt.Printf("%#v\n", res)
	fmt.Println(findMax(res))
	fmt.Println("Safe Region size: ", safeRegionSize(board, points, 10000))
}

func safeRegionSize(board *boardData, points []*point, safeLimit int) (size int) {
	for j, row := range board.points {
		for i := range row {
			x, y := board.coord(i, j)
			var sum int

			for _, p := range points {
				sum += p.dist(x, y)
			}

			if sum < safeLimit {
				size++
			}
		}
	}

	return
}

func findMax(res result) (id string, max int) {
	for k, v := range res {
		if v > max {
			id, max = k, v
		}
	}

	return
}

func truncInfinites(b *boardData, m result) {
	for _, p := range b.points[0] {
		// fmt.Print(p.id)
		if p != nil {
			m[p.id] = 0
		}
	}
	// fmt.Println("")
	for _, p := range b.points[len(b.points)-1] {
		// fmt.Print(p.id)
		if p != nil {
			m[p.id] = 0
		}
	}
	// fmt.Println("")
	for _, row := range b.points {
		first := row[0]
		last := row[len(row)-1]

		if first != nil {
			m[first.id] = 0
		}
		if last != nil {
			m[last.id] = 0
		}
		// fmt.Println(first.id)
		// fmt.Println(last.id)
	}
}

func fillDistances(board *boardData, points []*point) result {
	m := make(result)
	for j, row := range board.points {
		for i := range row {
			x, y := board.coord(i, j)
			max := len(row)

			for _, p := range points {
				dist := p.dist(x, y)

				if dist == max {
					if row[i] != nil && m[row[i].id] > 0 {
						m[row[i].id]--
					}

					row[i] = nil
					continue
				}

				if dist < max {
					if row[i] != nil {
						m[row[i].id]--
					}

					row[i] = p
					m[p.id]++

					max = dist
				}
			}
		}
	}

	return m
}

func createBoard(points []*point) *boardData {
	var maxX, maxY int
	minX := points[0].x
	minY := points[0].y

	for _, p := range points {
		if p.x > maxX {
			maxX = p.x
		}
		if p.y > maxY {
			maxY = p.y
		}
		if p.x < minX {
			minX = p.x
		}
		if p.y < minY {
			minY = p.y
		}
	}

	height := maxY - minY + 1
	width := maxX - minX + 1
	area := make(boardArea, height)

	for i := range area {
		area[i] = make(boardRow, width)
	}

	return &boardData{area, minX, maxX, minY, maxY}
}

func createPoints(f *os.File) []*point {
	f.Seek(0, 0)
	points := make([]*point, 0)

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	id := 'A'

	for scanner.Scan() {
		var x, y int
		text := scanner.Text()

		fmt.Sscanf(text, "%d, %d", &x, &y)

		points = append(points, &point{x, y, string(id)})
		id++
	}

	return points
}
