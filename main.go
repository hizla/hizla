package main

import (
	_ "embed"
	"flag"
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/hizla/hizla/internal"
)

//go:embed LICENSE
var license string

type command struct {
	name, description string
}

var (
	flagVerbose bool
	commands    = []command{
		{"version", "Show hizla version"},
		{"license", "Show full license text"},
		{"help", "Show this help message"},
		{"serve", "Start the web server (internal)"},
	}
)

func main() {
	flag.BoolVar(&flagVerbose, "v", false, "Enable verbose output")
	flag.Usage = printUsage
	flag.Parse()

	log.SetFlags(0)
	log.SetPrefix("hizla: ")
	if flagVerbose {
		log.Println("Verbose mode enabled")
	}

	args := flag.Args()
	if len(args) == 0 {
		printUsage()
		os.Exit(0)
	}

	handleCommand(args)
}

func handleCommand(args []string) {
	cmd := args[0]
	if flagVerbose {
		log.Printf("Executing command: %s", cmd)
	}

	switch cmd {
	case "version":
		printVersion()
	case "license":
		fmt.Println(license)
	case "help":
		printUsage()
	case "serve":
		if flagVerbose {
			log.Println("Starting server with verbose mode")
		}
		doServe(args)
	default:
		log.Printf("unknown command: %q", cmd)
		printUsage()
		os.Exit(1)
	}
}

func printVersion() {
	if v, ok := internal.Check(internal.Version); ok {
		fmt.Println(v)
	} else if flagVerbose {
		log.Println("Impure build detected")
		fmt.Println("impure")
	} else {
		fmt.Println("impure")
	}
}

func printUsage() {
	fmt.Println("\nUsage:\thizla [-v] COMMAND [OPTIONS]")
	fmt.Println("Options:")
	fmt.Println("  -v\tVerbose output")
	fmt.Println("Commands:")

	w := tabwriter.NewWriter(os.Stdout, 0, 1, 4, ' ', 0)
	for _, cmd := range commands {
		fmt.Fprintf(w, "\t%s\t%s\n", cmd.name, cmd.description)
	}
	w.Flush()

	fmt.Println("\nUse 'hizla help' for more information")
}
