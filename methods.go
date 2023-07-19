package main

import (
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	hel "github.com/hamza72x/go-helper"
)

var (
	regexSlugify = regexp.MustCompile("[^a-z0-9]+")
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
func getAudioLengthMS(path string) int64 {

	dur, err := strconv.ParseFloat(trimSpaces(execute(
		"ffprobe", "-i "+path+" -show_entries format=duration -v quiet -of csv=p=0",
	)), 64)

	panics("error in ParseFloat ", err)

	return int64(dur * 1000)
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

func getVbvAyaFileName(sura int, aya int) string {
	return getPartName(sura) + getPartName(aya) + ".mp3"
}

func getVbvAyaFilePath(sura int, aya int) string {
	return path.Join(getSuraDir(sura), getVbvAyaFileName(sura, aya))
}

func getGaplessMp3SuraFilePath(sura int) string {
	return path.Join(dirOutSura, getPartName(sura)+".mp3")
}

func getGaplessOpusSuraFilePath(sura int) string {
	return path.Join(dirOutSura, getPartName(sura)+".opus")
}

func slugify(s string) string {
	return strings.Trim(regexSlugify.ReplaceAllString(strings.ToLower(s), "-"), "-")
}

func getFfmpegConcatFilePath(sura int) string {
	return dirOutBuild + "/" + getPartName(sura) + ".txt"
}

func getAbs(path string) string {
	abs, err := filepath.Abs(path)
	panics("Error in getting absolute path of "+path, err)
	return abs
}

func getSuraDir(sura int) string {
	if isVbvAyaFileInSuraDir {
		return path.Join(dirVbvAudio, strconv.Itoa(sura))
	}
	return dirVbvAudio
}

func isDirExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}
