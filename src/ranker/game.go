package ranker

import (
	"strings"

	"github.com/go-msvc/errors/v2"
)

type Game interface {
	String() string
	Parse(string, Teams) error
	TeamScore(index int) TeamScore
}

func GameFromString(s string, teams Teams) (Game, error) {
	g := &game{}
	if err := g.Parse(s, teams); err != nil {
		return nil, errors.Join(ErrInvalidGame, errors.Wrapf(err, "failed to parse game"))
	}
	return g, nil
}

type game struct {
	teamScores [2]TeamScore
}

func (g game) String() string { return g.teamScores[0].String() + ", " + g.teamScores[1].String() }

func (g game) TeamScore(index int) TeamScore {
	if index >= 0 && index < 2 {
		return g.teamScores[index]
	}
	return &teamScore{team: &team{}} //should not get here - just safer for caller
}

// expected CSV with two team scores
func (g *game) Parse(csv string, teams Teams) error {
	columns := strings.SplitN(csv, ",", 3) //split on the ',' into 2 columns (allow 3 to enable detection when more)
	if len(columns) != 2 {
		return errors.Join(ErrInvalidGame, errors.Errorf("line format is not \"<team 1 score>, <team 2 score>\""))
	}
	//parse both columns
	for col, teamScore := range columns {
		ts, err := TeamScoreFromString(teamScore, teams)
		if err != nil {
			return errors.Join(ErrInvalidGame, errors.Wrapf(err, "error in column %d", col+1))
		}
		log.Debugf("col[%d]: %s", col, ts)
		g.teamScores[col] = ts
	}

	//game ranking rules (for each team)
	if g.teamScores[0].Score() == g.teamScores[1].Score() {
		//a draw (tie) is worth 1 point for each team
		g.teamScores[0].Team().Award(1)
		g.teamScores[1].Team().Award(1)
	} else if g.teamScores[0].Score() > g.teamScores[1].Score() {
		//	- a win is worth 3 points
		//	- a loss is worth 0 points
		g.teamScores[0].Team().Award(3) //winner in column[0]
	} else {
		//	- a win is worth 3 points
		//	- a loss is worth 0 points
		g.teamScores[1].Team().Award(3) //winner in column[1]
	}
	return nil
}

var (
	ErrInvalidGame = errors.New("invalid game")
)
