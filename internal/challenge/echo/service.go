package echo

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Echo(value string) string {
	return value
}
