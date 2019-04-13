package crud

import (
	"net/http"
)

type JSONCRUDType interface {
}

type JSONCRUDObject interface {
	ID() int64
}

type CRUD struct {
}

func (c *CRUD) Wrap(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: check if there is a parameter
		hasParam, param := false, "" // FIXME: these should have real values!

		switch r.Method {
		case http.MethodPost && !hasParam:
			// TODO: verify posted data
			// TODO: create new object with data
			// TODO: return created object
			break
		case http.MethodGet && hasParam:
			// TODO: get object with ID = param
			// TODO: return object
			break
		case http.MethodGet && !hasParam:
			break
		case http.MethodPut:
			break
		case http.MethodDelete:
			break
		default:
			// TODO: fail
			break
		}
	})
}
