package hexamaps

import (
	"fmt"
	"testing"
	"time"
)

func TestInMemoryMaps(t *testing.T) {

	hexmap := CreateInMemoryMap(10, true)

	prop, _ := hexmap.get(Coordinate{1024, -251, 0})

	fmt.Println(prop)
	l := make([]*InMemoryMap, 10000)
	for i := 0; i < 10000; i++ {
		l[i] = CreateInMemoryMap(50, true)
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
