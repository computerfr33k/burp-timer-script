package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/ryanuber/go-glob"
)

func main() {
	argsWithoutProg := os.Args[1:]

	if len(argsWithoutProg) < 6 {
		fmt.Fprintln(os.Stderr, "Not enough arguments given.")
		os.Exit(1)
	}

	client := argsWithoutProg[0]
	current := argsWithoutProg[1]
	storage_dir := argsWithoutProg[2]
	//reserved1 := argsWithoutProg[3]
	//reserved2 := argsWithoutProg[4]
	interval := argsWithoutProg[5]
	timestamp := current + "/timestamp"

	// A 'backup' file placed in the storage directory tells this script that
	// a backup needs to be done right now.
	// This gives the 'server initiates a manual backup' feature.
	if force_manual_backup(storage_dir, client) {
		os.Exit(0)
	}

	// The rest of the arguments, if any, should be timebands.
	curDayHour := time.Now().Format("*Mon*15*")

	inTimeband := false // If no timebands given, default to not OK.
	timebands := argsWithoutProg[6:]
	if len(timebands) <= 0 {
		os.Exit(1)
	}

	for i := 0; i < len(timebands); i++ {
		inTimeband = false

		if glob.Glob(curDayHour, timebands[i]) {
			fmt.Println("In timeband: " + timebands[i])
			inTimeband = true
			break
		}

		fmt.Println("Out of timeband: " + timebands[i])
	}

	if inTimeband == false {
		os.Exit(1)
	}

	if get_intervals(current, client, timestamp, interval) {
		fmt.Println("Do a backup of " + client + " now.")
		os.Exit(0)
	}

	fmt.Println("Not yet time for a backup of " + client)
	os.Exit(1)
}

func force_manual_backup(storage_dir string, client string) bool {
	manual_file := storage_dir + "/" + client + "/backup"
	if _, err := os.Stat(manual_file); err == nil {
		fmt.Println("Found " + manual_file)
		fmt.Println("Do a backup of " + client + " now")

		err := os.Remove(manual_file)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return false
		}

		return true
	}

	return false
}

func get_intervals(current string, client string, timestamp string, interval string) bool {
	if _, err := os.Stat(current); os.IsNotExist(err) {
		fmt.Println("No Prior backup of " + client)

		return true
	}

	if _, err := os.Stat(timestamp); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Timestamp file missing for %s\n", client)

		return true
	}

	if len(interval) == 0 {
		fmt.Fprintf(os.Stderr, "No time interval given for %s\n", client)

		return false
	}

	intervalSecs := 0
	re := regexp.MustCompile("([0-9]+)")

	if res, _ := regexp.MatchString("[0-9]+s", interval); res {
		intervalSecs, _ = strconv.Atoi(re.FindString(interval))

	} else if res, _ := regexp.MatchString("[0-9]+.*m", interval); res {
		intervalSecs, _ = strconv.Atoi(re.FindString(interval))
		intervalSecs *= 60

	} else if res, _ := regexp.MatchString("[0-9]+h", interval); res {
		intervalSecs, _ = strconv.Atoi(re.FindString(interval))
		intervalSecs *= 60 * 60

	} else if res, _ := regexp.MatchString("[0-9]+.*d", interval); res {
		intervalSecs, _ = strconv.Atoi(re.FindString(interval))
		intervalSecs *= 60 * 60 * 24

	} else if res, _ := regexp.MatchString("[0-9]+.*w", interval); res {
		intervalSecs, _ = strconv.Atoi(re.FindString(interval))
		intervalSecs *= 60 * 60 * 24 * 7

	} else if res, _ := regexp.MatchString("[0-9]+.*n", interval); res {
		intervalSecs, _ = strconv.Atoi(re.FindString(interval))
		intervalSecs *= 60 * 60 * 24 * 30

	} else {
		fmt.Fprintf(os.Stderr, "interval %s not understood for %s\n", interval, client)

		return false
	}

	if intervalSecs == 0 {
		fmt.Fprintf(os.Stderr, "interval %s not understood for %s\n", interval, client)

		return false
	}

	lines, _ := readLines(timestamp)
	ts := lines[0]
	re = regexp.MustCompile("(\\d{4}-\\d{2}-\\d{2}\\s\\d{2}:\\d{2}:\\d{2})")
	ts = re.FindString(ts)

	const timeLayout = "2006-01-02 15:04:05"
	secs, _ := time.ParseInLocation(timeLayout, ts, time.Local) // YYYY-MM-DD hh:mm:ss
	now := time.Now()

	min_timesecs := secs.Unix() + int64(intervalSecs)
	min_time := time.Unix(min_timesecs, 0).Format(timeLayout)

	fmt.Printf("Last Backup: %s\n", ts)
	fmt.Printf("Next after: %s (interval %s)\n", min_time, interval)

	if min_timesecs < now.Unix() {
		fmt.Printf("%s < %s.\n", min_time, now.Format(timeLayout))

		return true
	}

	return false
}

// readLines reads a whole file into memory
// and returns a slice of its lines.
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
