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
	flag.StringVar(&outDir, "o", "", "output directory path (required)")
	flag.IntVar(&thread, "t", 10, "number of threads")

	flag.Parse()

	hel.Pl("verse by verse audio directory: ", col.Yellow(vbvAudioDir))
	hel.Pl("output directory: ", col.Red(outDir))

	if vbvAudioDir == "" || outDir == "" {
		flagExit()
	}

	validateVbvAudioDir()
	panics("Error creating outDir", hel.DirCreateIfNotExists(outDir))
	panics("Error creating outBuildDir", hel.DirCreateIfNotExists(outDir+"/build"))
	panics("Error creating outSuraDir", hel.DirCreateIfNotExists(outDir+"/sura"))

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

		for aya := 1; aya <= AYAH_COUNT[sura-1]; aya++ {

			ffmpegConcatData += "file \"" + getAyaFilePath(sura, aya) + "\"\n"

			wg.Add(1)
			go func(sura int, aya int) {
				c <- i
				// validateAyaFile(sura, aya)
				if i%100 == 0 {
					hel.Pl("Checking input files: " + strconv.Itoa(i))
				}
				<-c
				wg.Done()
				i++
			}(sura, aya)
		}

		hel.StrToFile(outBuildDir+"/"+getPartName(sura)+".txt", ffmpegConcatData)
	}

	wg.Wait()

	hel.Pl("All input audio files seems valid, creating gapless")
}

func validateAyaFile(sura int, aya int) {
	suraAya := SuraAya{sura, aya}
	fileName := getFileName(sura, aya)
	ayaFilePath := path.Join(vbvAudioDir, fileName)
	if !hel.FileExists(ayaFilePath) {
		panic("Audio file `" + fileName + "` doesn't exist")
	}
	vbvAyaLengths[suraAya.getAyaId()-1] = getAudioLength(ayaFilePath)
}

func getAyaFilePath(sura int, aya int) string {
	fileName := getFileName(sura, aya)
	return path.Join(vbvAudioDir, fileName)
}
