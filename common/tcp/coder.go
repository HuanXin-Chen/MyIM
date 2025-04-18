package tcp

import (
	"bytes"
	"encoding/binary"
)

type DataPgk struct {
	Len  uint32
	Data []byte
}

func (d *DataPgk) Marshal() []byte {
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, d.Len) //使用大端字节序
	return append(bytesBuffer.Bytes(), d.Data...)
}
