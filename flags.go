package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
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
	flag.BoolVar(&isOpusToo, "opus", false, "also create opus files? (default false)")

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
	missingSuras := []int{}

	// key: sura, value: missing ayas
	missingSuraAya := map[int][]int{}

	if !hel.PathExists(dirVbvAudio) {
		panic("directory `" + dirVbvAudio + "` doesn't exist")
	}

	totalMissingSuraAya := 0

	for sura := 1; sura <= TOTAL_SURA; sura++ {

		ffmpegConcatData := ""

		dirSura := getSuraDir(sura)

		// skip a sura if it's directory doesn't exist
		if !isDirExists(dirSura) {
			missingSuras = append(missingSuras, sura)
			continue
		}

		// skip a sura if it's incomplete
		isSuraIncomplete := false

		for aya := 1; aya <= AYAH_COUNT[sura-1]; aya++ {
			if isSuraIncomplete {
				totalMissingSuraAya++
				missingSuraAya[sura] = append(missingSuraAya[sura], aya)
				continue
			}

			vbvAyaPath := getVbvAyaFilePath(sura, aya)

			if !hel.FileExists(vbvAyaPath) {
				totalMissingSuraAya++
				isSuraIncomplete = true
				if _, ok := missingSuraAya[sura]; !ok {
					missingSuraAya[sura] = []int{}
				}
				missingSuraAya[sura] = append(missingSuraAya[sura], aya)
				continue
			}

			ffmpegConcatData += fmt.Sprintf("file '%s'\n", vbvAyaPath)
		}

		// skip a sura if it's incomplete
		if !isSuraIncomplete {
			suras = append(suras, sura)
			if err := hel.StrToFile(getFfmpegConcatFilePath(sura), ffmpegConcatData); err != nil {
				return []int{}, err
			}
		}
	}

	hel.Pl("missing sura ayas(s) =>")

	// sort by sura
	suraList := []int{}
	for sura := range missingSuraAya {
		suraList = append(suraList, sura)
	}
	sort.Ints(suraList)
	for sura := range missingSuraAya {
		for _, aya := range missingSuraAya[sura] {
			fmt.Println(getVbvAyaFileName(sura, aya))
		}
	}
	hel.Pl("total missing sura ayas: ", col.Red(totalMissingSuraAya))

	hel.Pl("missing full sura(s) =>")
	for _, sura := range missingSuras {
		fmt.Println(sura)
	}

	return suras, nil
}
