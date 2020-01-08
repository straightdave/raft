package event

// Event ...
type Event struct {
	typ  Type
	data interface{}
}

// Callback ...
type Callback func(e Event) error

// Listener ...
type Listener interface {
	Type() Type
	Notify(e *Event)
}

// NewListener ...
func NewListener(typ Type, cb Callback) *Listener {
	return &Listener{
		typ: typ,
		cb:  cb,
	}
}

// Type ...
func (l *Listener) Type() Type {
	return l.typ
}

// Callback ...
func (l *Listener) Callback() Callback {
	return l.cb
}
