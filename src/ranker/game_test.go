package ranker_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/jansemmelink/league-ranker/src/ranker"
	"github.com/stretchr/testify/assert"
)

func TestNewGame(t *testing.T) {
	teams := ranker.NewTeams()

	const (
		teamName1  = "A"
		teamName2  = "B"
		teamScore1 = 4
		teamScore2 = 5
	)
	g := must(ranker.NewGame(
		must(ranker.NewTeamScore(teams.GetByName(teamName1), teamScore1)),
		must(ranker.NewTeamScore(teams.GetByName(teamName2), teamScore2)),
	))
	assert.Equal(t, teamName1, g.TeamScore(0).Team().Name())
	assert.Equal(t, teamScore1, g.TeamScore(0).Score())
	assert.Equal(t, teamName2, g.TeamScore(1).Team().Name())
	assert.Equal(t, teamScore2, g.TeamScore(1).Score())
	assert.Nil(t, g.TeamScore(-1), "out of range index expected to return nil")
	assert.Nil(t, g.TeamScore(2), "out of range index expected to return nil")
}

func TestNewGameFailure(t *testing.T) {
	teams := ranker.NewTeams()
	teamScore, _ := ranker.NewTeamScore(teams.GetByName("A"), 1)

	{
		g, err := ranker.NewGame(teamScore, nil)
		assert.NotNil(t, err)
		assert.Nil(t, g)
		assert.True(t, errors.Is(err, ranker.ErrInvalidGame))
	}

	{
		g, err := ranker.NewGame(nil, teamScore)
		assert.NotNil(t, err)
		assert.Nil(t, g)
		assert.True(t, errors.Is(err, ranker.ErrInvalidGame))
	}
}

func TestNewGameFromStringValidLine(t *testing.T) {
	for _, test := range []struct {
		name          string
		line          string
		expTeamName1  string
		expTeamName2  string
		expTeamScore1 int
		expTeamScore2 int
	}{
		{"default", "Team A 3, Team B 2", "Team A", "Team B", 3, 2},
		{"with spaces, inner space is preserved", "      AAA  XXX    0   ,     BBB   1111    3      ", "AAA  XXX", "BBB   1111", 0, 3},
	} {
		t.Run(test.name, func(t *testing.T) {
			teams := ranker.NewTeams()
			g, err := ranker.NewGameFromString(test.line, teams)
			assert.Nil(t, err)
			assert.NotNil(t, g)
			assert.Equal(t, test.expTeamName1, g.TeamScore(0).Team().Name(), "wrong team[0] name")
			assert.Equal(t, test.expTeamScore1, g.TeamScore(0).Score(), "wrong team[0].score")
			assert.Equal(t, test.expTeamName2, g.TeamScore(1).Team().Name(), "wrong team[1] name")
			assert.Equal(t, test.expTeamScore2, g.TeamScore(1).Score(), "wrong team[1].score")
			assert.Equal(t, fmt.Sprintf("%s %d, %s %d", test.expTeamName1, test.expTeamScore1, test.expTeamName2, test.expTeamScore2), g.String(), "wrong game string")
			assert.True(t, true)
		})
	}
}

func TestGameInvalidLine(t *testing.T) {
	for _, test := range []struct {
		name          string
		line          string
		expectedError string
	}{
		{"fewer than 2 CSV columns", "Team A 1; Team B 2", ""},
		{"more than 2 CSV columns", "Team A 1, Team B 2, Team C 3", ""},
		{"error in 1st col", "Team A one, Team B 2", ""},
		{"error in 2nd col", "Team A 1, Team B two", ""},
	} {
		t.Run(test.name, func(t *testing.T) {
			teams := ranker.NewTeams()
			g, err := ranker.NewGameFromString(test.line, teams)
			assert.NotNil(t, err)
			assert.Nil(t, g)
			assert.True(t, strings.Contains(err.Error(), test.expectedError), "error does not contain "+test.expectedError+" in "+err.Error())
			assert.True(t, errors.Is(err, ranker.ErrInvalidGame), "error is not invalid game")
		})
	}
}
