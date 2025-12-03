package server

type Connection interface {
	Read() (string, error)
	Write(msg string) error
	Close() error
}
