package network

import (
	"bytes"
	"encoding/binary"
	"errors"
	"net"
)

// Hoffer struct
type Hoffer struct {
	buffer     bytes.Buffer
	readOffset int
}

// NewHoffer creates a new Hoffer instance
func NewHoffer() *Hoffer {
	return &Hoffer{}
}

// Write methods
func (h *Hoffer) PutNumber(value int32) {
	binary.Write(&h.buffer, binary.BigEndian, value)
}

func (h *Hoffer) PutString(value string) {
	stringBuffer := []byte(value)
	length := int32(len(stringBuffer))
	binary.Write(&h.buffer, binary.BigEndian, length)
	h.buffer.Write(stringBuffer)
}

func (h *Hoffer) PutDouble(value float64) {
	binary.Write(&h.buffer, binary.BigEndian, value)
}

func (h *Hoffer) PutByte(value byte) {
	h.buffer.WriteByte(value)
}

func (h *Hoffer) PutByteArray(value []byte) {
	length := int32(len(value))
	binary.Write(&h.buffer, binary.BigEndian, length)
	h.buffer.Write(value)
}

func (h *Hoffer) PutValue(valueType string, value interface{}) error {
	switch valueType {
	case "number":
		if v, ok := value.(int32); ok {
			h.PutNumber(v)
		} else {
			return errors.New("invalid value type for number")
		}
	case "string":
		if v, ok := value.(string); ok {
			h.PutString(v)
		} else {
			return errors.New("invalid value type for string")
		}
	case "double":
		if v, ok := value.(float64); ok {
			h.PutDouble(v)
		} else {
			return errors.New("invalid value type for double")
		}
	case "byte":
		if v, ok := value.(byte); ok {
			h.PutByte(v)
		} else {
			return errors.New("invalid value type for byte")
		}
	case "byteArray":
		if v, ok := value.([]byte); ok {
			h.PutByteArray(v)
		} else {
			return errors.New("invalid value type for byteArray")
		}
	default:
		return errors.New("unsupported type")
	}
	return nil
}

// Read methods
func (h *Hoffer) GetNumber() (int32, error) {
	if h.readOffset+4 > h.buffer.Len() {
		return 0, errors.New("buffer underflow")
	}
	var value int32
	err := binary.Read(bytes.NewReader(h.buffer.Bytes()[h.readOffset:]), binary.BigEndian, &value)
	if err != nil {
		return 0, err
	}
	h.readOffset += 4
	return value, nil
}

func (h *Hoffer) GetString() (string, error) {
	length, err := h.GetNumber()
	if err != nil {
		return "", err
	}
	if h.readOffset+int(length) > h.buffer.Len() {
		return "", errors.New("buffer underflow")
	}
	value := string(h.buffer.Bytes()[h.readOffset : h.readOffset+int(length)])
	h.readOffset += int(length)
	return value, nil
}

func (h *Hoffer) GetDouble() (float64, error) {
	if h.readOffset+8 > h.buffer.Len() {
		return 0, errors.New("buffer underflow")
	}
	var value float64
	err := binary.Read(bytes.NewReader(h.buffer.Bytes()[h.readOffset:]), binary.BigEndian, &value)
	if err != nil {
		return 0, err
	}
	h.readOffset += 8
	return value, nil
}

func (h *Hoffer) GetByte() (byte, error) {
	if h.readOffset+1 > h.buffer.Len() {
		return 0, errors.New("buffer underflow")
	}
	value := h.buffer.Bytes()[h.readOffset]
	h.readOffset++
	return value, nil
}

func (h *Hoffer) GetByteArray() ([]byte, error) {
	length, err := h.GetNumber()
	if err != nil {
		return nil, err
	}
	if h.readOffset+int(length) > h.buffer.Len() {
		return nil, errors.New("buffer underflow")
	}
	value := h.buffer.Bytes()[h.readOffset : h.readOffset+int(length)]
	h.readOffset += int(length)
	return value, nil
}

func (h *Hoffer) GetValue(valueType string) (interface{}, error) {
	switch valueType {
	case "number":
		return h.GetNumber()
	case "string":
		return h.GetString()
	case "double":
		return h.GetDouble()
	case "byte":
		return h.GetByte()
	case "byteArray":
		return h.GetByteArray()
	default:
		return nil, errors.New("unsupported type")
	}
}

// Reset buffer
func (h *Hoffer) Reset() {
	h.buffer.Reset()
	h.readOffset = 0
}

// GetData returns the buffer's bytes
func (h *Hoffer) GetData() []byte {
	return h.buffer.Bytes()
}

// SetData sets the buffer's bytes
func (h *Hoffer) SetData(data []byte) {
	h.buffer = *bytes.NewBuffer(data)
	h.readOffset = 0
}

// SendData sends the buffer's data over a net.Conn
func (h *Hoffer) SendData(conn net.Conn) error {
	_, err := conn.Write(h.buffer.Bytes())
	return err
}
