package ranker

import "fmt"

type Team interface {
	Name() string
	String() string
	Award(points int)
	Points() int
}

func NewTeam(name string) Team {
	return &team{name: name, points: 0}
}

type team struct {
	name   string
	points int
}

func (t *team) Award(points int) {
	t.points += points
}

func (t team) Name() string   { return t.name }
func (t team) Points() int    { return t.points }
func (t team) String() string { return fmt.Sprintf("%s %d", t.name, t.points) }
