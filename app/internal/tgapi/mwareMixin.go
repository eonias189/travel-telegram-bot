package tgapi

type Middleware func(next HandlerFunc) HandlerFunc

type MiddlewareMixin struct {
	mwares []Middleware
}

func (mm *MiddlewareMixin) Use(mware Middleware) {
	if mware != nil {
		mm.mwares = append(mm.mwares, mware)
	}
}

func (mm *MiddlewareMixin) WithMiddlewares(handler HandlerFunc) HandlerFunc {
	next := handler
	for _, mw := range mm.mwares {
		next = mw(next)
	}
	return next
}
