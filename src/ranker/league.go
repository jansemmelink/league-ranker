package ranker

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

type League interface {
	Teams() Teams
	Games() Games
}

func NewLeague() League {
	return newLeague()
}

func newLeague() *league {
	return &league{
		teams: NewTeams(),
		games: NewGames(),
	}
}

func NewLeagueFromFile(filename string) (League, error) {
	r := newLeague()

	//without a filename, default to stdin
	file := os.Stdin
	if filename != "" {
		//open the specified file
		var err error
		file, err = os.Open(filename)
		if err != nil {
			return nil, errors.Join(err, fmt.Errorf("failed to open input file \"%s\"", filename))
		}
		defer file.Close()
	}

	// The bufio.Scanner in Go by default uses bufio.ScanLines as its splitting function.
	// This function defines a line as a sequence of characters followed by an optional carriage return
	// and a mandatory newline (\r?\n).	scanner := bufio.NewScanner(r.file)
	// This should work for both windows and linux file formats
	// Also not expecting long lines, so should be ok to work with default line buffer size
	scanner := bufio.NewScanner(file)
	lineNr := 0
	for scanner.Scan() {
		lineNr++
		line := scanner.Text()

		//decode the line
		if g, err := NewGameFromString(line, r.teams); err != nil {
			return nil, errors.Join(err, fmt.Errorf("failed to parse line %d", lineNr))
		} else {
			r.games.Add(g)
		}
	}
	return r, nil
} //LeagueFromFile()

type league struct {
	teams Teams
	games Games
}

func (r league) Teams() Teams { return r.teams }
func (r league) Games() Games { return r.games }
