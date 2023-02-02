package api

// Error represents a handler error. It provides methods for a HTTP status
// code and embeds the built-in error interface.
type Error interface {
	error
	Status() int
}

type StatusError struct {
	Code int
	Err  error
}
