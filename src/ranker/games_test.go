package ranker_test

import (
	"testing"

	"github.com/jansemmelink/league-ranker/src/ranker"
	"github.com/stretchr/testify/assert"
)

func TestGames(t *testing.T) {
	teams := ranker.NewTeams()
	assert.Equal(t, 0, len(teams.All()))
	games := ranker.NewGames()
	assert.Equal(t, 0, len(games.All()))

	//add 1st game
	g1 := must(ranker.NewGame(
		must(ranker.NewTeamScore(teams.GetByName("A"), 1)),
		must(ranker.NewTeamScore(teams.GetByName("B"), 1)),
	))
	games.Add(g1)

	//add 2nd game
	g2 := must(ranker.NewGame(
		must(ranker.NewTeamScore(teams.GetByName("C"), 3)),
		must(ranker.NewTeamScore(teams.GetByName("D"), 0)),
	))
	games.Add(g2)

	allGames := games.All()
	assert.Equal(t, 2, len(allGames))
	assert.Equal(t, g1.String(), allGames[0].String())
	assert.Equal(t, g2.String(), allGames[1].String())
}

func must[Type any](value Type, err error) Type {
	if err != nil {
		panic(err)
	}
	return value
}
