package sequencereader

import (
	"errors"
	"io"
)

// SequenceReader reads from sources sequencely.
type SequenceReader []io.Reader

var _ io.Reader = (*SequenceReader)(nil)

// Read implements io.Reader.
func (r *SequenceReader) Read(p []byte) (int, error) {
	if len(*r) > 0 {
		n, err := (*r)[0].Read(p)
		if err != nil {
			if !errors.Is(err, io.EOF) {
				return n, err
			}
			*r = (*r)[1:]
			return n, nil
		}
	}
	return 0, io.EOF
}
