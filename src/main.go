package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"strings"

	"github.com/fatih/color"
	"github.com/jkassis/pokemoncli/eztelemetry"
	"github.com/jkassis/pokemoncli/niantic"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	trace "go.opentelemetry.io/otel/trace"
)

const otTracerName = "github.com/jkassis/pokemoncli"

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
	// setup opentelemetry
	{
		// make the application resource identifier
		otResource, err := eztelemetry.NewResource(otTracerName, "v0.1.0", "demo")
		if err != nil {
			Red.Fprintf(os.Stderr, "err setting up ot exporter", err)
		}

		// make a simple exporter to the /tmp/telemetry.txt file
		f, err := os.Create("/tmp/telemetry.txt")
		if err != nil {
			Red.Fprintf(os.Stderr, "err creating ot trace file", err)
			os.Exit(1)
		}
		defer f.Close()
		otExporter, err := eztelemetry.NewPrettySimpleStreamExporter(f)
		if err != nil {
			Red.Fprintf(os.Stderr, "err setting up ot exporter", err)
		}

		// make and set the trace provider
		tp := sdktrace.NewTracerProvider(
			sdktrace.WithBatcher(otExporter),
			sdktrace.WithResource(otResource),
			sdktrace.WithSampler(sdktrace.AlwaysSample()))
		defer func() {
			if err := tp.Shutdown(context.Background()); err != nil {
				Red.Fprintf(os.Stderr, "err shutting down ot trace provider", err)
				os.Exit(1)
			}
		}()
		otel.SetTracerProvider(tp)
	}

	// intercept interrupts
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for range signalChan {
			fmt.Fprintln(os.Stderr, "\ntype 'exit' to cancel the shell")
		}
	}()

	// make the niantic api
	api := &niantic.API{
		BaseURL: "https://api.pokemontcg.io/v2/",
		Headers: map[string]string{
			"X-Api-Key": "",
		},
	}

	// welcome!
	Green.Fprintln(os.Stdout, "welcome to pokemon cli!")

	// loop forever
	background := context.Background()
	reader := bufio.NewReader(os.Stdin)
	for {
		func() {
			ctx, span := otel.Tracer(otTracerName).Start(background, "Loop") // OpenTelemetry Span
			defer span.End()

			// prompt
			Blue.Printf("your wish is my command (-h for help) > ")

			cmd, err := reader.ReadString('\n')
			span.AddEvent("User Input", trace.WithAttributes(attribute.String("cmd", cmd)))
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
				Red.Fprintf(os.Stderr, "err reading input: %v", err)
				return
			}

			// only supports 2 command right now
			if strings.HasPrefix(cmd, "exit") {
				os.Exit(0)
			} else {
				req := niantic.CardsSearchReq{}
				req.Init()
				if err := req.Parse(cmd); err != nil {
					return
				}
				fmt.Fprintf(os.Stdout, "CardSearch: %s\n", req.String())
				resp, apiErr, err := api.CardsSearch(ctx, &req)
				// we defer to the api to decide what apiErrs need logging
				// but we definitely want to log other errors
				if err != nil {
					span.RecordError(err)
					span.SetStatus(codes.Error, err.Error())
				}
				Report(resp, apiErr, err)
			}
		}()
	}
}
