package ranker

import (
	"bufio"
	"io"
	"os"

	"github.com/go-msvc/errors/v2"
	"github.com/go-msvc/logger"
)

type Ranker interface {
	io.Closer
	Teams() Teams
	Process() error
}

func RankerFromFile(filename string) (Ranker, error) {
	r := &ranker{
		teams: NewTeams(),
	}

	//without a filename, default to stdin
	if filename == "" {
		r.file = os.Stdin
		r.mustClose = false
		return r, nil
	}

	//open the specified file
	var err error
	r.file, err = os.Open(filename)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open input file \"%s\"", filename)
	}

	log.Debugf("Opened named file %s", filename)
	r.mustClose = true
	return r, nil
} //RankerFromFile()

type ranker struct {
	file      *os.File
	mustClose bool //false when reading stdin
	teams     Teams
}

func (r ranker) Teams() Teams { return r.teams }

func (r *ranker) Close() error {
	if r.mustClose && r.file != nil {
		if err := r.file.Close(); err != nil {
			return errors.Wrapf(err, "failed to close file")
		}
		log.Debugf("Closed named file")
		r.mustClose = false
	}
	r.file = nil
	return nil
} //ranker.Close()

func (r *ranker) Process() error {
	if r.file == nil {
		return errors.Errorf("processing file is already closed")
	}

	// The bufio.Scanner in Go by default uses bufio.ScanLines as its splitting function.
	// This function defines a line as a sequence of characters followed by an optional carriage return
	// and a mandatory newline (\r?\n).	scanner := bufio.NewScanner(r.file)
	// This should work for both windows and linux file formats
	// Also not expecting long lines, so should be ok to work with default line buffer size
	scanner := bufio.NewScanner(r.file)
	lineNr := 0
	for scanner.Scan() {
		lineNr++
		line := scanner.Text()
		log.Debugf("Read line: %s", line)

		//decode the line
		var g game
		if err := g.Parse(line, r.teams); err != nil {
			return errors.Wrapf(err, "failed to parse line %d", lineNr)
		}

		log.Debugf("Line %d: %s", lineNr, g)
	}
	log.Debugf("Successfully parsed %d lines.", lineNr)

	//close after processing
	if err := r.Close(); err != nil {
		return errors.Wrapf(err, "failed to close after processing")
	}
	return nil
} //ranker.Process()

// change to LevelDebug to get debug output from ranker when you run this locally
// if changed to live service, remove this and use context logger
var log = logger.New().WithLevel(logger.LevelError)
