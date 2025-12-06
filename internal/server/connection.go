package server

type Connection interface {
	Read() (string, error)
	ReadAndGetData() (*IncomingMessage, error)
	Write(v interface{}) error
	Close() error
}
