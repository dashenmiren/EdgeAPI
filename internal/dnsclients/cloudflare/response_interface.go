package cloudflare

type ResponseInterface interface {
	IsOk() bool
	LastError() (code int, message string)
}
