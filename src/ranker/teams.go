package ranker

import (
	"sort"
	"sync"
)

type Teams interface {
	GetByName(name string) Team //create if not found
	All() []Team                //unsorted list of team
	Rankings() []TeamRanking
}

func NewTeams() Teams {
	return &teams{
		byName: map[string]Team{},
	}
}

type teams struct {
	sync.Mutex
	byName map[string]Team
}

func (teams *teams) GetByName(name string) Team {
	teams.Lock()
	defer teams.Unlock()
	team, ok := teams.byName[name]
	if !ok {
		team = NewTeam(name)
		teams.byName[name] = team
	}
	return team
}

func (teams *teams) All() []Team {
	teams.Lock()
	defer teams.Unlock() //would be more optimal to use a read-only lock here...

	//make an array of all teams
	result := []Team{}
	for _, team := range teams.byName {
		result = append(result, team)
	}
	return result
}

// League Rules:
//   - Output is ordered on team points (descending)
//   - If two or more teams have the same number of points:
//     they should have the same rank and be
//     printed in alphabetical order (ascending)
func (teams *teams) Rankings() []TeamRanking {
	//make a copy and sort the teams in the league
	allTeams := teams.All()
	sort.Slice(allTeams, func(i, j int) bool {
		if allTeams[i].Points() != allTeams[j].Points() {
			return allTeams[i].Points() > allTeams[j].Points()
		}
		return allTeams[i].Name() < allTeams[j].Name()
	})

	//award team ranks:
	//- rank is the same for all teams with the same score
	//- rank is skipped if multiple teams have the previous rank
	//		this is wierd... but the next team with a different nr of points gets the N't rank
	//		so if 3 teams was ranked 3rd, the no team will have rank 4 or 5. The next team will be rank 6
	//		e.g.:
	// 1. AAA, 6 pts
	// 2. BBB, 5 pts
	// 3. CCC, 1 pt
	// 3. DDD, 1 pt
	// 3. EEE, 1 pts
	// 6. FFF, 0 pts (rank 4 and 5 skipped)
	teamRankings := make([]TeamRanking, len(allTeams))
	for pos, team := range allTeams {
		teamRanking := TeamRanking{
			TeamName:   team.Name(),
			TeamPoints: team.Points(),
			Rank:       pos + 1,
		}
		if pos > 0 && team.Points() == allTeams[pos-1].Points() {
			//same rank
			teamRanking.Rank = teamRankings[pos-1].Rank
		}

		teamRankings[pos] = teamRanking
	}
	return teamRankings
}
