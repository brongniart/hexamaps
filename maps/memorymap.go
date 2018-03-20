package maps

import (
	"errors"
	"sync"
)

import . "github.com/brongniart/hexamaps"

type InMemory struct {
	radius     int
	properties [][]*Property

	mutex sync.RWMutex
}

func (this *InMemory) get(coord Coordinate) (property Property, err error) {
	if coord.Distance(Coordinate{}) > int(this.radius) {
		err = errors.New("Index out of range")
		return
	}

	x, y := getIndices(coord, int(this.radius))

	this.mutex.RLock()
	result := this.properties[x][y]
	property = *result
	this.mutex.RUnlock()

	return
}

func (this *InMemory) set(coord Coordinate, property Property) (err error) {
	if coord.Distance(Coordinate{}) > int(this.radius) {
		return errors.New("Index out of range")
	}

	x, y := getIndices(coord, int(this.radius))

	this.mutex.Lock()
	this.properties[x][y] = &property
	this.mutex.Unlock()

	return
}

func CreateInMemoryMap(radius int) *InMemory {

	result := InMemory{}
	result.radius = radius

	max_x, max_y := getOffset(CreateCoordinate(radius, 0), radius)

	defaultProp := new(Property)

	result.properties = make([][]*Property, 2*radius+1)
	slicing := make([]*Property, max_x+max_y+1)

	for x := 0; x <= 2*radius; x++ {
		sizeY := -1
		if x < radius {
			sizeY = radius + x + 1
		} else {
			sizeY = 3*radius - x + 1
		}
		result.properties[x] = slicing[:sizeY]
		slicing = slicing[sizeY:]

		for y := 0; y < sizeY; y++ {
			result.properties[x][y] = defaultProp
		}
	}

	return &result
}
