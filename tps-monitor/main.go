package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/QuorumEngineering/quorum-test/tps-monitor/tpsmon"
	"gopkg.in/urfave/cli.v1"
)

var (
	app   = cli.NewApp()
	flags = []cli.Flag{
		tpsmon.ConsensusFlag,
		tpsmon.DebugFlag,
		tpsmon.TpsPortFlag,
		tpsmon.HttpEndpointFlag,
		tpsmon.ReportFileFlag,
		tpsmon.FromBlockFlag,
		tpsmon.ToBlockFlag,
		tpsmon.AwsMetricsEnabledFlag,
		tpsmon.AwsRegionFlag,
		tpsmon.AwsNwNameFlag,
		tpsmon.AwsInstanceFlag,
		tpsmon.PrometheusPortFlag,
	}
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
	app.Action = tps
	app.Before = func(c *cli.Context) error {
		log.Info("starting tps monitor")
		return nil
	}
	app.After = func(c *cli.Context) error {
		log.Info("exiting tps monitor")
		return nil
	}
	app.Flags = flags
	app.Usage = "tpsmonitor connects to a geth client (enabled with JSON-RPC endpoint) and monitors the TPS or calculate TPS for given block range"
}

func tps(ctx *cli.Context) error {
	httpendpoint := ctx.GlobalString(tpsmon.HttpEndpointFlag.Name)
	consensus := ctx.GlobalString(tpsmon.ConsensusFlag.Name)
	awsEnabled := ctx.GlobalBool(tpsmon.AwsMetricsEnabledFlag.Name)
	awsRegion := ctx.GlobalString(tpsmon.AwsRegionFlag.Name)
	awsNwName := ctx.GlobalString(tpsmon.AwsNwNameFlag.Name)
	awsInstance := ctx.GlobalString(tpsmon.AwsInstanceFlag.Name)
	debugMode := ctx.GlobalBool(tpsmon.DebugFlag.Name)
	prometheusPort := ctx.GlobalInt(tpsmon.PrometheusPortFlag.Name)
	if httpendpoint == "" {
		return errors.New("httpendpoint is empty")
	}

	if consensus == "" || (consensus != "raft" && consensus != "ibft") {
		return errors.New("invalid consensus. should be raft or ibft")
	}

	if debugMode {
		log.SetLevel(log.DebugLevel)
	}
	var awsService *tpsmon.AwsCloudwatchService
	var promethService *tpsmon.PrometheusMetricsService
	if awsEnabled {
		awsService = tpsmon.NewCloudwatchService(awsRegion, awsNwName, awsInstance)
	}

	if prometheusPort > 0 {
		promethService = tpsmon.NewPrometheusMetricsService(prometheusPort)
	}

	fromBlk := ctx.GlobalUint64(tpsmon.FromBlockFlag.Name)
	toBlk := ctx.GlobalUint64(tpsmon.ToBlockFlag.Name)
	if fromBlk > toBlk {
		log.Fatalf("from block is less than to block no")
	}

	tm := tpsmon.NewTPSMonitor(awsService, promethService, ctx.GlobalString(tpsmon.ConsensusFlag.Name) == "raft", ctx.GlobalString(tpsmon.ReportFileFlag.Name),
		fromBlk, toBlk, httpendpoint)
	startTps(tm)
	tpsPort := ctx.GlobalInt(tpsmon.TpsPortFlag.Name)
	tpsmon.NewTPSServer(tm, tpsPort)
	tm.Wait()
	return nil
}

func startTps(monitor *tpsmon.TPSMonitor) {
	if monitor.IfBlockRangeGiven() {
		go monitor.StartTpsForBlockRange()
	} else {
		monitor.StartTpsForNewBlocksFromChain()
		go func() {
			sigc := make(chan os.Signal, 1)
			signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
			defer signal.Stop(sigc)
			<-sigc
			log.Error("Got interrupt, shutting down...")
			go monitor.Stop()
			for i := 10; i > 0; i-- {
				<-sigc
				if i > 1 {
					log.Warning("Already shutting down, interrupt more to panic.", "times", i-1)
				}
			}
		}()
	}
}

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
