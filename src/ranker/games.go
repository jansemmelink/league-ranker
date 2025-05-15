package ranker

import (
	"slices"
	"sync"
)

type Games interface {
	Add(Game)
	All() []Game //in order they were added
}

func NewGames() Games {
	return &games{
		all: []Game{},
	}
}

type games struct {
	sync.Mutex
	all []Game
}

func (games *games) Add(game Game) {
	games.Lock()
	defer games.Unlock()
	games.all = append(games.all, game)
}

func (games *games) All() []Game {
	games.Lock()
	defer games.Unlock() //would be more optimal to use a read-only lock here...
	return slices.Clone(games.all)
}
