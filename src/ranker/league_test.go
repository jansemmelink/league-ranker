package ranker

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRankerFromStdin(t *testing.T) {
	//mock stdin:
	originalStdin := os.Stdin
	r, w, err := os.Pipe()
	assert.Nil(t, err, "failed to create pipe")
	os.Stdin = r
	defer func() {
		os.Stdin = originalStdin
	}()

	//write lines to stdin to parse by the ranker to create 4 different teams
	//note: this blocks in the pipe until stdin is read... and close when done
	//signaling EOF to the parser
	go func() {
		for _, line := range []string{
			"A 1,B 2",
			"C 3, D 4",
			"A 0, C 1",
		} {
			_, err := w.Write([]byte(line + "\n"))
			assert.Nil(t, err)
		}
		w.Close()
	}()

	//create rankings from stdin
	league, err := NewLeagueFromFile("")
	assert.Nil(t, err)
	assert.NotNil(t, league)
	assert.Equal(t, 4, len(league.Teams().All()))
}

func TestRankerFromNamedFile(t *testing.T) {
	for _, test := range []struct {
		name            string
		lines           []string
		expectedError   string //"" for none
		expectedNrTeams int    //checked if no error
	}{
		{"valid contents 4 teams", []string{"A 1,B 2", "C 3, D 4", "A 0, C 1"}, "", 4},
		{"error on line 1", []string{"A 1;B 2", "C 3, D 4", "A 0, C 1"}, "failed to parse line 1", -1},
		{"error on line 2", []string{"A 1,B 2", "C 3;;; D 4", "A 0, C 1"}, "failed to parse line 2", -1},
		{"error on line 3", []string{"A 1,B 2", "C 3, D 4", "A 0 C 1"}, "failed to parse line 3", -1},
	} {
		t.Run(test.name, func(t *testing.T) {
			//make temp test file
			tmpFile, err := os.CreateTemp("", "test-ranger.*.txt")
			assert.Nil(t, err, "failed to create temp file")
			defer func() {
				os.Remove(tmpFile.Name()) // Clean up the file when done
			}()
			t.Logf("Created temporary file: %s", tmpFile.Name())

			// Get absolute path
			absPath, err := filepath.Abs(tmpFile.Name())
			assert.Nil(t, err, "failed to get absolute filename")
			t.Logf("Absolute path: %s", absPath)

			//write test lines to named file to parse by the league to create 4 different teams
			for _, line := range test.lines {
				_, err := tmpFile.Write([]byte(line + "\n"))
				assert.Nil(t, err, "failed to write to test file")
				t.Logf("Written a line")
			}
			tmpFile.Close()

			//create rankings from that named file - should work for all these tests
			league, err := NewLeagueFromFile(absPath)

			//assert it worked (not checking complete ranking logic - that is tested in other tests)
			if test.expectedError == "" {
				//expecting process to succeed (i.e. valid file contents)
				assert.Nil(t, err, "failed to create league")
				assert.NotNil(t, league, "did not create league")
				assert.Equal(t, 4, len(league.Teams().All()), "not correct nr of teams after processing the file")
			} else {
				//expecting process to fail (i.e. invalid file contents)
				assert.NotNil(t, err, "create league did not fail as expected")
				assert.Nil(t, league, "created league when expected error")
				assert.True(t, strings.Contains(err.Error(), test.expectedError), "expected error "+test.expectedError+" not found in "+err.Error())
			}
		}) //test func
	} //for each test
}
