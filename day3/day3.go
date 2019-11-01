package main

import (
	"bufio"
	"fmt"
	"os"
)

type claim struct {
	id, left, top, width, height int
	overlaps                     bool
}

type fabric [][]*claim

func newClaim(row string) *claim {
	format := "#%d @ %d,%d: %dx%d"
	var id, left, top, width, height int

	fmt.Sscanf(row, format, &id, &left, &top, &width, &height)

	aClaim := claim{id, left, top, width, height, false}

	return &aClaim
}

func main() {
	f, err := os.Open("./input.txt")
	check(err)
	defer f.Close()

	claims := getClaims(f)
	fab := getFabric(claims)

	markClaims(fab, claims)
	overlaps := countOverlaps(fab)
	original := findOriginal(claims)

	fmt.Println("Overlapping square inches:", overlaps)
	fmt.Println("Not overlapping piece: ", original)
}

func findOriginal(claims []*claim) *claim {
	for _, o := range claims {
		if !o.overlaps {
			return o
		}
	}

	return nil
}

func countOverlaps(fab fabric) (overlaps int) {
	for _, row := range fab {
		for _, cell := range row {
			if cell != nil && cell.overlaps {
				overlaps++
			}
		}
	}

	return
}

func markClaims(fab fabric, claims []*claim) {
	for _, c := range claims {
		for x := c.left; x < c.left+c.width; x++ {
			for y := c.top; y < c.top+c.height; y++ {
				if fab[x][y] != nil {
					fab[x][y].overlaps = true
					c.overlaps = true
				} else {
					fab[x][y] = c
				}
			}
		}
	}
}

func getFabric(claims []*claim) fabric {
	var width, height int

	for _, c := range claims {
		if width < c.left+c.width {
			width = c.left + c.width
		}

		if height < c.top+c.height {
			height = c.top + c.height
		}
	}

	fabric := make(fabric, width)

	for i := range fabric {
		fabric[i] = make([]*claim, height)
	}

	return fabric
}

func getClaims(f *os.File) []*claim {
	f.Seek(0, 0)
	claims := make([]*claim, 0)

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		text := scanner.Text()
		claims = append(claims, newClaim(text))
	}

	return claims
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
