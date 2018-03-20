package hexamaps

import (
	"bytes"
	"encoding/binary"
	"errors"
	"os"
)

type OnFileMap struct {
	Radius uint16
	file   *os.File
}

type header struct {
	Radius uint16
}

func (this *OnFileMap) get(coord Coordinate) (*Property, error) {
	if coord.Distance(Coordinate{}) > int(this.Radius) {
		return nil, errors.New("Index out of range")
	}

	offsetX, offsetY := getOffset(coord, int(this.Radius))
	return this.readProperty(offsetX + offsetY)
}

func (this *OnFileMap) set(coord Coordinate, property *Property) error {
	if coord.Distance(Coordinate{}) > int(this.Radius) {
		return errors.New("Index out of range")
	}

	offsetX, offsetY := getOffset(coord, int(this.Radius))
	return this.writeProperty(property, offsetX+offsetY)
}

func (this *OnFileMap) readHeader() (head header, err error) {
	sizeHeader := binary.Size(header{})

	buf := make([]byte, sizeHeader)

	n, err := this.file.ReadAt(buf, 0)
	if err != nil {
		return
	}

	if n != binary.Size(header{}) {
		err = errors.New("Incorrect header size")
		return
	}
	binary.Read(bytes.NewBuffer(buf), binary.LittleEndian, &head)

	return
}

func (this *OnFileMap) writeHeader(head header) (err error) {
	buf := new(bytes.Buffer)

	err = binary.Write(buf, binary.LittleEndian, head)
	if err != nil {
		return
	}

	n, err := this.file.WriteAt(buf.Bytes(), 0)
	if err != nil {
		return
	}
	if n != binary.Size(header{}) {
		err = errors.New("Incorrect header size")
		return
	}
	this.file.Sync()

	return nil
}

func (this *OnFileMap) readProperty(nb int) (*Property, error) {
	prop := Property{}

	sizeProperty := binary.Size(Property{})
	sizeHeader := binary.Size(header{})

	buf := make([]byte, sizeProperty)

	n, err := this.file.ReadAt(buf, int64(sizeHeader+nb*sizeProperty))
	if err != nil {
		return nil, err
	}
	if n != binary.Size(Property{}) {
		return nil, errors.New("Incorrect property size")
	}

	binary.Read(bytes.NewBuffer(buf), binary.LittleEndian, &prop)

	return &prop, nil
}

func (this *OnFileMap) writeProperty(property *Property, nb int) error {

	sizeProperty := binary.Size(Property{})
	sizeHeader := binary.Size(header{})

	buf := new(bytes.Buffer)

	err := binary.Write(buf, binary.LittleEndian, property)
	if err != nil {
		return err
	}

	n, err := this.file.WriteAt(buf.Bytes(), int64(sizeHeader+nb*sizeProperty))
	if err != nil {
		return err
	}
	if n != binary.Size(Property{}) {
		return errors.New("Incorrect property size")
	}

	return nil
}

func (this *OnFileMap) writeBunchEmptyProperties(index int) error {
	offsetX, nbProperties := -1, -1

	if index <= int(this.Radius) {
		coord := Coordinate{index - int(this.Radius), int(this.Radius), -index}
		offsetX, nbProperties = getOffset(coord, int(this.Radius))
	} else {
		coord := Coordinate{index - int(this.Radius), -index, int(this.Radius)}
		offsetX, nbProperties = getOffset(coord, int(this.Radius))
	}

	sizeProperty := binary.Size(Property{})
	sizeHeader := binary.Size(header{})

	buf := new(bytes.Buffer)
	for i := 0; i < nbProperties; i++ {
		err := binary.Write(buf, binary.LittleEndian, Property{})
		if err != nil {
			return err
		}
	}

	this.file.WriteAt(buf.Bytes(), int64(sizeHeader+offsetX*sizeProperty))
	return nil
}

func CreateEmptyOnFileMap(path string, radius uint16) *OnFileMap {
	result := OnFileMap{}
	result.Radius = radius

	file, error := os.Create(path)
	if error != nil {
		panic(error)
	}
	result.file = file
	head := header{result.Radius}
	error = result.writeHeader(head)
	if error != nil {
		panic(error)
	}

	for i := 0; i < 2*int(result.Radius)+1; i++ {
		error = result.writeBunchEmptyProperties(i)
		if error != nil {
			panic(error)
		}
	}

	return &result
}

func CreateOnFileMap(path string) *OnFileMap {
	result := OnFileMap{}
	result.Radius = 0

	file, error := os.OpenFile(path, os.O_RDWR, 0755)
	if error != nil {
		panic(error)
	}
	result.file = file

	head, err := result.readHeader()
	if err != nil {
		panic(err)
	}
	result.Radius = head.Radius

	return &result
}
