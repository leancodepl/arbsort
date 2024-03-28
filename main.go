package main

import (
	"encoding/json"
	"os"
	"sort"
	"strings"

	orderedmap "github.com/wk8/go-ordered-map/v2"
)

type jsonLine struct {
	Key      string
	Val      any
	Metadata any
}

func main() {
	if len(os.Args) != 2 {
		panic("usage: arbsort <arb filename>")
	}

	filename := os.Args[1]
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var jsonFile map[string]any
	err = json.NewDecoder(file).Decode(&jsonFile)
	if err != nil {
		panic(err)
	}

	jsonLines := make([]jsonLine, 0, len(jsonFile))
	for key, val := range jsonFile {
		if strings.HasPrefix(key, "@") && !strings.HasPrefix(key, "@@") {
			continue
		}

		jsonLines = append(jsonLines, jsonLine{key, val, jsonFile["@"+key]})
	}

	sort.Slice(jsonLines, func(i, j int) bool {
		iKey, jKey := jsonLines[i].Key, jsonLines[j].Key

		if strings.HasPrefix(iKey, "@@") {
			if strings.HasPrefix(jKey, "@@") {
				return iKey < jKey
			}

			return true
		} else if strings.HasPrefix(jKey, "@@") {
			return false
		}

		return iKey < jKey
	})

	newJsonContents := orderedmap.New[string, any]()
	for _, line := range jsonLines {
		newJsonContents.Set(line.Key, line.Val)
		if line.Metadata != nil {
			newJsonContents.Set("@"+line.Key, line.Metadata)
		}
	}

	file, err = os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	err = encoder.Encode(newJsonContents)
	if err != nil {
		panic(err)
	}
}
