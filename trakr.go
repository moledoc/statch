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

// logfile is a variable that holds the log file path.
// structure of the file: label start end
var logfile string = ".trakr.csv"
var label string = "general"
var compare time.Time

// trak is a structure that holds each logged item's label, start, end and duration.
type trak struct {
	label    string
	start    time.Time
	end      time.Time
	duration time.Duration
}

// help is a function that prints help.
func help() {
	fmt.Println("TODO:")
}

// start is a function that starts a new insert for given label.
// IDEA: // If any previous insert was still open for given label, then that insert get's a closed value.
func start(label string) {
	srt := time.Now()
	line := fmt.Sprintf("%v,%v,%v\n", label, srt.Unix(), "")
	f, err := os.OpenFile(logfile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}
	_, err = f.WriteString(line)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Added start time '%v' to label '%v'.\n", srt.String(), label)
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
		if openLabel == -1 && contents[2] == "" && label == contents[0] {
			openLabel = i
		}
		traks = append(traks, trak{contents[0], srtTime, endTime, duration})
		i++
	}
	return traks, openLabel
}

// end is a function that closes the last opened insert for corresponding label.
func end(traks *[]trak, openLabel int) {
	cur := (*traks)[openLabel]
	cur.end = time.Now()
	cur.duration = cur.end.Sub(cur.start)
	(*traks)[openLabel] = cur
	fmt.Printf("Added end time '%v' to label '%v' with start time '%v'.\n", cur.end.String(), cur.label, cur.start.String())

	f, err := os.OpenFile(logfile, os.O_WRONLY|os.O_TRUNC, 0600)
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}
	for _, elem := range *traks {
		var saveEnd string
		if elem.end != compare {
			saveEnd = strconv.FormatInt(elem.end.Unix(), 10)
		}
		line := fmt.Sprintf("%v,%v,%v\n", elem.label, elem.start.Unix(), saveEnd)
		_, err = f.WriteString(line)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// trakr [action] [subaction|label]

func main() {
	if len(os.Args) > 2 {
		label = os.Args[2]
	}
	traks, openLabel := logged(label)
	switch os.Args[1] {
	case "help":
		help()
	case "start":
		if label == "open" {
			log.Fatal("Label 'open' not allowed")
		}
		start(label)
	case "show":
		if len(traks) == 0 {
			fmt.Println("Nothing logged yet")
		}
		for _, elem := range traks {
			if label == "all" || (label == "open" && elem.end == compare) || elem.label == label {
				fmt.Printf("%-10v %-30v %-30v %5v\n", elem.label, elem.start.String(), elem.end.String(), elem.duration)
			}
		}
	case "end":
		end(&traks, openLabel)
	default:
		fmt.Println("Unknown action")
		help()
	}
}
