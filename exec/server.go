package exec

type ServeConfig interface {
	GetPort() int
}

func Serve(config ServeConfig) error {
	return nil
}
