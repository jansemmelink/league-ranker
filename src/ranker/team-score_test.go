package ranker

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTeamScoreFromStringWithValidEntries(t *testing.T) {
	for _, test := range []struct {
		input       string
		expectName  string
		expectScore int
	}{
		{"Abc 1", "Abc", 1},
		{"A-bc 2", "A-bc", 2},
		{"Name with spaces 3", "Name with spaces", 3},
		{"  Name with outer spaces   3", "Name with outer spaces", 3},
	} {
		t.Run(test.input, func(t *testing.T) {
			teams := NewTeams()
			teamScore, err := NewTeamScoreFromString(test.input, teams)
			assert.Nil(t, err)
			assert.NotNil(t, teamScore.Team())
			assert.Equal(t, test.expectName, teamScore.Team().Name())
			assert.Equal(t, test.expectScore, teamScore.Score())
			assert.Equal(t, fmt.Sprintf("%s %d", test.expectName, test.expectScore), teamScore.String())
			assert.Equal(t, 1, len(teams.All()))
		})
	}
}

func TestNewTeamScoreFromStringWithInvalidEntries(t *testing.T) {
	for _, test := range []struct {
		input       string
		expectError string
	}{
		{"Abc,1", "not <name> <score>"},
		{"Abc -1", "not integer >= 0"},
		{"Abc two", "not integer >= 0"},
		//{" 2", "missing team name"},	unreachable because of trim space at start
	} {
		t.Run(test.input, func(t *testing.T) {
			teams := NewTeams()
			teamScore, err := NewTeamScoreFromString(test.input, teams)
			assert.NotNil(t, err)
			assert.Nil(t, teamScore)
			assert.True(t, errors.Is(err, ErrInvalidTeamScore))
			assert.True(t, strings.Contains(err.Error(), test.expectError), "error does not contain "+test.expectError+" in "+err.Error())
		})
	}
}
