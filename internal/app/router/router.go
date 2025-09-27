package router

import (
	"TuruBot-Go/internal/app/types"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"os"
	"reflect"
	"runtime"
	"strings"
)

type Route struct {
	handler     types.CommandHandler
	cmd         string
	description string
	aliases     []string
	middlewares []types.Middleware
	router      *Router
}

// Set cmd (command name) dan update registry
func (r *Route) Cmd(cmd string) *Route {
	prevCmd := r.cmd
	r.cmd = strings.ToLower(cmd)
	r.router.routes[r.cmd] = r.handler

	if prevCmd != "" && prevCmd != r.cmd {
		delete(r.router.routes, prevCmd)
		for k, v := range r.router.aliases {
			if v == prevCmd {
				r.router.aliases[k] = r.cmd
			}
		}
	}

	for _, a := range r.aliases {
		r.router.aliases[strings.ToLower(a)] = r.cmd
	}

	return r
}

func (r *Route) Description(desc string) *Route {
	r.description = desc
	return r
}

func (r *Route) Aliases(a ...string) *Route {
	r.aliases = append(r.aliases, a...)
	for _, alias := range a {
		r.router.aliases[strings.ToLower(alias)] = strings.ToLower(r.cmd)
	}
	return r
}

func (r *Route) Use(m ...types.Middleware) *Route {
	r.middlewares = append(r.middlewares, m...)
	return r
}

type Router struct {
	routes      map[string]types.CommandHandler
	aliases     map[string]string
	middlewares []types.Middleware // global middleware
	allRoutes   []*Route           // semua Route dengan metadata
}

func New() *Router {
	return &Router{
		routes:  make(map[string]types.CommandHandler),
		aliases: make(map[string]string),
	}
}

func (r *Router) Use(m ...types.Middleware) {
	r.middlewares = append(r.middlewares, m...)
}

func (r *Router) Handle(handler types.CommandHandler) *Route {
	route := &Route{
		handler: handler,
		router:  r,
	}
	r.allRoutes = append(r.allRoutes, route)
	return route
}

func (r *Router) Exec(cmd string, ctx *types.BotContext) error {
	key := strings.ToLower(cmd)
	if main, ok := r.aliases[key]; ok {
		key = main
	}

	handler, ok := r.routes[key]

	var routeMiddlewares []types.Middleware

	if ok {
		for _, route := range r.allRoutes {
			if route.cmd == key {
				routeMiddlewares = route.middlewares
				break
			}
		}
	} else {
		handler = func(ctx *types.BotContext) error {
			// command tidak ditemukan, jangan return error supaya bot diam
			return nil
		}
	}

	allMiddleware := append(append([]types.Middleware(nil), r.middlewares...), routeMiddlewares...)
	wrapped := chain(handler, allMiddleware...)
	return wrapped(ctx)
}

// getFuncName get func name as string
func getFuncName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

// PrintRoutes: print all routes to console
func (r *Router) PrintRoutes() {
	header := []string{"Type", "Command", "Alias", "Description", "Handler/Middleware"}
	var data [][]string

	totalCommands := 0
	totalMiddleware := 0

	for _, route := range r.allRoutes {
		aliases := "<none>"
		if len(route.aliases) > 0 {
			aliases = strings.Join(route.aliases, ", ")
		}

		// Command handler row
		data = append(data, []string{
			"Command",
			route.cmd,
			aliases,
			route.description,
			getFuncName(route.handler),
		})
		totalCommands++

		// Middleware per route
		for _, mw := range route.middlewares {
			data = append(data, []string{
				"Middleware",
				route.cmd,
				"<Not a Command>",
				"",
				getFuncName(mw),
			})
			totalMiddleware++
		}
	}

	// Middleware global
	for _, mw := range r.middlewares {
		data = append(data, []string{
			"Middleware",
			"(global)",
			"<Not a Command>",
			"",
			getFuncName(mw),
		})
		totalMiddleware++
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.Header(header)
	_ = table.Bulk(data)

	// Footer, total
	table.Footer([]string{
		"Total",
		fmt.Sprintf("%d Command(s)", totalCommands),
		fmt.Sprintf("%d Middleware(s)", totalMiddleware),
	})

	_ = table.Render()
}

func chain(h types.CommandHandler, m ...types.Middleware) types.CommandHandler {
	for i := len(m) - 1; i >= 0; i-- {
		h = m[i](h)
	}
	return h
}
