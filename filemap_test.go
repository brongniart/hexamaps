package hexamaps

import (
	"fmt"
	"testing"
	"time"
)

func TestOnFileMaps(t *testing.T) {

	start := time.Now()

	//onFileMap := CreateOnFileMap("C:\\tmp\\test.data")
	onFileMap := CreateEmptyOnFileMap("C:\\tmp\\test.data", 1024)

	a := Coordinate{0, 0, 0}
	next := a.Neighboors(1)
	test := make(map[Coordinate]Property)
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
			onFileMap.set(hex, &prop)

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
			}*/
		test[hex] = Property{} //*propR
		//fmt.Println(hex)
	}
	elapsed := time.Since(start)
	fmt.Println(" took ", elapsed, i)

	onFileMap.file.Close()

	//time.Sleep(time.Second * 20)
}
