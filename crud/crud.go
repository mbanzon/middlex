package crud

import (
	"net/http"
	"strconv"
	"strings"
)

type ListerFunc func() ([]interface{}, error)
type GetterFunc func(int64) (interface{}, error)
type PosterFunc func(interface{}) (int64, interface{}, error)
type PutterFunc func(int64, interface{}) (interface{}, error)
type DeleterFunc func(int64) error

type CRUD struct {
	lister  http.Handler
	getter  http.Handler
	poster  http.Handler
	putter  http.Handler
	deleter http.Handler
}

func New() *CRUD {
	return &CRUD{}
}

func (c *CRUD) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		found, id, err := lastURIArgAsInt64(r)

		if r.Method == http.MethodGet {
			if !found && c.lister != nil {
				c.lister.ServeHTTP(w, r)
				return
			} else if err == nil && c.getter != nil {
				c.getter.ServeHTTP(w, r)
				return
			}
		} else if r.Method == http.MethodPost {

		} else if r.Method == http.MethodPut {

		} else if r.Method == http.MethodPut {

		} else if r.Method == http.MethodDelete {

		}
	})
}

func lastURIArg(r *http.Request) string {
	uri := strings.Trim(r.RequestURI, "/")

	if uri == "" {
		return ""
	}

	var arg string
	index := strings.LastIndex(uri, "/")
	if index != -1 {
		arg = uri[:index]
	} else {
		arg = uri
	}
	return arg
}

func lastURIArgAsInt64(r *http.Request) (bool, int64, error) {
	arg := lastURIArg(r)
	if arg == "" {
		return false, -1, nil
	}
	id, err := strconv.ParseInt(arg, 10, 64)
	return true, id, err
}
