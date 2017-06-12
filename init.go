package quranize

import (
	"encoding/xml"
	"io/ioutil"
	"strings"
)

func init() {
	loadHijaiyas()
	loadQuran()
	buildIndex()
}

func loadHijaiyas() {
	filePath := "corpus/arabic-to-alphabet"
	raw, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	trimmed := strings.TrimSpace(string(raw))
	for _, line := range strings.Split(trimmed, "\n") {
		components := strings.Split(line, " ")
		arabic := components[0]
		for _, alphabet := range components[1:] {
			hijaiyas[alphabet] = append(hijaiyas[alphabet], arabic)
			if maxWidth < len(alphabet) {
				maxWidth = len(alphabet)
			}
		}
	}
}

func loadQuran() {
	filePath := "corpus/quran-simple-clean.xml"
	raw, err := ioutil.ReadFile(filePath)
	if err == nil {
		err = xml.Unmarshal(raw, &Quran)
	}
	if err != nil {
		panic(err)
	}
}

func buildIndex() {
	for s, sura := range Quran.Suras {
		for a, aya := range sura.Ayas {
			indexAya(aya.Text, Location{Sura: s, Aya: a})
		}
	}
}

func indexAya(text string, location Location) {
	start := 0
	for {
		location.Begin = start
		text = text[start:]
		buildTree([]rune(text), location, root)
		start = strings.Index(text, " ") + 1
		if start == 0 {
			break
		}
	}
}

func buildTree(harfs []rune, location Location, node *Node) {
	for i, harf := range harfs {
		location.End = location.Begin + i + 1
		if node.Children[harf] == nil {
			node.Children[harf] = &Node{Children: make(Children)}
		}
		node = node.Children[harf]
		node.Locations = insertLocation(node.Locations, location)
	}
}

func insertLocation(locations []Location, newLocation Location) []Location {
	for _, location := range locations {
		if newLocation == location {
			return locations
		}
	}
	return append(locations, newLocation)
}
