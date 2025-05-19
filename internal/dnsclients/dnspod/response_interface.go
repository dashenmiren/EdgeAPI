

package dnspod

type ResponseInterface interface {
	IsOk() bool
	LastError() (code string, message string)
}
