package binding

import "context"

type MultiWriter interface {
	Start(ctx context.Context) error
	Write(Message) error
	End(ctx context.Context) error
}
