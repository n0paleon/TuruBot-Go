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
	Cmd         string
	Description string
	Help        string
	Aliases     []string
	middlewares []types.Middleware
	router      *Router
}

// SetCmd set command name and update registry
func (r *Route) SetCmd(cmd string) *Route {
	prevCmd := r.Cmd
	r.Cmd = strings.ToLower(cmd)
	r.router.routes[r.Cmd] = r.handler

	if prevCmd != "" && prevCmd != r.Cmd {
		delete(r.router.routes, prevCmd)
		for k, v := range r.router.aliases {
			if v == prevCmd {
				r.router.aliases[k] = r.Cmd
			}
		}
	}

	for _, a := range r.Aliases {
		r.router.aliases[strings.ToLower(a)] = r.Cmd
	}

	return r
}

// SetDescription set command description
func (r *Route) SetDescription(desc string) *Route {
	r.Description = desc
	return r
}

// SetHelp set command help usage
func (r *Route) SetHelp(help string) *Route {
	r.Help = help
	return r
}

// SetAliases set command aliases
func (r *Route) SetAliases(a ...string) *Route {
	r.Aliases = append(r.Aliases, a...)
	for _, alias := range a {
		r.router.aliases[strings.ToLower(alias)] = strings.ToLower(r.Cmd)
	}
	return r
}

// Use register new middleware for specific route only
func (r *Route) Use(m ...types.Middleware) *Route {
	r.middlewares = append(r.middlewares, m...)
	return r
}

type Router struct {
	routes      map[string]types.CommandHandler
	aliases     map[string]string
	middlewares []types.Middleware // global middleware
	allRoutes   []*Route           // semua route dengan metadata
}

func New() *Router {
	return &Router{
		routes:  make(map[string]types.CommandHandler),
		aliases: make(map[string]string),
	}
}

// Use register middleware as global route
func (r *Router) Use(m ...types.Middleware) {
	r.middlewares = append(r.middlewares, m...)
}

// GetAll return all routes
func (r *Router) GetAll() []*Route {
	routesCopy := make([]*Route, len(r.allRoutes))
	copy(routesCopy, r.allRoutes)
	return routesCopy
}

// Handle register new handler
func (r *Router) Handle(handler types.CommandHandler) *Route {
	route := &Route{
		handler: handler,
		router:  r,
	}
	r.allRoutes = append(r.allRoutes, route)
	return route
}

// Exec execute message based on available routes
func (r *Router) Exec(cmd string, ctx *types.BotContext) error {
	key := strings.ToLower(cmd)
	if main, ok := r.aliases[key]; ok {
		key = main
	}

	var (
		handler          types.CommandHandler
		routeMiddlewares []types.Middleware
		cmdNotFound      bool
	)

	handler, ok := r.routes[key]
	if !ok {
		cmdNotFound = true
		handler = func(ctx *types.BotContext) error {
			return nil
		}
	} else {
		for _, route := range r.allRoutes {
			if route.Cmd == key {
				routeMiddlewares = route.middlewares
				break
			}
		}
	}

	allMiddleware := append(append([]types.Middleware(nil), r.middlewares...), routeMiddlewares...)
	wrapped := chain(handler, allMiddleware...)
	if err := wrapped(ctx); err != nil {
		return err
	}

	if cmdNotFound {
		return fmt.Errorf("COMMAND_NOT_FOUND")
	}
	return nil
}

// PrintRoutes print all routes to console
func (r *Router) PrintRoutes() {
	header := []string{"Type", "Command", "Alias", "Description", "Handler/Middleware"}
	var data [][]string

	totalCommands := 0
	totalMiddleware := 0

	for _, route := range r.allRoutes {
		aliases := "<none>"
		if len(route.Aliases) > 0 {
			aliases = strings.Join(route.Aliases, ", ")
		}

		// Command handler row
		data = append(data, []string{
			"Command",
			route.Cmd,
			aliases,
			route.Description,
			getFuncName(route.handler),
		})
		totalCommands++

		// Middleware per route
		for _, mw := range route.middlewares {
			data = append(data, []string{
				"Middleware",
				route.Cmd,
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

	// footer, total
	table.Footer([]string{
		"Total",
		fmt.Sprintf("%d Command(s)", totalCommands),
		fmt.Sprintf("%d Middleware(s)", totalMiddleware),
	})

	_ = table.Render()
}

// getFuncName return func name as string
func getFuncName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

// chain bootstrap handler and all middlewares
func chain(h types.CommandHandler, m ...types.Middleware) types.CommandHandler {
	for i := len(m) - 1; i >= 0; i-- {
		h = m[i](h)
	}
	return h
}
