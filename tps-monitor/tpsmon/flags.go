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

	InfluxdbEnabledFlag = cli.BoolFlag{
		Name:  "influxdb",
		Usage: `influxdb enabled`,
	}

	InfluxdbEndpointFlag = cli.StringFlag{
		Name:  "influxdb.endpoint",
		Usage: `influxdb endpoint`,
	}

	InfluxdbTokenFlag = cli.StringFlag{
		Name:  "influxdb.token",
		Usage: `influxdb token or username:password`,
		Value: ":",
	}
	InfluxdbOrgFlag = cli.StringFlag{
		Name:  "influxdb.org",
		Usage: `influxdb org name`,
		Value: "",
	}
	InfluxdbBucketFlag = cli.StringFlag{
		Name:  "influxdb.bucket",
		Usage: `influxdb bucket or database name`,
	}
	InfluxdbPointNameFlag = cli.StringFlag{
		Name:  "influxdb.point",
		Usage: `influxdb point name`,
		Value: "quorumTpsMon",
	}
	InfluxdbTagsFlag = cli.StringFlag{
		Name:  "influxdb.tags",
		Usage: `influxdb tags (comma separated list of key=value pairs)`,
		Value: "system=quorum,comp=tps",
	}
	ConsensusFlag = cli.StringFlag{
		Name:  "consensus",
		Usage: `name of consensus ("raft", "ibft")`,
	}

	DebugFlag = cli.BoolFlag{
		Name:  "debug",
		Usage: `debug mode`,
	}

	HttpEndpointFlag = cli.StringFlag{
		Name:  "httpendpoint",
		Usage: "geth chainReader's http endpoint",
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

	PrometheusPortFlag = cli.IntFlag{
		Name:  "prometheusport",
		Usage: "Enable prometheus metrics",
		Value: 0,
	}
)
