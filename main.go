package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"unicode"

	"github.com/gomarkdown/markdown"
	"gopkg.in/yaml.v2"
)

type config struct {
	Template templateConfig `yaml:"template"`
}

type templateConfig struct {
	Params params `yaml:"params"`
}

type params map[string]string

func readConfig(path string) config {
	b := readFile(path)
	config := config{}
	err := yaml.Unmarshal(b, &config)
	if err != nil {
		log.Fatal(err)
	}
	return config
}

func readFile(path string) []byte {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return []byte{}
	}
	return bytes
}

func writeFile(path string, b []byte) {
	if err := ioutil.WriteFile(path, b, 0644); err != nil {
		log.Fatal(err)
	}
}

func mkdir(path string) {
	if err := os.Mkdir("build", 0755); err != nil {
		if !os.IsNotExist(err) {
			log.Fatal(err)
		}
	}
}

// Another equivalent version of this logic when parsing
// {{x}} vars from left to right:
// return string(b1) == "}" ||  string(b2) == "}"
// However this is incorrect when there is an error of the form
// {{x} i.e. unmatched right bracket, so I prefer the below version.
func inVar(b1, b2 byte) bool {
	if string(b1) == "}" {
		return string(b2) != "}"
	}
	return true
}

func preproccess(md []byte, p params) []byte {
	return replaceVars(md, p)
}

func replaceVars(md []byte, p params) []byte {
	type loc struct {
		Start int
		End   int
	}
	mdStr := string(md)
	locs := map[string]loc{}
	for i := 0; i < len(mdStr); i++ {
		c1 := string(mdStr[i])
		c2 := ""
		if i+1 < len(mdStr) {
			c2 = string(mdStr[i+1])
		}
		if c1 == "{" && c2 == "{" {
			start := i
			i += 2
			vBytes := []byte{}
			for ; inVar(mdStr[i], mdStr[i+1]); i++ {
				if unicode.IsSpace(rune(mdStr[i])) {
					continue
				}
				vBytes = append(vBytes, mdStr[i])
			}
			v := string(vBytes)
			end := i + 2
			locs[v] = loc{start, end}
			vFmt := fmt.Sprintf("{{%s}}", v)
			mdStr = strings.Replace(mdStr, mdStr[start:end], vFmt, 1)
			mdStr = strings.Replace(mdStr, vFmt, p[v], -1)
		}
	}
	return []byte(mdStr)
}

func main() {
	if err := os.RemoveAll("build"); err != nil {
		log.Fatal(err)
	}
	config := readConfig("resources/config/config.yml")

	md := readFile("resources/index.md")
	md = preproccess(md, config.Template.Params)
	output := markdown.ToHTML(md, nil, nil)
	mkdir("build")
	writeFile("build/index.html", output)
}
