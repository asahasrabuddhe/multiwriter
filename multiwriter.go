package multiwriter

import (
	"io"
	"sync"
)

type multiWriter struct {
	mw map[string]io.Writer
	mu sync.Mutex
}

func NewMultiWriter() *multiWriter {
	return &multiWriter{mw: make(map[string]io.Writer)}
}

func (m *multiWriter) Add(name string, writer io.Writer) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.mw[name] = writer
}

func (m *multiWriter) Remove(name string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.mw[name]; ok {
		delete(m.mw, name)
	}
}

func (m *multiWriter) Write(p []byte) (n int, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for name := range m.mw {
		if w, ok := m.mw[name]; ok {
			n, err = w.Write(p)
			if err != nil {
				return
			}
			if n != len(p) {
				err = io.ErrShortWrite
				return
			}
		}
	}
	return len(p), nil
}

func (m *multiWriter) WriteString(s string) (n int, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for name := range m.mw {
		if w, ok := m.mw[name]; ok {
			if sw, ok := w.(io.StringWriter); ok {
				n, err = sw.WriteString(s)
			} else {
				n, err = w.Write([]byte(s))
			}
			if err != nil {
				return
			}
			if n != len(s) {
				err = io.ErrShortWrite
				return
			}
		}
	}
	return len(s), nil
}
