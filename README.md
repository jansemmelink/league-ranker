# League Ranker
Command-line utility to rank teams after league games.

# League Rules
* Each game has two teams with a final score, e.g.:
`Falcons 5, Eagles 2`
* Points are awarded to both teams for each game:
    * 3 points for a win
    * 1 point each for a tie (even a 0-0 tie)
    * 0 points for a loss
* All points are summed to determine league ranking
* Teams with the same total points get the same ranking
    * They are reported in alphabetical order

Example of input:
```
Lions 3, Snakes 3
Tarantulas 1, FC Awesome 0
Lions 1, FC Awesome 1
Tarantulas 3, Snakes 1
Lions 4, Grouches 0
```

Example of ranking:
```
1. Tarantulas, 6 pts
2. Lions, 5 pts
3. FC Awesome, 1 pt
3. Snakes, 1 pt
5. Grouches, 0 pts
```

> Note: in the above two teams ranked 3rd because they have the same number of points. Therefore none are ranked 4th.

# Usage
The utility parses text input from stdin and writes to stdout:
```
cat test/data/minimum/input.txt | go run ./cmd
```
It can also process a named file:
```
go run ./cmd -f ./test/data/minimum/input.txt
```

To show command-line options:
```
cd cmd
go build
./cmd --help
Usage of ./cmd:
  -f string
    	Input file to read (default: stdin)
```

# Design Assumptions
* Game leagues are generally limited data sets on a season of games, i.e. no big data and need for streaming large batches of data. Still the code would perform reasonably on thousands and even millions of games, e.g. if this is used for online games. Without streaming however, it will need all games in memory at once.
* Used interfaces and private structs to make mocking possible at each level
* Defined a constructor from string for each item used for parsing, since we're not processing database content.
* Considered using mocks, but since there is no external dependencies - did not find much use for using mocks.
