package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/cugu/tempel"
	"github.com/cugu/tempel/theme/plain"
	"os"
)

func main() {
	generateCmd(context.Background(), os.Args[1:])
}

func generateCmd(ctx context.Context, args []string) {
	cmd := flag.NewFlagSet("generate", flag.ExitOnError)
	contentPathFlag := cmd.String("content", ".", "filesystem path to content directory (default: '.')")
	outputPathFlag := cmd.String("output", "public", "filesystem path to output directory (default: 'public')")
	themeFlag := cmd.String("theme", "plain", "theme to use (default: 'plain')")
	helpFlag := cmd.Bool("help", false, "Print help and exit.")

	if cmd.Parse(args) != nil || *helpFlag {
		cmd.PrintDefaults()

		return
	}

	if *themeFlag != "plain" {
		fmt.Println("Currently only the 'plain' template is supported.") //nolint:forbidigo
		os.Exit(1)
	}

	var theme tempel.Theme

	switch *themeFlag {
	case "plain":
		theme = plain.Theme{}
	default:
		fmt.Println("Unknown template.") //nolint:forbidigo
		os.Exit(1)
	}

	if err := tempel.Generate(ctx, *contentPathFlag, theme, *outputPathFlag); err != nil {
		fmt.Println(err.Error()) //nolint:forbidigo
		os.Exit(1)
	}
}
