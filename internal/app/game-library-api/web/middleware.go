package web

// Middleware is a function to be run before or after another Handler
type Middleware func(Handler) Handler

// chainMiddleware wraps middleware around handler
func chainMiddleware(h Handler, mw []Middleware) Handler {
	for i := len(mw) - 1; i >= 0; i-- {
		handler := mw[i]
		if h != nil {
			h = handler(h)
		}
	}

	return h
}
