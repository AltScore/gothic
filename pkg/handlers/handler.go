package handlers

type HandlerSpec struct {
	Name        string
	Route       string
	Method      string
	Description string
}

type Handler[RequestBody any, ResponseBody any] interface {
}

type handler[RequestBody any, ResponseBody any] struct {
	handlerSpec HandlerSpec
}

func New[RequestBody any, ResponseBody any](
	HandlerSpec HandlerSpec,
) Handler[RequestBody, ResponseBody] {
	return &handler[RequestBody, ResponseBody]{}
}
