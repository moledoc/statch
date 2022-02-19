package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// TODO: documentation
// NOTE: only one open label at a time is allowed

// logfile is a variable that holds the log file path.
// structure of the file: label start end
var logfile string = ".trakr.csv"
var label string = "all"
var compare time.Time

// trak is a structure that holds each logged item's label, start, end and duration.
type trak struct {
	label    string
	start    time.Time
	end      time.Time
	duration time.Duration
}

func (t trak) Store() string {
	var saveEnd string
	if t.end != compare {
		saveEnd = strconv.FormatInt(t.end.Unix(), 10)
	}
	return fmt.Sprintf("%v,%v,%v\n", t.label, t.start.Unix(), saveEnd)
}

var format string = "%-10v %-30v %-30v %5v"
var header string = fmt.Sprintf(format, "label", "start", "end", "duration")

func (t trak) String() string {
	return fmt.Sprintf(format, t.label, t.start.String(), t.end.String(), t.duration)
}

// logged is a function that reads and parses the contents of the logfile.
func logged(label string) ([]trak, int) {
	var traks []trak
	var openLabel int = -1
	if _, err := os.Stat(logfile); err != nil {
		return traks, openLabel
	}
	f, err := os.Open(logfile)
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(f)
	defer func() {
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}()
	var i int
	for scanner.Scan() {
		contents := strings.Split(scanner.Text(), ",")
		srt, err := strconv.ParseInt(contents[1], 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		srtTime := time.Unix(srt, 0)
		var endTime time.Time
		var duration time.Duration
		if contents[2] != "" {
			end, err := strconv.ParseInt(contents[2], 10, 64)
			if err != nil {
				log.Fatal(err)
			}
			endTime = time.Unix(end, 0)
			duration = endTime.Sub(srtTime)
		}
		if openLabel == -1 && contents[2] == "" {
			openLabel = i
		}
		traks = append(traks, trak{contents[0], srtTime, endTime, duration})
		i++
	}
	return traks, openLabel
}

// help is a function that prints help.
func help() {
	fmt.Println("TODO:")
}

func save(traks *[]trak) {
	f, err := os.OpenFile(logfile, os.O_WRONLY|os.O_TRUNC, 0600)
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}
	for _, elem := range *traks {
		_, err = f.WriteString(elem.Store())
		if err != nil {
			log.Fatal(err)
		}
	}
}

// end is a function that closes the last opened insert for corresponding label.
func end(traks *[]trak, openLabel int) {
	if openLabel != -1 {
		cur := (*traks)[openLabel]
		cur.end = time.Now()
		cur.duration = cur.end.Sub(cur.start)
		(*traks)[openLabel] = cur
		fmt.Printf("Closed '%v'\n", cur.label)

	}
	save(traks)
}

// start is a function that starts a new insert for given label.
// If any previous insert was still open for given label, then that insert gets closed.
func start(label string, traks *[]trak, openLabel int) {
	*traks = append(*traks, trak{label, time.Now(), compare, time.Duration(0)})
	end(traks, openLabel)
	fmt.Printf("Started '%v'\n", label)
}

// trakr [action] (label)

func main() {
	if len(os.Args) > 2 {
		label = os.Args[2]
	}
	traks, openLabel := logged(label)
	switch os.Args[1] {
	case "help":
		help()
	case "show":
		if len(traks) == 0 {
			fmt.Println("Nothing logged yet")
			return
		}
		fmt.Println(header)
		for _, elem := range traks {
			if label == "all" || elem.label == label {
				fmt.Println(elem.String())
			}
		}
	case "start":
		start(label, &traks, openLabel)
	case "end":
		if openLabel == -1 {
			fmt.Println("No trak to close")
			return
		}
		end(&traks, openLabel)
	case "summary":
		log.Fatal("TODO:")
	default:
		fmt.Println("Unknown action")
		help()
	}
	fmt.Printf("--------------------------------------------------------------------------------\nDONE\n")
}
