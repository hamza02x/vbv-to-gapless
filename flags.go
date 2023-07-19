package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	col "github.com/hamza72x/go-color"
	hel "github.com/hamza72x/go-helper"
)

func handleFlags() {

	flag.StringVar(&dirVbvAudio, "vd", "", "verse by verse audio directory (required)")
	flag.StringVar(&name, "n", "", "database name (required) (ex: husary)")
	flag.StringVar(&dirOut, "o", "", "output directory path (required)")
	flag.IntVar(&thread, "t", 10, "number of threads")
	flag.BoolVar(&isVbvAyaFileInSuraDir, "visd", false, "is vbv file in their sura directory? (default false); ex: 2/002001.mp3")

	flag.Parse()

	if dirVbvAudio == "" || dirOut == "" || name == "" {
		flagExit()
	}

	name = slugify(name)

	dirOut = getAbs(dirOut)
	dirVbvAudio = getAbs(dirVbvAudio)

	dirOutBuild = getAbs(dirOut + "/build")
	dirOutSura = getAbs(dirOut + "/" + name)

	os.RemoveAll(dirOutBuild)

	panics("Error creating dirOut", hel.DirCreateIfNotExists(dirOut))
	panics("Error creating dirOutBuild", hel.DirCreateIfNotExists(dirOutBuild))
	panics("Error creating dirOutSura", hel.DirCreateIfNotExists(dirOutSura))

	hel.Pl("verse by verse audio directory: ", col.Yellow(dirVbvAudio))
	hel.Pl("output directory: ", col.Red(dirOut))

	time.Sleep(1 * time.Second)

}

func flagExit() {
	flag.Usage()
	os.Exit(1)
}

// getSuras get all suras and validate them
func getSuras() []int {

	suras := []int{}

	if !hel.PathExists(dirVbvAudio) {
		panic("directory `" + dirVbvAudio + "` doesn't exist")
	}

	var wg sync.WaitGroup
	var c = make(chan int, thread)
	var i = 0

	for sura := 1; sura <= TOTAL_SURA; sura++ {

		ffmpegConcatData := ""

		dirSura := getSuraDir(sura)

		// skip a sura if it's directory doesn't exist
		if !isDirExists(dirSura) {
			fmt.Printf("sura directory `%s` doesn't exist; skipping\n", dirSura)
			continue
		}

		// skip a sura if it's incomplete
		isSuraIncomplete := false

		for aya := 1; aya <= AYAH_COUNT[sura-1]; aya++ {

			vbvAyaPath := getVbvAyaFilePath(sura, aya)
			if !hel.FileExists(vbvAyaPath) {
				isSuraIncomplete = true
				break
			}

			ffmpegConcatData += fmt.Sprintf("file '%s'\n", vbvAyaPath)

			wg.Add(1)
			go func(sura int, aya int) {
				c <- i
				validateAyaFile(sura, aya)
				if i%500 == 0 {
					hel.Pl("Checking input files: " + strconv.Itoa(i))
				}
				<-c
				wg.Done()
				i++
			}(sura, aya)
		}

		// skip a sura if it's incomplete
		if isSuraIncomplete {
			fmt.Printf("sura `%d` is incomplete; skipping\n", sura)
			continue
		}

		suras = append(suras, sura)
		panics("Error creating contact data file", hel.StrToFile(getFfmpegConcatFilePath(sura), ffmpegConcatData))
	}

	wg.Wait()
	close(c)

	hel.Pl("All input audio files seems valid, creating gapless")

	return suras
}

func validateAyaFile(sura int, aya int) {
	if !hel.FileExists(getVbvAyaFilePath(sura, aya)) {
		panic("Audio file `" + getVbvAyaFileName(sura, aya) + "` doesn't exist")
	}
}
