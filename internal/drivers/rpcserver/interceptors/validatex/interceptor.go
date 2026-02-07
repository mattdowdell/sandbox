// ...
package validatex

import (
	"connectrpc.com/validate"
)

// ...
func New() *validate.Interceptor {
	return validate.NewInterceptor()
}
