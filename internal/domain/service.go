package domain

type RTOBLL interface {
	HelloWorld() string
}

type Service struct {
	target float32
	test   string
}

func NewService(
	target float32,
	test string,
) RTOBLL {

	service := Service{
		target: target,
		test:   test,
	}

	return &service
}
