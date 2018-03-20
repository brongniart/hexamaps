package maps

import . "github.com/brongniart/hexamaps"

func getOffset(coord Coordinate, radius int) (int, int) {
	x := coord.GetX() + radius
	y := coord.GetY()

	offsetx, offsety := -1, -1
	if x < radius+1 {
		offsetx = (x)*radius + (x*(x+1))/2
		offsety = y + x
	} else {
		tmp := radius*radius + radius + (radius*radius+3*radius+2)/2
		x = x - radius - 1
		offsetx = tmp + 2*radius*x - (x*(x-1))/2
		offsety = y + radius
	}

	return offsetx, offsety
}

func getIndices(coord Coordinate, radius int) (int, int) {
	x, y := coord.GetX()+radius, coord.GetY()
	if coord.GetX() < 0 {
		y += coord.GetX() + radius
	} else {
		y += radius
	}

	return x, y
}
