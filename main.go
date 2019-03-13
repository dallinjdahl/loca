package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	// "fmt"
	"io"
	"log"
	"os"
	"strings"
	"unicode/utf8"
	"path/filepath"
)

func check(e error, s string) {
	if e != nil {
		log.Fatal(e, s)
	}
}

func parseStrings(r *csv.Reader, lang string) map[string]string {
	head, err := r.Read()
	check(err, "parser head")

	index := 0
	dict := make(map[string]string)

	for i := range head {
		if head[i] == lang {
			index = i + 1
			break
		}
	}
	if index == 0 {
		log.Fatal("Language not defined: ", lang)
	}

	r.FieldsPerRecord = 0

	body, err := r.ReadAll()
	check(err, "parser body")

	for _, f := range body {
		dict[f[0]] = f[index]
	}
	return dict
}

func processString(r *bufio.Reader, spe rune) string {
	s := ""
	for {
		c, _, err := r.ReadRune()
		if err != io.EOF {
			check(err, "process string")
		}

		if strings.ContainsRune(
			" \n\t~!@#$%^&*()`=,<.>/?;:'\"\\|[{]}"+ string(spe), c) || err == io.EOF {
			// to permit most programming languages
			// and still allow for international names

			r.UnreadRune()
			return s
		} else {
			s = s + string(c)
		}
	}
	return "" // unreachable
}

func main() {
	charPtr := flag.String("c", "`", "string interpolation character")
	langPtr := flag.String("l", "eng", "language to interpolate")
	outPtr := flag.String("o", "/", "output filename")

	flag.Parse()

	srcName := flag.Arg(0)
	stringsName := flag.Arg(1)
	if stringsName == "" {
		stringsName = "locali"
	}

	srcFile, err := os.Open(srcName)
	check(err, "open src")
	defer srcFile.Close()

	stringsFile, err := os.Open(stringsName)
	check(err, "open strings")
	defer stringsFile.Close()

	outFile := os.Stdout
	if *outPtr == "/" {
		outName := strings.TrimSuffix(srcName,
								filepath.Ext(srcName))

		outFile, err = os.Create(*langPtr + "." + outName)
		check(err, "open default out")
		defer outFile.Close()
	} else if *outPtr != "-" {
		outFile, err = os.Create(*outPtr)
		check(err, "open output")
		defer outFile.Close()
	}

	src := bufio.NewReader(srcFile)
	strings := csv.NewReader(stringsFile)
	strings.FieldsPerRecord = -1
	out := bufio.NewWriter(outFile)

	char, _ := utf8.DecodeRuneInString(*charPtr)

	dict := parseStrings(strings, *langPtr)
	for {
		c, _, err := src.ReadRune()
		if err == io.EOF {
			break
		} else {
			check(err, "Read main body")
		}
		if c == char {
			s := processString(src, char)
			_, err = out.WriteString(dict[s])
			check(err, "writing new strings")
		} else {
			out.WriteRune(c)
		}
	}
	out.Flush()
}
