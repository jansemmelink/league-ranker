package ranker

import (
	"errors"
	"fmt"
	"strings"
)

// Game represents the score of the two teams that played each other
type Game interface {
	String() string
	TeamScore(index int) TeamScore //index = 0..1
}

func NewGameFromString(csv string, teams Teams) (Game, error) {
	columns := strings.SplitN(csv, ",", 3) //split on the ',' into 2 columns (allow 3 to enable detection when more)
	if len(columns) != 2 {
		return nil, errors.Join(ErrInvalidGame, errors.New("not \"<team 1 score>, <team 2 score>\""))
	}

	//parse both columns
	teamScores := [2]TeamScore{}
	for col, teamScore := range columns {
		ts, err := NewTeamScoreFromString(teamScore, teams)
		if err != nil {
			return nil, errors.Join(ErrInvalidGame, err, fmt.Errorf("error in column %d", col+1))
		}
		teamScores[col] = ts
	}
	return newGame(teamScores[0], teamScores[1])
} //NewGameFromString()

func NewGame(score1, score2 TeamScore) (Game, error) {
	if score1 == nil || score2 == nil {
		return nil, ErrInvalidGame
	}
	return newGame(score1, score2)
}

func newGame(score1, score2 TeamScore) (Game, error) {
	g := &game{
		teamScores: [2]TeamScore{score1, score2},
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
	return g, nil
} //NewGame()

type game struct {
	teamScores [2]TeamScore
}

func (g game) String() string { return g.teamScores[0].String() + ", " + g.teamScores[1].String() }

func (g game) TeamScore(index int) TeamScore {
	if index >= 0 && index < 2 {
		return g.teamScores[index]
	}
	return nil
}

// expected CSV with two team scores
var (
	ErrInvalidGame = errors.New("invalid game")
)
