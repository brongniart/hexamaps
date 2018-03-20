package maps

import (
	"bytes"
	"encoding/binary"
	"errors"
	"os"
	"sync"
)
import . "github.com/brongniart/hexamaps"

type OnFile struct {
	radius int
	file   *os.File

	mutex sync.Mutex
}

type header struct {
	Radius uint16
}

func (this *OnFile) get(coord Coordinate) (property Property, err error) {
	if coord.Distance(Coordinate{}) > this.radius {
		return Property{}, errors.New("Index out of range")
	}

	offsetX, offsetY := getOffset(coord, this.radius)

	this.mutex.Lock()
	property, err = this.readProperty(offsetX + offsetY)
	this.mutex.Unlock()

	return
}

func (this *OnFile) set(coord Coordinate, property Property) (err error) {
	if coord.Distance(Coordinate{}) > this.radius {
		return errors.New("Index out of range")
	}

	offsetX, offsetY := getOffset(coord, int(this.radius))

	this.mutex.Lock()
	err = this.writeProperty(property, offsetX+offsetY)
	this.mutex.Unlock()

	return
}

func (this *OnFile) readHeader() (head header, err error) {
	sizeHeader := binary.Size(header{})

	buf := make([]byte, sizeHeader)

	n, err := this.file.ReadAt(buf, 0)
	if err != nil {
		return
	}
	if n != sizeHeader {
		err = errors.New("Incorrect header size")
		return
	}
	binary.Read(bytes.NewBuffer(buf), binary.LittleEndian, &head)
	return
}

func (this *OnFile) writeHeader(head header) (err error) {
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
		return errors.New("Incorrect header size")
	}
	this.file.Sync()

	return nil
}

func (this *OnFile) readProperty(index int) (prop Property, err error) {
	sizeProperty := binary.Size(Property{})
	sizeHeader := binary.Size(header{})

	buf := make([]byte, sizeProperty)

	n, err := this.file.ReadAt(buf, int64(sizeHeader+index*sizeProperty))
	if err != nil {
		return
	}
	if n != sizeProperty {
		err = errors.New("Incorrect property size")
		return
	}

	binary.Read(bytes.NewBuffer(buf), binary.LittleEndian, &prop)

	return prop, nil
}

func (this *OnFile) writeProperty(property Property, index int) error {

	sizeProperty := binary.Size(Property{})
	sizeHeader := binary.Size(header{})

	buf := new(bytes.Buffer)

	err := binary.Write(buf, binary.LittleEndian, property)
	if err != nil {
		return err
	}

	n, err := this.file.WriteAt(buf.Bytes(), int64(sizeHeader+index*sizeProperty))
	if err != nil {
		return err
	}
	if n != sizeProperty {
		return errors.New("Incorrect property size")
	}

	return nil
}

func (this *OnFile) readProperties(index int, nb int) (properties []Property, err error) {
	sizeProperty := binary.Size(Property{})
	sizeHeader := binary.Size(header{})

	buf := make([]byte, sizeProperty*nb)

	n, err := this.file.ReadAt(buf, int64(sizeHeader+index*sizeProperty))
	if err != nil {
		return
	}
	if n != sizeProperty*nb {
		err = errors.New("Incorrect properties size")
		return
	}

	binary.Read(bytes.NewBuffer(buf), binary.LittleEndian, &properties)

	return
}

func (this *OnFile) writeProperties(index int, nb int) (properties []Property, err error) {
	sizeProperty := binary.Size(Property{})
	sizeHeader := binary.Size(header{})

	buf := make([]byte, sizeProperty*nb)

	n, err := this.file.ReadAt(buf, int64(sizeHeader+index*sizeProperty))
	if err != nil {
		return
	}
	if n != sizeProperty*nb {
		err = errors.New("Incorrect properties size")
		return
	}

	binary.Read(bytes.NewBuffer(buf), binary.LittleEndian, &properties)

	return
}

func (this *OnFile) writeBunchEmptyProperties(index int) error {
	offsetX, nbProperties := -1, -1

	if index <= int(this.radius) {
		coord := CreateCoordinate(index-this.radius, this.radius)
		offsetX, nbProperties = getOffset(coord, this.radius)
	} else {
		coord := CreateCoordinate(index-this.radius, -index)
		offsetX, nbProperties = getOffset(coord, this.radius)
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

func CreateEmptyOnFileMap(path string, radius int) (result *OnFile, err error) {
	result = new(OnFile)

	if radius > 0 {
		result.radius = radius
	} else {
		err = errors.New("Radius must > 0")
		return
	}
	file, err := os.Create(path)
	if err != nil {
		return
	}
	result.file = file
	head := header{uint16(result.radius)}
	err = result.writeHeader(head)
	if err != nil {
		return
	}

	for i := 0; i < 2*int(result.radius)+1; i++ {
		err = result.writeBunchEmptyProperties(i)
		if err != nil {
			return
		}
	}

	return
}

func CreateOnFileMap(path string) (result *OnFile, err error) {
	result = new(OnFile)
	result.radius = -1

	file, err := os.OpenFile(path, os.O_RDWR, 0755)
	if err != nil {
		return
	}
	result.file = file

	head, err := result.readHeader()
	if err != nil {
		return
	}

	if int(head.Radius) > 0 {
		result.radius = int(head.Radius)
	} else {
		err = errors.New("Radius must >0")
		return
	}

	return
}
