package events

type Fetcher interface {
	Fetch() ([]Event, error)
}

type Processor interface {
	Process(Event) error
}

type Type int

const (
	Unknown Type = iota
	Message
)

type Event struct {
	Text string
	Type Type
	Meta interface{}
}
