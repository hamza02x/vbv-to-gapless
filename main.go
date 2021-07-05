package main

import (
	"os"
	"path"
	"sync"

	col "github.com/hamza02x/go-color"
	hel "github.com/hamza02x/go-helper"
)

var (
	vbvAudioDir string // flag
	outDir      string // flag
	name        string // flag
	outBuildDir string // outDir + "/build"
	outSuraDir  string // outDir + "/sura"
	thread      int    // flag
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

	moveTimingUnorderedToMain()

	os.RemoveAll(outBuildDir)
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
		// bismillah
		endTime = getAudioLengthMS(getAyaFilePath(SURA_FATIHA, 1))
	}

	for aya := 1; aya <= AYAH_COUNT[sura-1]; aya++ {

		endTime += getAudioLengthMS(getAyaFilePath(sura, aya))

		db.Create(&TimingUnordered{Sura: sura, Ayah: aya, Time: endTime})
	}

	lengthFullSura := getAudioLengthMS(getSuraFilePath(sura))

	if endTime > lengthFullSura {
		db.Save(&TimingUnordered{Sura: sura, Ayah: AYAH_COUNT[sura-1], Time: lengthFullSura})
	}

	db.Create(&TimingUnordered{Sura: sura, Ayah: 999, Time: lengthFullSura})
}

func moveTimingUnorderedToMain() {

	for sura := 1; sura <= TOTAL_SURA; sura++ {
		for aya := 1; aya <= AYAH_COUNT[sura-1]; aya++ {
			var t TimingUnordered
			db.Where("sura = ? and ayah = ?", sura, aya).First(&t)
			db.Create(&Timing{Sura: t.Sura, Ayah: t.Ayah, Time: t.Time})
		}
		var t TimingUnordered
		db.Where("sura = ? and ayah = ?", sura, 999).First(&t)
		db.Create(&Timing{Sura: t.Sura, Ayah: t.Ayah, Time: t.Time})
	}

	db.Exec("drop table " + TimingUnordered{}.TableName())
	db.Exec("VACUUM")
}
