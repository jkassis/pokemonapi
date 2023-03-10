package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"strings"

	"github.com/fatih/color"
	"github.com/jkassis/pokemoncli/niantic"
)

const tracerName = "github.com/jkassis/pokemoncli"

var Blue = color.New(color.FgBlue)
var Red = color.New(color.FgRed)
var Yellow = color.New(color.FgYellow)
var Green = color.New(color.FgGreen)
var White = color.New(color.FgWhite)
var Spaces = regexp.MustCompile(`\s+`)

func Report(resp niantic.Res, apiErr niantic.Err, err error) {
	if err != nil {
		fmt.Println(err.Error())
	} else if apiErr != nil {
		out, err := json.MarshalIndent(apiErr, "", "  ")
		if err != nil {
			fmt.Printf("error marshaling response in test code: %v\n", err)
		} else {
			fmt.Println(string(out))
		}
	} else {
		out, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			fmt.Printf("error marshaling response in test code: %v\n", err)
		} else {
			fmt.Println(string(out))
		}
	}
}

func main() {
	// intercept interrupts
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for range signalChan {
			fmt.Fprintln(os.Stderr, "\ntype 'exit' to cancel the shell")
		}
	}()

	// make the api
	api := &niantic.API{
		BaseURL: "https://api.pokemontcg.io/v2/",
		Headers: map[string]string{
			"X-Api-Key": "",
		},
	}

	// welcome!
	Green.Fprintln(os.Stdout, "welcome to pokemon cli!")

	// loop forever
	reader := bufio.NewReader(os.Stdin)
	for {
		Blue.Printf("your wish is my command (-h for help) > ")

		cmd, err := reader.ReadString('\n')
		if err != nil {
			Red.Fprintf(os.Stderr, "err reading input: %v", err)
			continue
		}

		// only supports 2 command right now
		if strings.HasPrefix(cmd, "exit") {
			os.Exit(0)
		} else {
			req := niantic.CardsSearchReq{}
			req.Init()
			err := req.Parse(cmd)
			if err != nil {
				continue
			}
			resp, apiErr, err := api.CardsSearch(&req)
			Report(resp, apiErr, err)
		}
	}
}
