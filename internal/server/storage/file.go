package storage

import "bytes"

type File struct {
	name   string
	buffer *bytes.Buffer
}

func (f *File) Write(chunk []byte) error {
	_, err := f.buffer.Write(chunk)

	return err
}
