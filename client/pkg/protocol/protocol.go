//Package protocol -
//Every message must be written and read as a tuple (*n, **content).
//*n - 4 bytes, represents a length of content
//**content - slice of bytes
package protocol

import (
	"fmt"
	"io"
)

const (
	_defaultMaxReadSize = 1024
)

type ReadWriter interface {
	// Read - reads by protocol
	Read(io.Reader) ([]byte, error)
	// Write - writes by protocol
	Write(io.Writer, []byte) error
}

type protocolT struct {
	maxReadSize uint32
}

var protocol protocolT

func init() {
	protocol = protocolT{
		maxReadSize: _defaultMaxReadSize,
	}
}

func SetMaxReadSize(maxReadSize uint32) {
	protocol.maxReadSize = maxReadSize
}

func Get() ReadWriter {
	return protocol
}

func (p protocolT) Read(r io.Reader) (content []byte, err error) {

	// init buffer for read size
	btReadSize := make([]byte, 4)
	_, err = r.Read(btReadSize)
	if err != nil {
		return
	}

	// get read size from buffer
	var readSize uint32
	var b uint32 = 1
	for i := 0; i < 4; i++ {
		readSize += b * uint32(btReadSize[i])
		b *= 256
	}

	// check read size
	if p.maxReadSize < readSize {
		err = fmt.Errorf("content size more than max read size: %d < %d", p.maxReadSize, readSize)
		return
	}

	// create buffer for content
	content = make([]byte, readSize)
	_, err = r.Read(content)
	if err != nil {
		return
	}

	return
}

func (p protocolT) Write(w io.Writer, content []byte) (err error) {

	// set write size of length of content
	writeSize := uint32(len(content))

	// convert write size to bytes
	btWriteSize := make([]byte, 4)
	for i := 0; i < 4; i++ {
		btWriteSize[i] = byte(writeSize % 256)
		writeSize /= 256
	}

	// write a size
	_, err = w.Write(btWriteSize)
	if err != nil {
		return
	}

	// write a content
	_, err = w.Write(content)
	if err != nil {
		return
	}
	return
}
