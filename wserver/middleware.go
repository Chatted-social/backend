package wserver

type MiddlewareFunc func(HandlerFunc) HandlerFunc

// nothing to see here for now