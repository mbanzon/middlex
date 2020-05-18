package rest

import "net/http"

type RESTAction int

const (
	GetAllAction = RESTAction(iota)
	GetOneAction
	CreateAction
	UpdateAction
	DeleteAction
)

type RESTHandler struct {
	resourceName string
}

type RESTResourceID struct{}

func (id RESTResourceID) Present() bool {
	// TODO: implement
	return false
}

func (id RESTResourceID) Get() (int64, error) {
	// TODO: implement
	return -1, nil
}

func (id RESTResourceID) GetAsString() (string, error) {
	// TODO: implement
	return "", nil
}

type GetAllFn func() ([]interface{}, error)

type GetOneFn func(int64) (interface{}, error)

type CreateFn func(interface{}) (int64, interface{}, error)

type UpdateFn func(int64, interface{}) (interface{}, error)

type DeleteFn func(int64) error

type ValidateFn func(string, RESTAction, int64) (bool, error)

type ConfigFunc func(*RESTHandler)

func New(name string) *RESTHandler {
	r := &RESTHandler{
		resourceName: name,
	}
	return r
}

func (h *RESTHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := RESTResourceID{}

	if id.Present() {
		switch r.Method {
		case http.MethodGet:
			break
		case http.MethodPut:
			break
		case http.MethodDelete:
			break
		default:
			http.Error(w, "", http.StatusBadRequest)
			break
		}
	}

	switch r.Method {
	case http.MethodGet:
		break
	case http.MethodPost:
		break
	default:
		http.Error(w, "", http.StatusBadRequest)
		break
	}
}
