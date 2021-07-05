package main

import (
	"os"
	"path"
	"strconv"
	"sync"

	col "github.com/hamza02x/go-color"
	hel "github.com/hamza02x/go-helper"
)

var (
	vbvAudioDir  string // flag
	outDir       string // flag
	name         string // flag
	outBuildDir  string // outDir + "/build"
	outSuraDir   string // outDir + "/sura"
	thread       int    // flag
	createdCount int
)

func main() {

	handleFlags()
	setDB(path.Join(outDir, name+".db"))

	var wg sync.WaitGroup
	var c = make(chan int, thread)

	for sura := 1; sura <= TOTAL_SURA; sura++ {
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

	os.RemoveAll(outBuildDir)
}

func concatSuraAudio(sura int) {
	outSuraFile := getSuraFilePath(sura)

	hel.Pl("🔪 Creating: " + col.Red(outSuraFile))
	execute("ffmpeg", "-f concat -safe 0 -i "+getFfmpegConcatFile(sura)+" "+outSuraFile+" -v quiet -y")
	hel.Pl("✅ " + strconv.Itoa(createdCount+1) + ". Created: " + col.Green(outSuraFile))

	createdCount++
}

func insertTimingRows(sura int) {

	var startTime int64 = 0

	if sura != SURA_FATIHA && sura != SURA_TAWBA {
		// bismillah
		startTime = getAudioLengthMS(getAyaFilePath(SURA_FATIHA, 1))
	}

	for aya := 1; aya <= AYAH_COUNT[sura-1]; aya++ {
		dbCreateOrSave(sura, aya, startTime, false)
		startTime += getAudioLengthMS(getAyaFilePath(sura, aya))
	}

	dbCreateOrSave(sura, 999, getAudioLengthMS(getSuraFilePath(sura)), false)
}
