# Development

# Design

The data is modelled as:
* Ranker: Read file into memory as Games and Teams
* Game: Consists of 2 teams, each with a score for the game, modeled as TeamScore
* Games: A collection of chronological games played
* TeamScore: Links a team to a score in a game
* Team: Represents a named team with accumulation of points
* Teams: Manages the collection of teams

All the models are defined as an interface and an implementation to protect (scope) the internals from the users of the model.

Ranker can be created and then add games as they are played and see the league change, or process all games at once from a file.

# Unit Test Coverage
Check test coverage with:
```
go test ./... -coverprofile=./coverage.out
go tool cover -html coverage.out
```
That should open in your browser then use dropdown to check coverage in each file.

At time of writing it is 100% for all src/ranker modules except 2 of 3 lines of unexpected file errors that are handled and not tested. Also cmd/main.go is not tested with unit testing, but can be verified on the command line with a script:
```
go run ./cmd -f ./test/data/minimum/input.txt > /tmp/a
diff /tmp/a test/data/minimum/expected-output.txt
[ $? -ne 0 ] && echo "ERROR ************************* Differenced Found" || echo "DONE: ALL GOOD"
rm -f /tmp/a
```
Expected Output:
```
DONE: ALL GOOD
```

# Load Testing
Run a load test with generated input:
```
go run . --nr-teams 100 --nr-games 1000000
                START Generating test data
1.57304575s       DONE Generating test data
                START Making ranker
 331.584µs       DONE Making ranker
                START Running ranker
474.649917ms       DONE Running ranker
1. Team 89, 28560 pts
2. Team 27, 28540 pts
3. Team 24, 28534 pts
4. Team 34, 28467 pts
5. Team 94, 28459 pts
...
96. Team 69, 27439 pts
97. Team 2, 27400 pts
98. Team 49, 27328 pts
99. Team 51, 27266 pts
100. Team 4, 27165 pts
```
