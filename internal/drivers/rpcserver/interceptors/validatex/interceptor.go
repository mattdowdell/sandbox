// ...
package validatex

import (
	"connectrpc.com/validate"
)

// ...
type Interceptor = validate.Interceptor

// ...
func New() (*Interceptor, error) {
	return validate.NewInterceptor()
}
