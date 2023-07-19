package main

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"sync"

	col "github.com/hamza72x/go-color"
	hel "github.com/hamza72x/go-helper"
)

var (
	name                  string // flag
	dirVbvAudio           string // flag
	dirOut                string // flag
	dirOutBuild           string // dirOut + "/build"
	dirOutSura            string // dirOut + "/$name"
	thread                int    // flag
	isVbvAyaFileInSuraDir bool   // flag
	createdCount          int
)

func main() {

	handleFlags()
	suras := getSuras()
	setDB(path.Join(dirOut, name+".db"))

	var wg sync.WaitGroup
	var c = make(chan int, thread)

	for _, sura := range suras {
		wg.Add(1)
		go func(sura int) {
			c <- sura
			concatSuraAudio(sura)
			insertTimingRows(sura)
			<-c
			wg.Done()
		}(sura)
	}

	wg.Wait()
	close(c)

	dbVaccum()

	os.RemoveAll(dirOutBuild)
}

func concatSuraAudio(sura int) {
	outMp3File := getGaplessMp3SuraFilePath(sura)
	contactFile := getFfmpegConcatFilePath(sura)

	hel.Pl("🔪 Creating: " + col.Red(outMp3File))
	execute("ffmpeg", fmt.Sprintf(
		"-f concat -safe 0 -i %s %s -v quiet -y",
		contactFile, outMp3File,
	))
	hel.Pl("✅ " + strconv.Itoa(createdCount+1) + ". Created: " + col.Green(outMp3File))

	// also create opus version
	outOpusFile := getGaplessOpusSuraFilePath(sura)

	hel.Pl("🔪 Creating: " + col.Red(outOpusFile))
	execute("ffmpeg", fmt.Sprintf(
		"-i %s -c:a libopus -b:a 16k -vbr on -compression_level 10 -frame_duration 60 -application audio -v quiet -y %s",
		outMp3File,
		outOpusFile,
	))
	hel.Pl("✅ " + strconv.Itoa(createdCount+1) + ". Created: " + col.Green(outOpusFile))

	createdCount++
}

func insertTimingRows(sura int) {

	var startTime int64 = 0

	if sura != SURA_FATIHA && sura != SURA_TAWBA {
		// bismillah
		startTime = getAudioLengthMS(getVbvAyaFilePath(SURA_FATIHA, 1))
	}

	for aya := 1; aya <= AYAH_COUNT[sura-1]; aya++ {
		dbUpdateTiming(sura, aya, startTime)
		startTime += getAudioLengthMS(getVbvAyaFilePath(sura, aya))
	}

	dbUpdateTiming(sura, 999, getAudioLengthMS(getGaplessMp3SuraFilePath(sura)))
}
