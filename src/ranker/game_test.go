package ranker_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/jansemmelink/league-ranker/src/ranker"
	"github.com/stretchr/testify/assert"
)

func TestGameValidLine(t *testing.T) {
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
			g, err := ranker.GameFromString(test.line, teams)
			assert.Nil(t, err)
			assert.NotNil(t, g)
			assert.Equal(t, test.expTeamName1, g.TeamScore(0).Team().Name(), "wrong team[0] name")
			assert.Equal(t, test.expTeamScore1, g.TeamScore(0).Score(), "wrong team[0].score")
			assert.Equal(t, test.expTeamName2, g.TeamScore(1).Team().Name(), "wrong team[1] name")
			assert.Equal(t, test.expTeamScore2, g.TeamScore(1).Score(), "wrong team[1].score")
			assert.Equal(t, fmt.Sprintf("%s %d, %s %d", test.expTeamName1, test.expTeamScore1, test.expTeamName2, test.expTeamScore2), g.String(), "wrong game string")
			assert.True(t, true)

			//out of range index does not return nil:
			assert.NotNil(t, g.TeamScore(-1))                  //out of range index
			assert.NotNil(t, g.TeamScore(-1).Team())           //out of range index
			assert.Equal(t, "", g.TeamScore(-1).Team().Name()) //out of range index
			assert.Equal(t, 0, g.TeamScore(-1).Score())        //out of range index
			assert.NotNil(t, g.TeamScore(2))                   //out of range index
			assert.NotNil(t, g.TeamScore(2).Team())            //out of range index
			assert.Equal(t, "", g.TeamScore(2).Team().Name())  //out of range index
			assert.Equal(t, 0, g.TeamScore(2).Score())         //out of range index
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
			g, err := ranker.GameFromString(test.line, teams)
			assert.NotNil(t, err)
			assert.Nil(t, g)
			assert.True(t, strings.Contains(err.Error(), test.expectedError), "error does not contain "+test.expectedError+" in "+err.Error())
			assert.True(t, errors.Is(err, ranker.ErrInvalidGame), "error is not invalid game")
		})
	}
}
