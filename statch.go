package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	_ "strconv"
	"time"
)

// logfile is a variable that holds the log file path.
// structure of the file: label start end duration
var logfile string = "./.statch.csv"

// statch is a structure that holds each logged item's label, start, end and duration.
type statch struct {
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
//IDEA: // If any previous insert was still open for given label, then that insert get's a closed value.
func start(label string) {
	srt := time.Now()
	line := fmt.Sprintf("%v,%v,%v,%v\n", label, srt.Unix(), "", "")
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

// logged is a function that CURRENTLY prints the contents of the logfile.
func logged(label string, printing bool) {
	if _, err := os.Stat(logfile); err != nil {
		fmt.Println("Nothing logged yet")
		return
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
	for scanner.Scan() {
		if printing {
			fmt.Println(scanner.Text())
		}
	}
}

func main() {
	labelFlag := flag.String("label", "general", "Label of saved time")
	flag.Parse()
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	switch action := os.Args[flag.NFlag()+1]; action {
	case "help":
		help()
	case "start":
		start(*labelFlag)
	case "logged":
		logged(*labelFlag, true)
	default:
		log.Fatal("Not defined")
	}
}
