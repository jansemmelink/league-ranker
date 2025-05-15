package ranker

import "fmt"

type TeamRanking struct {
	Rank       int
	TeamName   string
	TeamPoints int
}

func (tr TeamRanking) String() string {
	return fmt.Sprintf("%d. %s, %d pt%s", tr.Rank, tr.TeamName, tr.TeamPoints, plural(tr.TeamPoints))
}

func plural(i int) string {
	if i != 1 {
		return "s"
	}
	return ""
}
