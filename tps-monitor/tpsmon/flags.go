package tpsmon

import "gopkg.in/urfave/cli.v1"

var (
	AwsMetricsEnabledFlag = cli.BoolFlag{
		Name:  "awsmetrics",
		Usage: `aws metrics enabled`,
	}

	AwsRegionFlag = cli.StringFlag{
		Name:  "awsregion",
		Usage: `aws region`,
	}

	AwsNwNameFlag = cli.StringFlag{
		Name:  "awsnetwork",
		Usage: `aws network name`,
	}

	AwsInstanceFlag = cli.StringFlag{
		Name:  "awsinst",
		Usage: `aws instance name`,
	}

	ConsensusFlag = cli.StringFlag{
		Name:  "consensus",
		Usage: `name of consensus ("raft", "ibft")`,
	}

	WSEndpointFlag = cli.StringFlag{
		Name:  "wsendpoint",
		Usage: "geth clients WS endpoint",
	}

	ReportFileFlag = cli.StringFlag{
		Name:  "report",
		Usage: "full path of the file to write the report",
		Value: "tps-report.csv",
	}

	TpsPortFlag = cli.IntFlag{
		Name:  "port",
		Usage: "port for tps monitor",
		Value: 7575,
	}

	FromBlockFlag = cli.Uint64Flag{
		Name:  "from",
		Usage: "from block no",
		Value: 0,
	}

	ToBlockFlag = cli.Uint64Flag{
		Name:  "to",
		Usage: "to block no",
		Value: 0,
	}
)
