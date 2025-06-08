package api

// StrictServerInterface being defined like this causes unimplemented handler methods to be a compilation error.
var _ StrictServerInterface = (*server)(nil)

type server struct{}

func NewServer() server {
	return server{}
}
