package server

type NoMethodError struct {
}

func (*NoMethodError) Error() string {
	return "no such method error"
}

func (a *NoMethodError) String() string {
	return a.Error()
}
