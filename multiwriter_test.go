package multiwriter

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"testing"
)

func TestMultiWriter_Write(t *testing.T) {
	mw := NewMultiWriter()

	bufOne := new(bytes.Buffer)
	bufTwo := new(bytes.Buffer)
	bufThree := new(bytes.Buffer)

	mw.Add("one", bufOne)
	mw.Add("two", bufTwo)

	for i := 0; i <= 10; i++ {
		if i%2 == 0 {
			mw.Remove("one")
			mw.Add("three", bufThree)
		} else {
			mw.Remove("three")
			mw.Add("one", bufOne)
		}

		n, err := mw.Write([]byte(fmt.Sprintf("%v", i)))
		if err != nil {
			t.Error("fail:", err)
		}

		if n != len(fmt.Sprintf("%v", i)) {
			t.Errorf("fail: expected written bytes = %d, got %d", len(fmt.Sprintf("%v", i)), n)
		}
	}

	if bufOne.String() != "13579" {
		t.Errorf("fail: expected %s, got %s", "13579", bufOne.String())
	}

	if bufTwo.String() != "012345678910" {
		t.Errorf("fail: expected %s, got %s", "012345678910", bufTwo.String())
	}

	if bufThree.String() != "0246810" {
		t.Errorf("fail: expected %s, got %s", "0246810", bufThree.String())
	}
}

func TestMultiWriter_WriteString(t *testing.T) {
	mw := NewMultiWriter()

	bufOne := new(bytes.Buffer)
	bufTwo := new(bytes.Buffer)
	bufThree := new(bytes.Buffer)

	mw.Add("one", bufOne)
	mw.Add("two", bufTwo)

	for i := 0; i <= 10; i++ {
		if i%2 == 0 {
			mw.Remove("one")
			mw.Add("three", bufThree)
		} else {
			mw.Remove("three")
			mw.Add("one", bufOne)
		}

		n, err := mw.WriteString(fmt.Sprintf("%v", i))
		if err != nil {
			t.Error("fail:", err)
		}

		if n != len(fmt.Sprintf("%v", i)) {
			t.Errorf("fail: expected written bytes = %d, got %d", len(fmt.Sprintf("%v", i)), n)
		}
	}

	if bufOne.String() != "13579" {
		t.Errorf("fail: expected %s, got %s", "13579", bufOne.String())
	}

	if bufTwo.String() != "012345678910" {
		t.Errorf("fail: expected %s, got %s", "012345678910", bufTwo.String())
	}

	if bufThree.String() != "0246810" {
		t.Errorf("fail: expected %s, got %s", "0246810", bufThree.String())
	}
}

type failWriter struct {
	written        int
	failAfterBytes int
	w              io.Writer
}

func (f *failWriter) Write(p []byte) (n int, err error) {
	if f.written >= f.failAfterBytes {
		return f.written, errors.New("write failed")
	}

	if len(p) > f.failAfterBytes-f.written {
		p = p[:f.failAfterBytes-f.written]
	}

	n, err = f.w.Write(p)
	if err != nil {
		return
	}

	f.written += n

	return
}

func (f *failWriter) WriteString(s string) (n int, err error) {
	return f.Write([]byte(s))
}

func TestMultiWriter_WriteError(t *testing.T) {
	mw := NewMultiWriter()

	bufOne := new(bytes.Buffer)
	bufTwo := new(bytes.Buffer)

	wOne := &failWriter{failAfterBytes: 3, w: bufOne}

	mw.Add("one", wOne)
	mw.Add("two", bufTwo)

	n, err := mw.Write([]byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ"))
	if n != 3 {
		t.Errorf("fail: expected %d bytes to be written, got %d", 3, n)
	}

	if err == nil {
		t.Error("fail: expected error")
	}

	if bufOne.String() != "ABC" && bufOne.String() != bufTwo.String() {
		t.Errorf("fail: expected %s, got %s", "ABC", bufOne.String())
	}
}

func TestMultiWriter_WriteStringError(t *testing.T) {
	mw := NewMultiWriter()

	bufOne := new(bytes.Buffer)
	bufTwo := new(bytes.Buffer)

	wOne := &failWriter{failAfterBytes: 6, w: bufOne}

	mw.Add("one", wOne)
	mw.Add("two", bufTwo)

	n, err := mw.WriteString("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	if n != 6 {
		t.Errorf("fail: expected %d bytes to be written, got %d", 6, n)
	}

	if err == nil {
		t.Error("fail: expected error")
	}

	if bufOne.String() != "ABCDEF" && bufOne.String() != bufTwo.String() {
		t.Errorf("fail: expected %s, got %s", "ABCDEF", bufOne.String())
	}
}
