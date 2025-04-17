// ...
package validatex

import (
	"connectrpc.com/validate"
)

// ...
func New() (*validate.Interceptor, error) {
	return validate.NewInterceptor()
}
