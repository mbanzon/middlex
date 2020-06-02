package rest

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
