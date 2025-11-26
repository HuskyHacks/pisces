package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/mjc-gh/pisces"
	"github.com/mjc-gh/pisces/engine"
	"github.com/mjc-gh/pisces/internal/browser"
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v3"
)

var logger *zerolog.Logger

type taskCallbackFn = func(*cli.Command, *engine.Engine) error

func main() {
	baseArgs := []cli.Argument{
		&cli.StringArgs{Name: "url", Min: 1, Max: -1},
	}

	baseFlags := []cli.Flag{
		&cli.BoolFlag{Name: "debug", Aliases: []string{"d"}},
		&cli.BoolFlag{Name: "remote", Aliases: []string{"r"}},
		&cli.IntFlag{Name: "concurrency", Aliases: []string{"c"}},
		&cli.IntFlag{Name: "port", Value: 9222},
		&cli.StringFlag{Name: "device-type", Value: "desktop", Action: validDeviceType},
		&cli.StringFlag{Name: "device-size", Value: "large", Action: validDeviceSize},
		&cli.StringFlag{Name: "host", Value: "127.0.0.1"},
		&cli.StringFlag{Name: "user-agent", Value: "chrome"},
	}

	withOutputFlags := append([]cli.Flag{
		&cli.StringFlag{Name: "output", Value: "pisces.ndjson", Aliases: []string{"o"}},
	}, baseFlags...)

	cmd := &cli.Command{
		Name:    "pisces",
		Version: "0.0.0",
		Usage:   "A tool for analyzing phishing sites",
		Commands: []*cli.Command{
			{
				Name:      "analyze",
				Usage:     "Analyze and interact one or more URLs for phishing",
				Arguments: baseArgs,
				Flags:     withOutputFlags,
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return runTask(ctx, cmd, "analyze", outputResultJson)
				},
			},
			{
				Name:      "collect",
				Usage:     "Collect HTML and assets for one or more URLs",
				Arguments: baseArgs,
				Flags:     withOutputFlags,
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return runTask(ctx, cmd, "collect", outputResultJson)
				},
			},
			{
				Name:      "screenshot",
				Usage:     "Screenshot one or more URLs",
				Arguments: baseArgs,
				Flags:     baseFlags,
				Action: func(ctx context.Context, cmd *cli.Command) error {
					//return runTask(ctx, cmd, "screenshot")
					return nil
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func runTask(ctx context.Context, cmd *cli.Command, name string, callback taskCallbackFn) error {
	logger = pisces.SetupLogger(cmd.Bool("debug"))

	deviceSize := cmd.StringArg("device-size")
	deviceType := cmd.StringArg("device-type")
	host := cmd.String("host")
	port := cmd.Int("port")
	urls := cmd.StringArgs("url")

	opts := []engine.Option{engine.WithLogger(pisces.Logger())}

	if cmd.Bool("remote") && host != "" && port != 0 {
		opts = append(opts, engine.WithRemoteAllocator(host, port))
	}

	e := engine.New(cmd.Int("concurrency"), opts...)
	e.Start(ctx)

	for _, url := range urls {
		t := engine.NewTask(name, url)
		t.SetDevice(deviceType, deviceSize)
		t.SetUserAgent(deviceType, cmd.StringArg("user-agent"))

		e.Add(t)
	}

	go e.Shutdown()

	// TODO handle interrupt signal and wait for shutdown

	return callback(cmd, e)
}

func outputResultJson(cmd *cli.Command, e *engine.Engine) error {
	output := cmd.String("output")
	out, err := os.Create(output)
	if err != nil {
		panic(err)
	}

	for r := range e.Results() {
		if r.Error != nil {
			logger.Warn().Msgf("result error: %v", r.Error)
			continue
		}

		logger.Info().Msgf("result for %s (duration %s)", r.URL, r.Elapsed.String())

		line, err := json.Marshal(r.Result)
		if err != nil {
			logger.Warn().Msgf("result json marshal error: %v", r.Error)
		}

		if _, err = out.Write(line); err != nil {
			logger.Warn().Msgf("result json write error: %v", r.Error)
		}

		if _, err = out.Write([]byte("\n")); err != nil {
			logger.Warn().Msgf("result json write error: %v", r.Error)
		}

		logger.Debug().Msgf("wrote to file %s", output)
	}

	if err = out.Close(); err != nil {
		logger.Warn().Msgf("file close error: %v", err)
	}

	logger.Info().Msg("done")

	return nil
}

func validDeviceType(ctx context.Context, cmd *cli.Command, v string) error {
	if !browser.IsValidDeviceType(v) {
		return fmt.Errorf("invalid device type: %v", v)
	}

	return nil
}

func validDeviceSize(ctx context.Context, cmd *cli.Command, v string) error {
	if !browser.IsValidDeviceSize(v) {
		return fmt.Errorf("invalid device size: %v", v)
	}

	return nil
}
