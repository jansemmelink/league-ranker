package ranker_test

import (
	"fmt"
	"testing"

	"github.com/jansemmelink/league-ranker/src/ranker"
	"github.com/stretchr/testify/assert"
)

func TestTeam(t *testing.T) {
	//new team without points
	const teamName = "some-name"
	team := ranker.NewTeam(teamName)
	assert.NotNil(t, team)
	assert.Equal(t, teamName, team.Name())
	assert.Equal(t, 0, team.Points())
	assert.Equal(t, teamName+" 0", team.String())

	//award first points
	const firstPoints = 3
	team.Award(firstPoints)
	assert.Equal(t, teamName, team.Name())
	assert.Equal(t, firstPoints, team.Points())
	assert.Equal(t, fmt.Sprintf("%s %d", teamName, firstPoints), team.String())

	//awart more points
	const morePoints = 9
	team.Award(morePoints)
	assert.Equal(t, teamName, team.Name())
	assert.Equal(t, firstPoints+morePoints, team.Points())
	assert.Equal(t, fmt.Sprintf("%s %d", teamName, firstPoints+morePoints), team.String())
}
