package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	col "github.com/hamza72x/go-color"
	hel "github.com/hamza72x/go-helper"
)

func handleFlags() error {

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

	var err error
	dirOut, err = getAbs(dirOut)
	if err != nil {
		return err
	}

	dirVbvAudio, err = getAbs(dirVbvAudio)
	if err != nil {
		return err
	}

	dirOutBuild, err = getAbs(dirOut + "/build")
	if err != nil {
		return err
	}

	dirOutSura, err = getAbs(dirOut + "/" + name)
	if err != nil {
		return err
	}

	if err := os.RemoveAll(dirOutBuild); err != nil {
		return err
	}
	if err := hel.DirCreateIfNotExists(dirOut); err != nil {
		return err
	}
	if err := hel.DirCreateIfNotExists(dirOutBuild); err != nil {
		return err
	}
	if err := hel.DirCreateIfNotExists(dirOutSura); err != nil {
		return err
	}

	hel.Pl("verse by verse audio directory: ", col.Yellow(dirVbvAudio))
	hel.Pl("output directory: ", col.Red(dirOut))

	time.Sleep(1 * time.Second)

	return err
}

func flagExit() {
	flag.Usage()
	os.Exit(1)
}

// getSuras get all suras and validate them
func getSuras() ([]int, error) {

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
			log.Println(col.Red(fmt.Sprintf("sura directory `%s` doesn't exist; skipping\n", dirSura)))
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
			log.Println(col.Yellow(fmt.Sprintf("sura `%d` is incomplete; skipping\n", sura)))
			continue
		}

		suras = append(suras, sura)
		if err := hel.StrToFile(getFfmpegConcatFilePath(sura), ffmpegConcatData); err != nil {
			return []int{}, err
		}
	}

	wg.Wait()
	close(c)

	hel.Pl("audio files checked, valid sura(s) =>", suras)

	return suras, nil
}

func validateAyaFile(sura int, aya int) {
	if !hel.FileExists(getVbvAyaFilePath(sura, aya)) {
		panic("Audio file `" + getVbvAyaFileName(sura, aya) + "` doesn't exist")
	}
}
