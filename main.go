package main

import (
	"fmt"
	"log"
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

	err := handleFlags()
	if err != nil {
		log.Fatal(err)
	}

	suras, err := getSuras()
	if err != nil {
		log.Fatal(err)
	}

	if err := setDB(path.Join(dirOut, name+".db")); err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	var c = make(chan int, thread)

	for _, sura := range suras {
		wg.Add(1)
		go func(sura int) {
			defer wg.Done()
			c <- sura
			if err := concatSuraAudio(sura); err != nil {
				log.Printf("Error in concat sura audio, sura: %d, error: %v", sura, err)
			} else if err := insertTimingRows(sura); err != nil {
				log.Printf("Error in inserting timing rows, sura: %d, error: %v", sura, err)
			}
			<-c
		}(sura)
	}

	wg.Wait()
	close(c)

	if err := dbVaccum(); err != nil {
		log.Printf("Error in vacuuming database, error: %v", err)
	}

	if err := os.RemoveAll(dirOutBuild); err != nil {
		log.Printf("Error in removing build directory, error: %v", err)
	}
}

func concatSuraAudio(sura int) error {
	outMp3File := getGaplessMp3SuraFilePath(sura)
	concatFile := getFfmpegConcatFilePath(sura)

	hel.Pl("🔪 Creating: " + col.Red(outMp3File))
	if _, err := execute("ffmpeg", fmt.Sprintf(
		"-f concat -safe 0 -i %s %s -v quiet -y",
		concatFile, outMp3File,
	)); err != nil {
		return err
	}

	hel.Pl("✅ " + strconv.Itoa(createdCount+1) + ". Created: " + col.Green(outMp3File))

	// also create opus version
	outOpusFile := getGaplessOpusSuraFilePath(sura)

	hel.Pl("🔪 Creating: " + col.Red(outOpusFile))
	if _, err := execute("ffmpeg", fmt.Sprintf(
		"-i %s -c:a libopus -vbr on -compression_level 10 -frame_duration 60 -application audio -v quiet -y %s",
		outMp3File,
		outOpusFile,
	)); err != nil {
		return err
	}
	hel.Pl("✅ " + strconv.Itoa(createdCount+1) + ". Created: " + col.Green(outOpusFile))

	createdCount++

	return nil
}

func insertTimingRows(sura int) error {

	var startTime int64 = 0
	var err error

	for aya := 1; aya <= AYAH_COUNT[sura-1]; aya++ {
		if err := dbUpdateTiming(sura, aya, startTime); err != nil {
			return err
		}
		time, err := getAudioLengthMS(getVbvAyaFilePath(sura, aya))
		if err != nil {
			return err
		}
		startTime += time
	}

	audioLength, err := getAudioLengthMS(getGaplessMp3SuraFilePath(sura))
	if err != nil {
		return err
	}

	return dbUpdateTiming(sura, 999, audioLength)
}
