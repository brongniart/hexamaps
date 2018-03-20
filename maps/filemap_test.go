package maps

import (
	"fmt"
	"testing"
	"time"
)

import . "github.com/brongniart/hexamaps"

func TestOnFileMaps(t *testing.T) {

	start := time.Now()

	onFileMap, err := CreateOnFileMap("C:\\tmp\\test.data")
	if err != nil {
		panic(err)
	}
	//onFileMap := CreateEmptyOnFileMap("C:\\tmp\\test.data", 1024)

	a := CreateCoordinate(0, 0)
	next := a.Neighboors(10)
	start = time.Now()
	i := 0
	for {
		stop, hex := next()
		if stop {
			break
		}
		i++
		/*
			prop := Property{int16(hex.GetX()), int16(hex.GetY()), int16(hex.GetZ())}
			onFileMap.set(hex, prop)
		*/
		propR, err := onFileMap.get(hex)
		if err != nil {
			fmt.Println(propR, hex)
			panic(err)
		}
		if int(propR.Value1) != hex.GetX() ||
			int(propR.Value2) != hex.GetY() ||
			int(propR.Value3) != hex.GetZ() {
			fmt.Println(propR, hex)
			panic("error")
		}
		//fmt.Println(hex)
	}
	elapsed := time.Since(start)
	fmt.Println(" took ", elapsed, i)

	onFileMap.file.Close()

	//time.Sleep(time.Second * 20)
}
