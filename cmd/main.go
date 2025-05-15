package main

import (
	"flag"
	"fmt"

	"github.com/jansemmelink/league-ranker/src/ranker"
)

func main() {
	inFilenameFlag := flag.String("f", "", "Input file to read (default: stdin)")
	flag.Parse()

	r, err := ranker.NewLeagueFromFile(*inFilenameFlag)
	if err != nil {
		panic(fmt.Sprintf("failed to start: %+v", err))
	}

	rankings := r.Teams().Rankings()
	for _, ranking := range rankings {
		fmt.Println(ranking.String())
	}
} //main()
