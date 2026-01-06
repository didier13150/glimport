package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/didier13150/gllib"
)

func main() {

	var inputFile = flag.String("input", ".env", "File which contains vars to import.")
	var outputFile = flag.String("output", ".gitlab-vars.new.json", "JSON output file which can be imported by glcli.")
	var prefix = flag.String("prefix", "", "Prefix to apply to var name")
	var verbose = flag.Bool("verbose", false, "Make application more talkative.")
	var isCsv = flag.Bool("csv", false, "Input file is a CSV")

	if *verbose {
		log.Print("Verbose mode is active")
	}

	flag.Usage = func() {
		fmt.Print("Convert shell env file to glcli ready to import file\n\n")
		fmt.Printf("Usage: " + os.Args[0] + " [options]\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	vars := gllib.NewGitlabVar("", "", *verbose)

	varfile, err := os.OpenFile(*outputFile, os.O_RDONLY, 0644)
	if err == nil {
		vars.ImportVars(*outputFile)
		err = varfile.Close()
		if err != nil {
			log.Fatalln("Cannot close var file")
		}
	}

	content, err := os.ReadFile(*inputFile)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Start to import %s", *inputFile)

	if *isCsv {
		//
		r := csv.NewReader(strings.NewReader(string(content)))
		r.Comma = ','
		r.Comment = '#'

		for {
			data, err := r.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			var newvar gllib.GitlabVarData
			newvar.Key = *prefix + data[0]
			newvar.Value = strings.ReplaceAll(data[1], "\\n", "\n")
			if data[2] == "all" {
				newvar.Env = "*"
			} else {
				newvar.Env = data[2]
			}
			if data[3] == "true" {
				newvar.IsRaw = true
			} else {
				newvar.IsRaw = false
			}
			vars.FileData = append(vars.FileData, newvar)
		}
	} else {
		lines := strings.Split(string(content), "\n")

		for _, line := range lines {
			if len(line) == 0 {
				continue
			}
			if line[0:1] == "#" {
				continue
			}
			log.Printf("Line: \"%s\"", line)

			data := strings.SplitN(line, "=", 2)
			var newvar gllib.GitlabVarData
			newvar.Key = *prefix + data[0]
			newvar.Value = data[1]
			newvar.Env = "*"
			newvar.IsRaw = true

			vars.FileData = append(vars.FileData, newvar)
		}
	}

	log.Printf("Export to %s", *outputFile)
	json, err := json.MarshalIndent(vars.FileData, "", "  ")
	if err != nil {
		log.Println(err)
		return
	}
	if *verbose {
		log.Printf("Try to write %s file", *outputFile)
	}
	err = os.WriteFile(*outputFile, json, 0644)
	if err != nil {
		log.Println("Export to file", *outputFile, err)
		return
	}
}
