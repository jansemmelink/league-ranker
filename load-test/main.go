//go:build !integration

package main

import (
	"flag"
	"fmt"
	"math/rand/v2"
	"os"
	"path/filepath"
	"time"

	"github.com/go-msvc/errors/v2"
	"github.com/go-msvc/logger"
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

	var r ranker.Ranker
	measure("Making ranker", func() error {
		r, err = ranker.RankerFromFile(tmpFilename)
		return err
	})
	defer r.Close()

	measure("Running ranker", r.Process)

	league := r.Teams().League()
	for _, team := range league {
		fmt.Printf("%d. %s, %d pt%s\n", team.Rank, team.Name, team.Points, plural(team.Points))
	}
} //main()

func plural(i int) string {
	if i != 1 {
		return "s"
	}
	return ""
}

func generateTestData(filename string, nrGames, nrTeams int) func() error {
	return func() error {
		f, err := os.Create(filename)
		if err != nil {
			return errors.Wrapf(err, "failed to create file(%s)", filename)
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
		log.Infof("Load test generated %d games for %d teams.", nrGames, nrTeams)
		return nil
	}
}

var log = logger.New().WithLevel(logger.LevelError)

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
