package main

import (
	_ "embed"
	"fmt"
	"git.gensokyo.uk/security/fortify/command"
	"github.com/hizla/hizla/internal"
	"log"
	"log/slog"
	"os"
)

//go:embed LICENSE
var license string

func main() {
	log.SetFlags(0)
	log.SetPrefix("hizla: ")

	var flagVerbose bool

	c := command.New(os.Stderr, log.Printf, "hizla", func(args []string) error {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		return nil
	}).Flag(&flagVerbose, "v", command.BoolFlag(false), "Verbose output")

	c.Command("version", "Show hizla version", func([]string) error {
		if v, ok := internal.Check(internal.Version); ok {
			fmt.Println(v)
		} else {
			fmt.Println("impure")
		}
		return nil
	}).Command("license", "Show full license text", func([]string) error {
		fmt.Println(license)
		return nil
	}).Command("help", "Show this help message", func([]string) error {
		c.PrintHelp()
		return nil
	}).Command("serve", command.UsageInternal, func(args []string) error {
		doServe(args)
		return nil
	})

	c.MustParse(os.Args[1:], func(err error) {
		if err != nil {
			log.Printf("error: %v", err)
			os.Exit(1)
		}
		os.Exit(0)
	})
}
