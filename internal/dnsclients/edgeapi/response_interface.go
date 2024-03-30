package edgeapi

type ResponseInterface interface {
	IsValid() bool
	Error() error
}
