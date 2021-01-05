package main

import (
	"github.com/krmckone/ksite/internal/templater"
)

func main() {
	if err := templater.BuildSite(); err != nil {
		panic(err)
	}
}
