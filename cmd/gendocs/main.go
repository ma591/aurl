package main

import (
	"fmt"
	"os"

	"github.com/shawnpana/aurl/cmd"
	cobradoc "github.com/spf13/cobra/doc"
)

func main() {
	root := cmd.GetRootCmd()

	// Generate man page
	header := &cobradoc.GenManHeader{
		Title:   "ARC",
		Section: "1",
		Source:  "aurl",
		Manual:  "Aurl Manual",
	}
	if err := cobradoc.GenManTree(root, header, "docs"); err != nil {
		fmt.Fprintf(os.Stderr, "Error generating man pages: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Generated man pages in docs/")

	// Generate shell completions
	if err := root.GenZshCompletionFile("docs/completions/aurl.zsh"); err != nil {
		fmt.Fprintf(os.Stderr, "Error generating zsh completion: %v\n", err)
		os.Exit(1)
	}

	if err := root.GenBashCompletionFileV2("docs/completions/aurl.bash", true); err != nil {
		fmt.Fprintf(os.Stderr, "Error generating bash completion: %v\n", err)
		os.Exit(1)
	}

	if err := root.GenFishCompletionFile("docs/completions/aurl.fish", true); err != nil {
		fmt.Fprintf(os.Stderr, "Error generating fish completion: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Generated shell completions in docs/completions/")
}
