package multiwriter

import (
	"io"
	"sync"
)

type multiWriter struct {
	mw   map[string]io.Writer
	keys []string
	mu   sync.Mutex
}

// Creates a new writer that duplicates its writes to all the attached writers.
//
// Each write is written to each attached writer, one at a time. If any writer
// returns an error, the entire write operation stops and the error is returned.
// The write is not written to the writers attached after the failing writer.
func MultiWriter() *multiWriter {
	return &multiWriter{mw: make(map[string]io.Writer)}
}

// Adds a new writer to the list of writers
func (m *multiWriter) Add(name string, writer io.Writer) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.mw[name] = writer
	m.keys = append(m.keys, name)
}

// Removes a writer from the list of writers
func (m *multiWriter) Remove(name string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for index, key := range m.keys {
		if key == name {
			m.keys = append(m.keys[:index], m.keys[index+1:]...)
		}
	}

	delete(m.mw, name)
}

// Write a slice of byte to each of the attached writers
func (m *multiWriter) Write(p []byte) (n int, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, key := range m.keys {
		if w, ok := m.mw[key]; ok {
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

// Write a string to each of the attached writers
func (m *multiWriter) WriteString(s string) (n int, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, key := range m.keys {
		if w, ok := m.mw[key]; ok {
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

func (m *multiWriter) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for name := range m.mw {
		if w, ok := m.mw[name].(io.Closer); ok {
			_ = w.Close()
		}
		delete(m.mw, name)
	}

	m.keys = nil

	return nil
}
