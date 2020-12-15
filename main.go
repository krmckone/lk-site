package main

import (
	"io/ioutil"
	"log"
	"os"
	"strings"

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

func preproccess(md []byte, p params) []byte {
	return replaceVars(md, p)
}

func replaceVars(md []byte, p params) []byte {
	type loc struct {
		Start int
		End   int
	}
	locs := map[string]loc{}
	i := 0
	for i < len(md) {
		c1 := string(md[i])
		c2 := ""
		if i+1 < len(md) {
			c2 = string(md[i+1])
		}
		if c1 == "{" && c2 == "{" {
			start := i
			i += 2
			var v string
			for string(md[i]) != "}" && string(md[i+1]) != "}" {
				v += string(md[i])
				i++
			}
			end := i + 2
			v = strings.TrimSpace(v)
			locs[v] = loc{start, end}
			md = []byte(strings.Replace(string(md), "{{ "+v+" }}", p[v], -1))
		}
		i++
	}
	return md
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
