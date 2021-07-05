package main

import (
	"flag"
	"os"
	"path"
	"strconv"
	"sync"

	col "github.com/hamza02x/go-color"
	hel "github.com/hamza02x/go-helper"
)

func handleFlags() {

	flag.StringVar(&vbvAudioDir, "vd", "", "verse by verse audio directory (required)")
	flag.StringVar(&name, "n", "", "database name (required)")
	flag.StringVar(&outDir, "o", "", "output directory path (required)")
	flag.IntVar(&thread, "t", 10, "number of threads")

	flag.Parse()

	hel.Pl("verse by verse audio directory: ", col.Yellow(vbvAudioDir))
	hel.Pl("output directory: ", col.Red(outDir))

	if vbvAudioDir == "" || outDir == "" || name == "" {
		flagExit()
	}

	outBuildDir = outDir + "/build"
	outSuraDir = outDir + "/sura"

	panics("Error creating outDir", hel.DirCreateIfNotExists(outDir))
	panics("Error creating outBuildDir", hel.DirCreateIfNotExists(outBuildDir))
	panics("Error creating outSuraDir", hel.DirCreateIfNotExists(outSuraDir))

	validateVbvAudioDir()
}

func flagExit() {
	flag.Usage()
	os.Exit(1)
}

func validateVbvAudioDir() {

	if !hel.PathExists(vbvAudioDir) {
		panic("directory `" + vbvAudioDir + "` doesn't exist")
	}

	var wg sync.WaitGroup
	var c = make(chan int, thread)
	var i = 0

	for sura := 1; sura <= TOTAL_SURA; sura++ {

		ffmpegConcatData := ""

		// contain bismillah in every other sura
		if sura != SURA_FATIHA && sura != SURA_TAWBA {
			ffmpegConcatData = "file '" + getAyaFilePath(SURA_FATIHA, 1) + "'\n"
		}

		for aya := 1; aya <= AYAH_COUNT[sura-1]; aya++ {

			ffmpegConcatData += "file '" + getAyaFilePath(sura, aya) + "'\n"

			wg.Add(1)
			go func(sura int, aya int) {
				c <- i
				validateAyaFile(sura, aya)
				if i%100 == 0 {
					hel.Pl("Checking input files: " + strconv.Itoa(i))
				}
				<-c
				wg.Done()
				i++
			}(sura, aya)
		}

		panics("Error creating contact data file", hel.StrToFile(getFfmpegConcatFile(sura), ffmpegConcatData))
	}

	wg.Wait()
	close(c)

	hel.Pl("All input audio files seems valid, creating gapless")
}

func validateAyaFile(sura int, aya int) {

	fileName := getSuraAyaFileName(sura, aya)
	ayaFilePath := getAyaFilePath(sura, aya)

	if !hel.FileExists(ayaFilePath) {
		panic("Audio file `" + fileName + "` doesn't exist")
	}

	timeMS := getAudioLengthMS(ayaFilePath)

	if sura == SURA_FATIHA && aya == 1 {
		lengthBismillah = timeMS
	}

	vbvAyaLengths[getAyaIndex(sura, aya)] = timeMS
}

func getAyaFilePath(sura int, aya int) string {
	fileName := getSuraAyaFileName(sura, aya)
	return path.Join(vbvAudioDir, fileName)
}

func getFfmpegConcatFile(sura int) string {
	return outBuildDir + "/" + getPartName(sura) + ".txt"
}
