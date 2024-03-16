package tgrouter

import (
	"context"
	"fmt"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type HandlerFunc func(ctx *Context) error
type Middleware func(next HandlerFunc) HandlerFunc

type Router interface {
	Handle(command string, handler HandlerFunc)
	Run(ctx context.Context, token string) error
	Close()
	Use(middleware Middleware) Router
}

type routerImpl struct {
	routes      map[string]HandlerFunc
	middlewares []Middleware
	cfg         *Config
	wg          sync.WaitGroup
	mu          sync.RWMutex
}

func (r *routerImpl) Handle(command string, handler HandlerFunc) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	r.routes[command] = handler
}

func (r *routerImpl) Use(mw Middleware) Router {
	r.mu.RLock()
	defer r.mu.RUnlock()
	r.middlewares = append(r.middlewares, mw)
	return r
}

func (r *routerImpl) withMWares(handler HandlerFunc) HandlerFunc {
	res := handler
	for _, mw := range r.middlewares {
		res = mw(res)
	}
	return res
}

func (r *routerImpl) handle(handler HandlerFunc, ctx *Context) {
	defer func() {
		err := recover()
		if err != nil {
			r.cfg.OnError(err.(error))
		}
	}()
	err := r.withMWares(handler)(ctx)
	if err != nil {
		r.cfg.OnError(err)
	}
}

func (r *routerImpl) Run(ctx context.Context, token string) error {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return err
	}

	updateCfg := tgbotapi.NewUpdate(0)
	ch := bot.GetUpdatesChan(updateCfg)

	for i := 0; i < r.cfg.Workers; i++ {
		r.wg.Add(1)
		go func() {
			defer r.wg.Done()
			for {
				select {
				case update := <-ch:
					if update.Message == nil {
						continue
					}

					ctx := &Context{
						Update:  update,
						Bot:     bot,
						Context: context.Background(),
					}

					if !update.Message.IsCommand() {
						r.handle(r.cfg.OnText, ctx)
						continue
					}

					handler, ok := r.routes[update.Message.Command()]
					if !ok {
						r.handle(r.cfg.OnNotFound, ctx)
						continue
					}

					r.handle(handler, ctx)
				case <-ctx.Done():
					fmt.Println("closing")
					return
				}
			}
		}()
	}
	r.wg.Wait()
	return nil
}

func (r *routerImpl) Close() {
	r.wg.Wait()
}

func NewRouter(cfg *Config) Router {
	return &routerImpl{
		routes:      map[string]HandlerFunc{},
		middlewares: make([]Middleware, 0),
		cfg:         handleConfig(cfg),
	}
}
