package rest

import (
	"errors"
	"strconv"
	"strings"
)

type RESTResourceID struct {
	raw           string
	asInt         int64
	asIntError    error
	asString      string
	asStringError error
}

var NoResourceIDError = errors.New("no resource id in uri")

func ParseRESTResourceID(uri string, resourceName string) RESTResourceID {
	rID := RESTResourceID{
		raw: uri,
	}

	uri = strings.Trim(uri, "/")
	parts := strings.Split(uri, "/")
	if len(parts) == 0 || parts[len(parts)-1] == resourceName || parts[len(parts)-1] == "" {
		rID.asInt, rID.asString = -1, ""
		rID.asIntError, rID.asStringError = NoResourceIDError, NoResourceIDError
		return rID
	}

	last := parts[len(parts)-1]
	rID.asString, rID.asStringError = last, nil
	rID.asInt, rID.asIntError = strconv.ParseInt(last, 10, 64)
	return rID
}

func (id RESTResourceID) Present() bool {
	return id.asInt != -1 || id.asString != ""
}

func (id RESTResourceID) Get() int64 {
	return id.asInt
}

func (id RESTResourceID) GetError() error {
	return id.asIntError
}

func (id RESTResourceID) GetAsString() string {
	return id.asString
}

func (id RESTResourceID) GetAsStringError() error {
	return id.asStringError
}
