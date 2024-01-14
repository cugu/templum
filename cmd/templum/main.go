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
	baseURLFlag := cmd.String("url", ".", "base URL for the site (default: '.')")
	helpFlag := cmd.Bool("help", false, "Print help and exit.")

	if cmd.Parse(args) != nil || *helpFlag {
		cmd.PrintDefaults()

		return
	}

	if err := templum.Generate(ctx, *baseURLFlag, *contentPathFlag, plain.Theme{}, *outputPathFlag); err != nil {
		fmt.Println(err.Error()) //nolint:forbidigo
		os.Exit(1)
	}
}
