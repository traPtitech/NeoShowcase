package handler

type Handler interface {
	HandleRequest(c Context) error
}
