package binding

import (
	"io"

	"github.com/cloudevents/sdk-go/v2/event"
)

type MultiMessageImpl struct {
	messages []Message
	index    int
}

func NewGenericMultiMessage(messages ...Message) MultiMessage {
	return &MultiMessageImpl{
		messages: messages,
		index:    0,
	}
}

func NewEventMultiMessage(events ...event.Event) MultiMessage {
	messages := make([]Message, len(events))
	for i, e := range events {
		messages[i] = (*EventMessage)(&e)
	}
	return &MultiMessageImpl{
		messages: messages,
		index:    0,
	}
}

func (m *MultiMessageImpl) Read() (Message, error) {
	if m.index >= len(m.messages) {
		return nil, io.EOF
	}
	message := m.messages[m.index]
	m.index++
	return message, nil
}

func (m *MultiMessageImpl) Finish(error) error {
	// Allow to reuse it
	m.index = 0
	return nil
}

var _ MultiMessage = (*MultiMessageImpl)(nil)
