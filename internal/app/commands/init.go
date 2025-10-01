package commands

import (
	"TuruBot-Go/internal/app/router"
	"TuruBot-Go/internal/port"
)

type Command struct {
	router           *router.Router
	storageAdapter   port.StorageProvider
	memeCraftAdapter port.MemeCraft
}

func Init(r *router.Router, storageAdapter port.StorageProvider, memeCraftAdapter port.MemeCraft) *Command {
	return &Command{
		router:           r,
		storageAdapter:   storageAdapter,
		memeCraftAdapter: memeCraftAdapter,
	}
}
