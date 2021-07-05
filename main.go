package main

import (
	"path"
	"sync"

	col "github.com/hamza02x/go-color"
	hel "github.com/hamza02x/go-helper"
)

var (
	vbvAudioDir     string               // flag
	outDir          string               // flag
	name            string               // flag
	outBuildDir     string               // outDir + "/build"
	outSuraDir      string               // outDir + "/sura"
	thread          int                  // flag
	vbvAyaLengths   = [TOTAL_AYA]int64{} // key/index: ayaId-1
	lengthBismillah int64
)

func main() {

	handleFlags()
	setDB(path.Join(outDir, name))

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
}

func concatSuraAudio(sura int) {
	outSuraFile := getSuraFilePath(sura)

	hel.Pl("Creating: " + col.Red(outSuraFile))
	execute("ffmpeg", "-f concat -safe 0 -i "+getFfmpegConcatFile(sura)+" "+outSuraFile+" -v quiet -y")
	hel.Pl("Created: " + col.Green(outSuraFile))
}

func insertTimingRows(sura int) {

	var endTime int64 = 0

	if sura != SURA_FATIHA && sura != SURA_TAWBA {
		endTime = lengthBismillah
	}

	for aya := 1; aya <= AYAH_COUNT[sura-1]; aya++ {

		endTime += vbvAyaLengths[getAyaIndex(sura, aya)]

		db.Save(&Timing{Sura: sura, Ayah: aya, Time: endTime})
	}

	db.Save(&Timing{Sura: sura, Ayah: 999, Time: getAudioLengthMS(getSuraFilePath(sura))})
}
