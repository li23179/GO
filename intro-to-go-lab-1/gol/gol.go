package main

const dead byte = 0
const live byte = 255

func calculateNextState(p golParams, world [][]byte) [][]byte {

	newWorld := make([][]byte, p.imageHeight)
	for i := range newWorld{
		newWorld[i] = make([]byte, p.imageHeight)
	}

	calculateLiveNbr := func (x, y int) int {
		counter := 0
		relativeDir := [8][2]int{
			{0, 1}, {0, -1},
			{1, 0}, {-1, 0},
			{-1, -1}, {-1, 1},
			{1, 1}, {1, -1},
		}

		for _, dir := range relativeDir{
			ny := (dir[1] + p.imageHeight + y) % p.imageHeight
			nx := (dir[0] + p.imageWidth + x) % p.imageWidth

			if world[ny][nx] == live{
				counter++
			}
		}

		return counter
	}

	for y := 0; y < p.imageHeight; y++{
		for x := 0; x < p.imageWidth; x++{
			counter := calculateLiveNbr(x, y)
			if world[y][x] == live{
				if counter < 2 || counter > 3 {
					newWorld[y][x] = dead
				} else {
					newWorld[y][x] = live
				}
			} else {
				if counter == 3 {
					newWorld[y][x] = live
				}
			}
		}
	}

	return newWorld
}

func calculateAliveCells(p golParams, world [][]byte) []cell {
	cells := []cell{}
	for y := 0; y < p.imageHeight; y++{
		for x := 0; x < p.imageWidth; x++{
			if world[y][x] == live {
				cells = append(cells, cell{x : x, y : y})
			}
		}
	}
	return cells
}
