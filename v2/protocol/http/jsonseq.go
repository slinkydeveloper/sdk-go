package http

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/jmank88/jsonseq"

	"github.com/cloudevents/sdk-go/v2/binding"
	"github.com/cloudevents/sdk-go/v2/binding/format"
)

type JsonSeqReader struct {
	scanner *bufio.Scanner
}

func NewJsonSeqReader(reader io.Reader) *JsonSeqReader {
	s := bufio.NewScanner(reader)
	s.Split(jsonseq.ScanRecord)
	return &JsonSeqReader{scanner: bufio.NewScanner(reader)}
}

func (j *JsonSeqReader) ReadNext() (binding.Message, error) {
	if !j.scanner.Scan() {
		if err := j.scanner.Err(); err != nil {
			return nil, err
		}
		return nil, io.EOF
	}
	b := j.scanner.Bytes()

	b, ok := jsonseq.RecordValue(b)
	if !ok {
		return nil, fmt.Errorf("invalid record: %q", string(b))
	}
	return &JsonSeqMessage{Bytes: b, Format: format.JSON}, nil
}

// Copy pasted from binding/test/mock_structured_message.go
type JsonSeqMessage struct {
	Format format.Format
	Bytes  []byte
}

func (s *JsonSeqMessage) ReadStructured(ctx context.Context, b binding.StructuredWriter) error {
	return b.SetStructuredEvent(ctx, s.Format, bytes.NewReader(s.Bytes))
}

func (s *JsonSeqMessage) ReadBinary(context.Context, binding.BinaryWriter) error {
	return binding.ErrNotBinary
}

func (s *JsonSeqMessage) ReadEncoding() binding.Encoding {
	return binding.EncodingStructured
}

func (s *JsonSeqMessage) Finish(error) error { return nil }

func (s *JsonSeqMessage) SetStructuredEvent(ctx context.Context, format format.Format, event io.Reader) (err error) {
	s.Format = format
	s.Bytes, err = ioutil.ReadAll(event)
	if err != nil {
		return
	}

	return nil
}

var _ binding.Message = (*JsonSeqMessage)(nil)
var _ binding.StructuredWriter = (*JsonSeqMessage)(nil)

const (
	rs = 0x1E
	lf = 0x0A
)

type JsonSeqWriter struct {
	writer io.Writer
}

func (j *JsonSeqWriter) SetStructuredEvent(ctx context.Context, format format.Format, event io.Reader) error {
	_, err := j.writer.Write([]byte{rs})
	if err != nil {
		return err
	}
	_, err = io.Copy(j.writer, event)
	if err != nil {
		return err
	}
	_, err = j.writer.Write([]byte{lf})
	if err != nil {
		return err
	}
	return nil
}
