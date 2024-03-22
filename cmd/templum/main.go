package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/cugu/templum"
	"github.com/cugu/templum/theme/plain"
)

func main() {
	generateCmd(context.Background(), os.Args[1:])
}

func generateCmd(ctx context.Context, args []string) {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	cmd := flag.NewFlagSet("generate", flag.ExitOnError)
	contentPathFlag := cmd.String("content", ".", "filesystem path to content directory (default: '.')")
	outputPathFlag := cmd.String("output", "public", "filesystem path to output directory (default: 'public')")
	themeFlag := cmd.String("theme", "plain", "selected theme (default: 'plain')")
	baseURLFlag := cmd.String("url", ".", "base URL for the site (default: '.')")
	helpFlag := cmd.Bool("help", false, "Print help and exit.")

	if cmd.Parse(args) != nil || *helpFlag {
		cmd.PrintDefaults()

		return
	}

	var theme templum.Theme

	switch *themeFlag {
	case "plain":
		theme = plain.Theme{}
	default:
		fmt.Printf("unknown theme: %s\n", *themeFlag) //nolint:forbidigo
		os.Exit(1)
	}

	if err := templum.Generate(ctx, *baseURLFlag, *contentPathFlag, theme, *outputPathFlag); err != nil {
		fmt.Println(err.Error()) //nolint:forbidigo
		os.Exit(1)
	}
}
