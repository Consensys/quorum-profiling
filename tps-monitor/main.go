package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/QuorumEngineering/tps-monitor/tpsmon"
	"github.com/ethereum/go-ethereum/console"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"gopkg.in/urfave/cli.v1"
)

var (
	app   = cli.NewApp()
	flags = []cli.Flag{
		tpsmon.WSEndpointFlag,
		tpsmon.ConsensusFlag,
		tpsmon.TpsPortFlag,
		tpsmon.ReportFileFlag,
		tpsmon.FromBlockFlag,
		tpsmon.ToBlockFlag,
		tpsmon.AwsMetricsEnabledFlag,
		tpsmon.AwsRegionFlag,
		tpsmon.AwsNwNameFlag,
		tpsmon.AwsInstanceFlag,
	}
)

func init() {
	app.Action = tps
	app.Before = func(c *cli.Context) error {
		log.Print("starting tps monitor")
		return nil
	}
	app.After = func(c *cli.Context) error {
		log.Print("exiting tps monitor")
		console.Stdin.Close()
		return nil
	}
	app.Flags = flags
	app.Usage = "tpsmonitor connects to a geth client (enabled with WS endpoint) and monitors the TPS or calculate TPS for given block range"
}

func tps(ctx *cli.Context) error {
	wsendpoint := ctx.GlobalString(tpsmon.WSEndpointFlag.Name)
	consensus := ctx.GlobalString(tpsmon.ConsensusFlag.Name)
	awsEnabled := ctx.GlobalBool(tpsmon.AwsMetricsEnabledFlag.Name)
	awsRegion := ctx.GlobalString(tpsmon.AwsRegionFlag.Name)
	awsNwName := ctx.GlobalString(tpsmon.AwsNwNameFlag.Name)
	awsInstance := ctx.GlobalString(tpsmon.AwsInstanceFlag.Name)
	var awsCfg *tpsmon.AwsCloudwatchService
	if awsEnabled {
		awsCfg = tpsmon.NewCloudwatchService(awsRegion, awsNwName, awsInstance)
	}

	if wsendpoint == "" {
		return errors.New("wsendpoint is empty")
	}

	if consensus == "" || (consensus != "raft" && consensus != "ibft") {
		return errors.New("invalid consensus. should be raft or ibft")
	}

	log.Printf("connecting to %s", wsendpoint)
	client, err := ethclient.Dial(wsendpoint)
	if err != nil {
		log.Fatal(err)
	}

	headers := make(chan *types.Header)
	sub, err := client.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		return err
	}

	fromBlk := ctx.GlobalUint64(tpsmon.FromBlockFlag.Name)
	toBlk := ctx.GlobalUint64(tpsmon.ToBlockFlag.Name)
	if fromBlk > toBlk {
		log.Fatalf("from block is less than to block no")
	}

	tm := tpsmon.NewTPSMonitor(awsCfg, ctx.GlobalString(tpsmon.ConsensusFlag.Name) == "raft", ctx.GlobalString(tpsmon.ReportFileFlag.Name),
		fromBlk, toBlk, sub, headers, client)
	startTps(tm)
	tpsPort := ctx.GlobalInt(tpsmon.TpsPortFlag.Name)
	tpsmon.NewTPSServer(tm, tpsPort)
	tm.Wait()
	return nil
}

func startTps(monitor *tpsmon.TPSMonitor) {
	monitor.Start()
	go func() {
		sigc := make(chan os.Signal, 1)
		signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
		defer signal.Stop(sigc)
		<-sigc
		log.Print("Got interrupt, shutting down...")
		go monitor.Stop()
		for i := 10; i > 0; i-- {
			<-sigc
			if i > 1 {
				log.Print("WARN: Already shutting down, interrupt more to panic.", "times", i-1)
			}
		}
	}()
}

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
