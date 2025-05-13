package ranker

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/go-msvc/errors/v2"
)

type TeamScore interface {
	Team() Team
	Score() int
	String() string
	Parse(s string, teams Teams) error
}

func TeamScoreFromString(s string, teams Teams) (TeamScore, error) {
	ts := &teamScore{}
	if err := ts.Parse(s, teams); err != nil {
		return nil, errors.Wrapf(err, "failed to parse team score")
	}
	return ts, nil
}

func (ts teamScore) String() string { return fmt.Sprintf("%s %d", ts.team.Name(), ts.score) }

type teamScore struct {
	team  Team
	score int
}

func (ts teamScore) Team() Team { return ts.team }

func (ts teamScore) Score() int { return ts.score }

func (ts *teamScore) Parse(s string, teams Teams) error {
	//remove white space
	s = strings.TrimSpace(s)

	//get the score value from the end of the string, after the last space
	spaceIndex := strings.LastIndex(s, " ")
	if spaceIndex < 0 {
		return errors.Join(ErrInvalidTeamScore, errors.Errorf("not <name> <score>"))
	}
	scoreString := s[spaceIndex+1:] //+1 to skip over the space separator before the score
	if i64, err := strconv.ParseInt(scoreString, 10, 64); err != nil || i64 < 0 {
		return errors.Join(ErrInvalidTeamScore, errors.Errorf("not integer >= 0"))
	} else {
		ts.score = int(i64)
	}

	//team name is everything before, less any outer space (e.g. double space before the score value)
	teamName := strings.TrimRight(s[:spaceIndex], " ")

	//no need to check: unreachable because spaces trimmed at the start :-)
	// if teamName == "" {
	// 	return errors.Join(ErrInvalidTeamScore, errors.Errorf("missing team name"))
	// }

	ts.team = teams.GetByName(teamName)
	return nil
}

var (
	ErrInvalidTeamScore = errors.New("invalid team score")
)
