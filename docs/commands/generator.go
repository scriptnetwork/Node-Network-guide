package main

import (
	"log"
	"strings"

	script "github.com/scripttoken/script/cmd/script/cmd"
	scriptcli "github.com/scripttoken/script/cmd/scriptcli/cmd"
	"github.com/spf13/cobra/doc"
)

func generateScriptCLIDoc(filePrepender, linkHandler func(string) string) {
	var all = scriptcli.RootCmd
	err := doc.GenMarkdownTreeCustom(all, "./wallet/", filePrepender, linkHandler)
	if err != nil {
		log.Fatal(err)
	}
}

func generateScriptDoc(filePrepender, linkHandler func(string) string) {
	var all = script.RootCmd
	err := doc.GenMarkdownTreeCustom(all, "./ledger/", filePrepender, linkHandler)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	filePrepender := func(filename string) string {
		return ""
	}

	linkHandler := func(name string) string {
		return strings.ToLower(name)
	}

	generateScriptCLIDoc(filePrepender, linkHandler)
	generateScriptDoc(filePrepender, linkHandler)
	Walk()
}
