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

var (
	flagVerbose bool

	//go:embed LICENSE
	license string
)

func init() {
	flag.BoolVar(&flagVerbose, "v", false, "Verbose output")
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("hizla: ")

	flag.CommandLine.Usage = func() {
		fmt.Println()
		fmt.Println("Usage:\thizla [-v] COMMAND [OPTIONS]")
		fmt.Println()
		fmt.Println("Commands:")
		w := tabwriter.NewWriter(os.Stdout, 0, 1, 4, ' ', 0)
		commands := [][2]string{
			{"version", "Show hizla version"},
			{"license", "Show full license text"},
			{"help", "Show this help message"},
		}
		for _, c := range commands {
			_, _ = fmt.Fprintf(w, "\t%s\t%s\n", c[0], c[1])
		}
		if err := w.Flush(); err != nil {
			log.Printf("cannot write command list: %v", err)
		}
		fmt.Println()
	}
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		flag.CommandLine.Usage()
		os.Exit(0)
	}

	switch args[0] {
	case "version": // print comp version string
		if v, ok := internal.Check(internal.Version); ok {
			fmt.Println(v)
		} else {
			fmt.Println("impure")
		}
		os.Exit(0)
	case "license": // print embedded license
		fmt.Println(license)
		os.Exit(0)
	case "help": // print help message
		flag.CommandLine.Usage()
		return

	// internal commands
	case "serve":
		doServe(args)
		os.Exit(0)
	}
}
