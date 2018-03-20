package hexamaps

type Coordinate struct {
	x, y, z int
}

func (this *Coordinate) SetX(x int) {
	this.x = x
	this.z = -this.x - this.y
}
func (this *Coordinate) SetY(y int) {
	this.y = y
	this.z = -this.x - this.y
}
func (this *Coordinate) GetX() int {
	return this.x
}
func (this *Coordinate) GetY() int {
	return this.y
}
func (this *Coordinate) GetZ() int {
	return this.z
}

func (this Coordinate) Distance(destination Coordinate) int {
	return max(abs(this.x-destination.x), max(abs(this.y-destination.y), abs(this.z-destination.z)))
}

func (this Coordinate) Neighboors(radius int) func() (bool, Coordinate) {

	dx := -radius - 1
	dy := min(radius, -dx+radius)
	return func() (bool, Coordinate) {
		if dy == min(radius, -dx+radius) {
			if dx == radius {
				return true, Coordinate{-1, -1, -1}
			} else {
				dx++
				dy = max(-radius, -dx-radius)
			}
		} else {
			dy++
		}
		return false, Coordinate{x: dx + this.x, y: dy + this.y,
			z: -(dx + this.x) - (dy + this.y)}
	}
}

func abs(num int) int {
	if num > 0 {
		return num
	} else {
		return -num
	}
}

func max(a int, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

func min(a int, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

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
