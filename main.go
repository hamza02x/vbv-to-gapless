package main

import (
	"path"

	hel "github.com/hamza02x/go-helper"
)

var (
	vbvAudioDir   string                 // flag
	outDir        string                 // flag
	outBuildDir   string                 // outDir + "/build"
	outSuraDir    string                 // outDir + "/sura"
	thread        int                    // flag
	vbvAyaLengths = [TOTAL_AYA]float64{} // key/index: ayaId-1
)

func main() {

	handleFlags()

	for ayaId := 1; ayaId <= TOTAL_AYA; ayaId++ {

		suraAya := suraAyaFromAyaId(ayaId)
		outSuraFile := path.Join(outDir, "audios/"+getPartName(suraAya.Sura)+".mp3")
		hel.Pl(outSuraFile)
	}
}
