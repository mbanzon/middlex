package rest

import (
	"encoding/json"
	"net/http"
)

type RESTAction int

const (
	GetAllAction = RESTAction(iota)
	GetOneAction
	CreateAction
	UpdateAction
	DeleteAction
)

type RESTHandler struct {
	resourceName   string
	GetAll         func() ([]interface{}, error)
	GetOne         func(RESTResourceID) (interface{}, error)
	Create         func(interface{}) (RESTResourceID, interface{}, error)
	Update         func(RESTResourceID, interface{}) (interface{}, error)
	Delete         func(RESTResourceID) error
	ValidateCreate func(json.RawMessage) (bool, interface{}, error)
	ValidateUpdate func(RESTResourceID, json.RawMessage) (bool, interface{}, error)
}

type ConfigFunc func(*RESTHandler)

func New(name string) *RESTHandler {
	r := &RESTHandler{
		resourceName: name,
	}
	return r
}

func (h *RESTHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := ParseRESTResourceID(r.RequestURI, h.resourceName)

	if id.Present() {
		switch r.Method {
		case http.MethodGet:
			if h.GetOne != nil {
				h.GetOne(id)
			}
			break
		case http.MethodPut:
			if h.Update != nil && h.ValidateUpdate != nil {
				// TODO: get data as raw JSON
				ok, data, err := h.ValidateUpdate(id, nil)
				if err != nil {
					// TODO: handle error
				}
				if !ok {
					// TODO: handle data not ok
				}
				h.Update(id, data)
			}
			break
		case http.MethodDelete:
			if h.Delete != nil {
				h.Delete(id)
			}
			break
		default:
			http.Error(w, "", http.StatusBadRequest)
			break
		}
	}

	switch r.Method {
	case http.MethodGet:
		if h.GetAll != nil {
			h.GetAll()
		}
		break
	case http.MethodPost:
		// TODO: get data as raw JSON
		ok, data, err := h.ValidateCreate(nil)
		if err != nil {
			// TODO: handle error
		}
		if !ok {
			// TODO: handle data not ok
		}
		h.Create(data)
		break
	default:
		http.Error(w, "", http.StatusBadRequest)
		break
	}

	// TODO: handle things if we get to this point = the request is bad
}
