package web

type Handler interface {
	HandleRequest(c Context) error
}
