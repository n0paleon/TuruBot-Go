package commands

import (
	"TuruBot-Go/internal/app/types"
	"fmt"
	"strings"
)

func (cmd *Command) ShowMenu(ctx *types.BotContext) error {
	routes := cmd.router.GetAll()

	var sb strings.Builder
	for _, route := range routes {
		sb.WriteString(fmt.Sprintf("/%s -> _%s_\n", route.Cmd, route.Description))
		sb.WriteString(fmt.Sprintf("Example: [%s]\n", route.Help))
		sb.WriteString(fmt.Sprintf("Alias: [%s]\n\n", strings.Join(route.Aliases, ", ")))
	}
	sb.WriteString("\n_Thank you for using TuruBot_")

	return ctx.Reply(sb.String())
}
