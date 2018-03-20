package hexamaps

import (
	"errors"
)

type InMemoryMap struct {
	Radius     uint16
	properties [][]*Property
}

func (this *InMemoryMap) get(coord Coordinate) (*Property, error) {
	if coord.Distance(Coordinate{}) > int(this.Radius) {
		return nil, errors.New("Index out of range")
	}

	x, y := getIndices(coord, int(this.Radius))
	result := this.properties[x][y]

	if result == nil {
		this.properties[x][y] = &Property{}
	}
	return this.properties[x][y], nil
}

func (this *InMemoryMap) set(coord Coordinate, property *Property) error {
	if coord.Distance(Coordinate{}) > int(this.Radius) {
		return errors.New("Index out of range")
	}

	x, y := getIndices(coord, int(this.Radius))
	this.properties[x][y] = property

	return nil
}

func CreateInMemoryMap(radius int, fill_optinal ...bool) *InMemoryMap {
	fill := false
	if len(fill_optinal) > 0 {
		fill = fill_optinal[0]
	}
	result := InMemoryMap{}
	result.Radius = uint16(radius)

	maxX, maxY := getOffset(Coordinate{int(radius), 0, -int(radius)}, int(radius))

	result.properties = make([][]*Property, 2*radius+1)
	tmp := make([]*Property, maxX+maxY+1)
	for x := 0; x <= 2*radius; x++ {
		sizeY := -1
		if x < radius {
			sizeY = radius + x + 1
		} else {
			sizeY = 3*radius - x + 1
		}
		result.properties[x] = tmp[:sizeY]
		tmp = tmp[sizeY:]

		if fill {
			for y := 0; y < sizeY; y++ {
				result.properties[x][y] = &Property{}
			}
		}
	}

	return &result
}
