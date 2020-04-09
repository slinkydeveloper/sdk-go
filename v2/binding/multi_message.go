package binding

type MultiMessage interface {
	Read() (Message, error)
	Finish(error) error
}
