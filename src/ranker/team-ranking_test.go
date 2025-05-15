package ranker

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTeamRankingString(t *testing.T) {
	for _, test := range []struct {
		name string
		tr   TeamRanking
		exp  string
	}{
		{"zero", TeamRanking{Rank: 4, TeamName: "X", TeamPoints: 0}, "4. X, 0 pts"},
		{"singular", TeamRanking{Rank: 5, TeamName: "X", TeamPoints: 1}, "5. X, 1 pt"},
		{"plural", TeamRanking{Rank: 6, TeamName: "X", TeamPoints: 2}, "6. X, 2 pts"},
	} {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.exp, test.tr.String(), "unexpected string value for TeamRanking")
		})
	}
}
