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

	manual_file := storage_dir + "/" + client + "/backup"
	if _, err := os.Stat(manual_file); err == nil {
		fmt.Println("Found " + manual_file)
		fmt.Println("Do a backup of " + client + " now")

		err := os.Remove(manual_file)
		if err != nil {
			fmt.Println(err)
			return
		}

		os.Exit(0)
	}

	// The rest of the arguments, if any, should be timebands.
	curDayHour := time.Now().Format("*Mon*03*")

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

	if get_intervals(current, client, timestamp, interval) == true {
		fmt.Println("Do a backup of " + client + " now.")
		os.Exit(0)
	}

	fmt.Println("Not yet time for a backup of " + client)
	os.Exit(1)
}

func get_intervals(current string, client string, timestamp string, interval string) bool {
	if _, err := os.Stat(current); os.IsNotExist(err) {
		fmt.Println("No Prior backup of " + client)

		return false
	}

	//if _,err := os.Stat(timestamp); os.IsNotExist(err) {
	//	fmt.Println("Timestamp file missing for " + client)
	//	return false
	//}

	if len(interval) == 0 {
		fmt.Println("No time interval given for " + client)

		return false
	}

	intervalSecs := 0
	re := regexp.MustCompile("([0-9]+)")

	if _, err := regexp.MatchString("[0-9]+.*s", interval); err == nil {
		intervalSecs, _ = strconv.Atoi(re.FindString(interval))

	} else if _, err := regexp.MatchString("[0-9]+.*m", interval); err == nil {
		intervalSecs, _ = strconv.Atoi(re.FindString(interval))
		intervalSecs *= 60

	} else if _, err := regexp.MatchString("[0-9]+.*h", interval); err == nil {
		intervalSecs, _ = strconv.Atoi(re.FindString(interval))
		intervalSecs *= 60 * 60

	} else if _, err := regexp.MatchString("[0-9]+.*d", interval); err == nil {
		intervalSecs, _ = strconv.Atoi(re.FindString(interval))
		intervalSecs *= 60 * 60 * 24

	} else if _, err := regexp.MatchString("[0-9]+.*w", interval); err == nil {
		intervalSecs, _ = strconv.Atoi(re.FindString(interval))
		intervalSecs *= 60 * 60 * 24 * 7

	} else if _, err := regexp.MatchString("[0-9]+.*n", interval); err == nil {
		intervalSecs, _ = strconv.Atoi(re.FindString(interval))
		intervalSecs *= 60 * 60 * 24 * 30

	} else {
		fmt.Println("interval " + interval + " not understood for " + client)

		return false
	}

	if intervalSecs == 0 {
		fmt.Println("interval " + interval + " not understood for " + client)

		return false
	}

	//lines, err := readLines(timestamp)
	//if err != nil {
	//	log.Fatalf("readLines: %s", err)
	//}
	//ts := lines[0]

	return true
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
