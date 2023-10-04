package server

type Config struct {
	Addr      string
	EnableLog bool
	Routes    []Route
}
