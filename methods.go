package main

import (
	"os/exec"
	"strconv"
	"strings"

	hel "github.com/hamza02x/go-helper"
)

// 1   => 001
// 12  => 012
// 122 => 122
func getPartName(suraOrAya int) string {

	var part = ""

	if suraOrAya < 10 {
		part = "00" + strconv.Itoa(suraOrAya)
	} else if suraOrAya < 100 {
		part = "0" + strconv.Itoa(suraOrAya)
	} else {
		part = strconv.Itoa(suraOrAya)
	}

	return part
}

// ffprobe -i 001001.mp3 -show_entries format=duration -v quiet -of csv="p=0"
func getAudioLength(path string) float64 {

	dur, err := strconv.ParseFloat(trimSpaces(execute(
		"ffprobe", "-i "+path+" -show_entries format=duration -v quiet -of csv=p=0",
	)), 64)

	panics("error in ParseFloat ", err)

	return dur
}

func execute(comm string, arg string) string {
	cmd := exec.Command(comm, hel.StrToArr(arg, " ")...)

	out, err := cmd.CombinedOutput()

	if err != nil {
		panics("Error in execute "+comm+" "+arg, err)
	}

	return string(out)
}

func panics(title string, err error) {
	if err != nil {
		hel.Pl(err)
		panic("`" + title + "`")
	}
}

func trimSpaces(str string) string {
	return strings.ReplaceAll(strings.TrimSpace(str), " ", "")
}

func getFileName(sura int, aya int) string {
	return getPartName(sura) + getPartName(aya) + ".mp3"
}

func inRange(val, min, max int) bool {
	return val >= min && val <= max
}
