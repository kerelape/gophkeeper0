package composedreadcloser

import "io"

// ComposedReadCloser is a composed io.ReadCloser.
type ComposedReadCloser struct {
	Reader io.Reader
	Closer io.Closer
}

// Read implements io.Reader.
func (rc *ComposedReadCloser) Read(p []byte) (int, error) {
	return rc.Reader.Read(p)
}

// Close implements io.ReadCloser.
func (rc *ComposedReadCloser) Close() error {
	return rc.Closer.Close()
}
