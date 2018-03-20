package maps

import (
	"fmt"
	"testing"
	"time"
)

import . "github.com/brongniart/hexamaps"

func TestInMemoryMaps(t *testing.T) {

	hexmap := CreateInMemoryMap(1024)

	prop, _ := hexmap.get(CreateCoordinate(1024, -251))

	fmt.Println(prop)
	l := make([]*InMemory, 10000)
	for i := 0; i < 10000; i++ {
		l[i] = CreateInMemoryMap(50)
	}

	a := Coordinate{}
	next := a.Neighboors(256)
	start := time.Now()
	i := 0
	for {
		stop, hex := next()
		if stop {
			break
		}
		i++
		hex.GetX()
		//fmt.Println(hex)
	}
	elapsed := time.Since(start)
	fmt.Println("iter: ", i, " took ", elapsed)
	time.Sleep(time.Second * 10)
}
