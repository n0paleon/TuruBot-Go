package types

type Middleware func(next CommandHandler) CommandHandler
