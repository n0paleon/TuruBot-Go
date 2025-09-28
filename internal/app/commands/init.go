package commands

import "TuruBot-Go/internal/app/router"

type Command struct {
	router *router.Router
}

func Init(r *router.Router) *Command {
	return &Command{r}
}
