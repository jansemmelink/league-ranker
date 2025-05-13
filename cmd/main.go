package main

import (
	"flag"
	"fmt"

	"github.com/jansemmelink/league-ranker/src/ranker"
)

func main() {
	inFilenameFlag := flag.String("f", "", "Input file to read (default: stdin)")
	flag.Parse()

	r, err := ranker.RankerFromFile(*inFilenameFlag)
	if err != nil {
		panic(fmt.Sprintf("failed to start: %+v", err))
	}
	defer r.Close()
	if err := r.Process(); err != nil {
		panic(fmt.Sprintf("failed to process: %+v", err))
	}

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
