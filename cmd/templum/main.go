package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/cugu/templum"
	"github.com/cugu/templum/theme/plain"
)

func main() {
	generateCmd(context.Background(), os.Args[1:])
}

func generateCmd(ctx context.Context, args []string) {
	cmd := flag.NewFlagSet("generate", flag.ExitOnError)
	contentPathFlag := cmd.String("content", ".", "filesystem path to content directory (default: '.')")
	outputPathFlag := cmd.String("output", "public", "filesystem path to output directory (default: 'public')")
	configPathFlag := cmd.String("config", "config.yaml", "filesystem path to config file, relative to the content directory (default: 'config.yaml')")
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

	var theme templum.Theme

	switch *themeFlag {
	case "plain":
		theme = plain.Theme{}
	default:
		fmt.Println("Unknown template.") //nolint:forbidigo
		os.Exit(1)
	}

	if err := templum.Generate(ctx, *configPathFlag, *contentPathFlag, theme, *outputPathFlag); err != nil {
		fmt.Println(err.Error()) //nolint:forbidigo
		os.Exit(1)
	}
}
