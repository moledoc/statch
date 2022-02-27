package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

// NOTE: only one open label at a time is allowed

// logfile is a variable that holds the log file path.
// structure of the file: label start end
var logfile string = ".trak.csv"

// label is a variable that holds the label value for trak
// Default trak label value is 'all'.
var label string = "all"

// compare is the default value of time.Time and is used to check if trak end is set or not.
var compare time.Time

// trak is a structure that holds each logged item's label, start, end and duration.
type trak struct {
	label    string
	start    time.Time
	end      time.Time
	duration time.Duration
}

// store is a method that formats trak into defined storing format.
func (t trak) store() string {
	var saveEnd string
	if t.end != compare {
		saveEnd = strconv.FormatInt(t.end.Unix(), 10)
	}
	return fmt.Sprintf("%v,%v,%v\n", t.label, t.start.Unix(), saveEnd)
}

// format is a variable that defines the trak printing format.
var format string = "%-10v %-30v %-30v %-10v"

// header is a variable that contains the header labels for printing traks.
var header string = fmt.Sprintf(format, "label", "start", "end", "duration")

// String is a method that prints trak to human readable format.
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

// TODO: help is a function that prints help.
func help() {
	fmt.Println("TODO:")
}

// save is a function that writes traks to the logfile.
func save(traks *[]trak) {
	f, err := os.OpenFile(logfile, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0600)
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}
	for _, elem := range *traks {
		_, err = f.WriteString(elem.store())
		if err != nil {
			log.Fatal(err)
		}
	}
}

// end is a function that closes the last opened insert for corresponding label.
func end(traks *[]trak, openLabel int) {
	if openLabel != -1 {
		cur := (*traks)[openLabel]
		cur.end = time.Unix(time.Now().Unix(), 0)
		cur.duration = cur.end.Sub(cur.start)
		(*traks)[openLabel] = cur
		fmt.Printf("Closed '%v'\n", cur.label)

	}
	save(traks)
	printer(traks, len(*traks)-5)
}

// start is a function that starts a new insert for given label.
// If any previous insert was still open for given label, then that insert gets closed.
func start(label string, traks *[]trak, openLabel int) {
	*traks = append(*traks, trak{label, time.Unix(time.Now().Unix(), 0), compare, time.Duration(0)})
	end(traks, openLabel)
	fmt.Printf("Started '%v'\n", label)
}

// printer is a function that prints traks in human readable format.
func printer(traks *[]trak, begInd int) {
	if begInd < 0 {
		begInd = 0
	}
	fmt.Println(header)
	for i := begInd; i < len(*traks); i++ {
		// 		for _, elem := range traks {
		if label == "all" || (*traks)[i].label == label {
			fmt.Println((*traks)[i].String())
		}
	}
}

// sumry is a structure to hold duration information about given label+timeframe
type sumry struct {
	label     string
	timeframe string
	duration  time.Duration
}

// mapKeys is a function that extracts and sorts all keys from given map.
func mapKeys(s map[string]sumry) []string {
	var re *regexp.Regexp = regexp.MustCompile(fmt.Sprintf("_%v$", label))
	var keys []string
	for k := range s {
		if label != "all" && !re.MatchString(k) {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// sumryPrinter is function that prints given sumry contents in certain format.
func sumryPrinter(s map[string]sumry, f string) {
	keys := mapKeys(s)
	fmt.Println("________________________________________________")
	fmt.Printf("%-10v %-10v %-10v\n", "label", f, "duration")
	fmt.Println("---------- ---------- ----------")
	for _, k := range keys {
		v := s[k]
		fmt.Printf("%-10v %-10v %-10v\n", v.label, v.timeframe, v.duration)
	}
}

// summary is a function that converts traks info into monthly, weekly and daily summary.
func summary(traks *[]trak, label string) {
	daily := make(map[string]sumry)
	monthly := make(map[string]sumry)
	weekly := make(map[string]sumry)
	for _, elem := range *traks {
		y, m, d := elem.start.Date()
		_, w := elem.start.ISOWeek()
		dayId := fmt.Sprintf("%v_%v", y*10000+int(m)*100+d, elem.label)
		monthId := fmt.Sprintf("%v_%v", y*100+int(m), elem.label)
		weekId := fmt.Sprintf("%v_%v", y*100+w, elem.label)

		if day, ok := daily[dayId]; !ok {
			daily[dayId] = sumry{elem.label, fmt.Sprintf("%v-%02d-%02d", y, int(m), d), elem.duration}
		} else {
			day.duration += elem.duration
			daily[dayId] = day
		}

		if month, ok := monthly[monthId]; !ok {
			monthly[monthId] = sumry{elem.label, fmt.Sprintf("%v-%02d", y, int(m)), elem.duration}
		} else {
			month.duration += elem.duration
			monthly[monthId] = month
		}

		if week, ok := weekly[weekId]; !ok {
			weekly[weekId] = sumry{elem.label, fmt.Sprintf("%v/%v", y, w), elem.duration}
		} else {
			week.duration += elem.duration
			weekly[weekId] = week
		}
	}
	sumryPrinter(monthly, "yyyy-mm")
	sumryPrinter(weekly, "yyyy/w")
	sumryPrinter(daily, "yyyy-mm-dd")
}

// trak [action] (label)
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
		printer(&traks, 0)
	case "start":
		start(label, &traks, openLabel)
	case "end":
		if openLabel == -1 {
			fmt.Println("No trak to close")
			return
		}
		end(&traks, openLabel)
	case "summary":
		summary(&traks, label)
	default:
		help()
	}
}
