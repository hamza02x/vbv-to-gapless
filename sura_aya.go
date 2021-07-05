package main

import hel "github.com/hamza02x/go-helper"

type SuraAya struct {
	Sura int
	Aya  int
}

func suraAyaFromAyaId(ayaId int) SuraAya {

	if !inRange(ayaId, 1, TOTAL_AYA) {
		panic("Invalid aya id")
	}

	var foundSura = 1
	var foundAya = 1

	if ayaId > 6230 {
		foundSura = 114
		foundAya = ayaId - 6230
	} else {
		for sura := 1; sura <= TOTAL_SURA; sura++ {
			var startAyaId = AYA_ID_START[sura-1]
			if startAyaId > ayaId {
				foundSura = sura - 1
				foundAya = ayaId - AYA_ID_START[foundSura-1] + 1
				break
			}
		}
	}

	s := SuraAya{Sura: foundSura, Aya: foundAya}

	if !s.isValidSura() {
		hel.Pl("Invalid suraAya, panicing", s)
		panic("Invalid suraAya")
	}

	return s
}

func (s *SuraAya) isValidSura() bool {
	suraZB := s.Sura - 1
	ayaZB := s.Aya - 1
	return suraZB >= 0 && suraZB <= 113 && ayaZB >= 0 && ayaZB <= AYAH_COUNT[suraZB]-1
}

func (s *SuraAya) getAyaId() int {
	return AYA_ID_START[s.Sura-1] + (s.Aya - 1)
}
