package ranker_test

import (
	"testing"

	"github.com/jansemmelink/league-ranker/src/ranker"
	"github.com/stretchr/testify/assert"
)

func TestTeams(t *testing.T) {
	teams := ranker.NewTeams()
	assert.Equal(t, 0, len(teams.All()))

	//add teams by fetching them - they are automatically created, without duplicates
	teamNames := []string{"team1", "team2", "team1", "team3"}
	for _, teamName := range teamNames {
		team := teams.GetByName(teamName)
		assert.NotNil(t, team)
		assert.Equal(t, teamName, team.Name())
	}

	//assert All() returns all without duplicates (in any order)
	assert.Equal(t, 3, len(teams.All()))
	exists := map[string]bool{}
	for _, team := range teams.All() {
		exists[team.Name()] = true
	}
	for _, teamName := range teamNames {
		_, found := exists[teamName]
		assert.True(t, found, "team "+teamName+" not found in teams.All()")
	}
}

func TestRankings(t *testing.T) {
	//create teams before the loop so we can check how ranking progress after each game
	teams := ranker.NewTeams()
	for _, test := range []struct {
		line             string
		expPointsByTeam  map[string]int
		expRankingByTeam map[string]int
	}{
		//group games A-B and C-D
		{"A 3, B 2", map[string]int{"A": 3, "B": 0}, map[string]int{"A": 1, "B": 2}},
		{"C 1, D 1", map[string]int{"A": 3, "B": 0, "C": 1, "D": 1}, map[string]int{"A": 1, "B": 4, "C": 2, "D": 2}},
		//play off: best of each: A-C B-D
		{"A 0, C 1", map[string]int{"A": 3, "B": 0, "C": 4, "D": 1}, map[string]int{"A": 2, "B": 4, "C": 1, "D": 3}},
		{"B 0, D 4", map[string]int{"A": 3, "B": 0, "C": 4, "D": 4}, map[string]int{"A": 3, "B": 4, "C": 1, "D": 1}},
		//semi
		{"A 1, B 1", map[string]int{"A": 4, "B": 1, "C": 4, "D": 4}, map[string]int{"A": 1, "B": 4, "C": 1, "D": 1}},
		{"B 0, D 2", map[string]int{"A": 4, "B": 1, "C": 4, "D": 7}, map[string]int{"A": 2, "B": 4, "C": 2, "D": 1}},
	} {
		t.Run(test.line, func(t *testing.T) {
			//add the game to the league
			_, err := ranker.NewGameFromString(test.line, teams)
			assert.Nil(t, err)

			//check expected points
			for teamName, expPoints := range test.expPointsByTeam {
				team := teams.GetByName(teamName)
				assert.Equal(t, expPoints, team.Points(), "incorrect points for "+team.Name())
			}

			//check league status after this game
			league := teams.Rankings()
			for _, ranking := range league {
				if expRank, ok := test.expRankingByTeam[ranking.TeamName]; ok {
					assert.Equal(t, expRank, ranking.Rank, "incorrect ranking for "+ranking.TeamName)
				}
			}

			//log the rank just for fun when running verbose: nice to explain
			t.Logf("----- AFTER GAME: %s -----", test.line)
			for _, ranking := range league {
				t.Logf("  %d. %s %d pts", ranking.Rank, ranking.TeamName, ranking.TeamPoints)
			}
		})
	}
}
