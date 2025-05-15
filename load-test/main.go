//go:build !integration

package main

import (
	"errors"
	"flag"
	"fmt"
	"math/rand/v2"
	"os"
	"path/filepath"
	"time"

	"github.com/jansemmelink/league-ranker/src/ranker"
)

func main() {
	loadTestNrGamesFlag := flag.Int("nr-games", 1000, "Write N records to -f <file> then run to measure speed")
	loadTestNrTeamsFlag := flag.Int("nr-teams", 12, "Use N different team names for load test")
	flag.Parse()

	//for load test
	if *loadTestNrGamesFlag <= 0 {
		panic("--nr-games <= 0")
	}
	if *loadTestNrTeamsFlag <= 0 {
		panic("--nr-teams <= 0")
	}

	tmpFile, err := os.CreateTemp("", "test-ranger.*.txt")
	if err != nil {
		panic(err)
	}
	defer func() {
		os.Remove(tmpFile.Name())
	}()
	tmpFilename, _ := filepath.Abs(tmpFile.Name())

	measure("Generating test data", generateTestData(tmpFilename, *loadTestNrGamesFlag, *loadTestNrTeamsFlag))

	var league ranker.League
	measure("Loading league", func() error {
		league, err = ranker.NewLeagueFromFile(tmpFilename)
		return err
	})
	rankings := league.Teams().Rankings()
	for _, ranking := range rankings {
		fmt.Println(ranking.String())
	}
} //main()

func generateTestData(filename string, nrGames, nrTeams int) func() error {
	return func() error {
		f, err := os.Create(filename)
		if err != nil {
			return errors.Join(err, fmt.Errorf("failed to create file(%s)", filename))
		}
		defer f.Close()

		for i := 0; i < nrGames; i++ {
			team1 := rand.IntN(nrTeams)
			team2 := rand.IntN(nrTeams)
			for team2 == team1 {
				team2 = rand.IntN(nrTeams)
			}
			score1 := rand.IntN(5)
			score2 := rand.IntN(5)
			fmt.Fprintf(f, "Team %d %d, Team %d %d\n", team1, score1, team2, score2)
		}
		return nil
	}
}

func measure(what string, fnc func() error) {
	t0 := time.Now()
	fmt.Fprintf(os.Stderr, "%10s %10s %s\n", "", "START", what)
	defer func() {
		dur := time.Since(t0)
		fmt.Fprintf(os.Stderr, "%10v %10s %s\n", dur, "DONE", what)
	}()
	err := fnc()
	if err != nil {
		panic(fmt.Sprintf("failed to %s: %+v", what, err))
	}
}
