package errors

type EntityAlreadyExists struct {
	StatusCode int `default:"400"`
}

func (e EntityAlreadyExists) Error() string {
	return "entity already exists"
}
