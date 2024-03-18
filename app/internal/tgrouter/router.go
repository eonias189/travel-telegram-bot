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
	On(command string, handler HandlerFunc)
	Run(ctx context.Context, token string) error
	Close()
	Use(middleware Middleware) Router
}

type routerImpl struct {
	routes      map[string]HandlerFunc
	middlewares []Middleware
	ctx         context.Context
	cfg         *Config
	wg          sync.WaitGroup
	mu          sync.RWMutex
}

func (r *routerImpl) On(command string, handler HandlerFunc) {
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
	err := r.withMWares(handler)(ctx)
	if err != nil {
		r.cfg.OnError(err)
	}
}

func (r *routerImpl) newContext(update tgbotapi.Update, bot *tgbotapi.BotAPI) *Context {
	return &Context{
		Update: update,
		Bot:    bot,
		ctx:    r.ctx,
		setContext: func(ctx context.Context) {
			r.mu.RLock()
			defer r.mu.RLock()
			r.ctx = ctx
		},
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

					ctx := r.newContext(update, bot)

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
		ctx:         context.Background(),
	}
}
