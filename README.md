## Quorum Test
> Quorum Test contains tools that can be used to run stress test in a Quorum network and measure/monitor TPS, Jmeter test and CPU/Memory usage metrics in different scenarios.

It can be used to do the following:
1. It can be used to spin up a Quorum network in AWS cloud and run stress test & measure TPS on it.  
1. It has Jmeter test profiles which can be executed manually in an existing Quorum network.
1. It has TPS Monitor which can be used to measure TPS in an existing Quorum network.
1. It can be used for development purposes to test quorum network running locally in laptop/desktop.


### Scenario1 - Testing new Quorum network in AWS with Jmeter test profile
Use this when you want to spin up a quorum network in AWS and run some tests on it and measure/monitor TPS, Jmeter test and CPU & Memory usage metrics.

Refer [testing Quorum in AWS](stresstest-aws/README.md) for more details on how to use it.

### Scenario2 - Testing an existing Quorum network with Jmeter test profile
Use this when you want to run Jmeter test profile in an existing Quorum network.

Refer [testing Quorum with JMeter Test Profiles](stresstest-aws/jmeter-test/README.md) for more details on how to use it.

### Scenario3 - Measuring TPS in an existing Quorum network
Use this when you want to measure TPS in an existing Quorum network.

Refer [measuring TPS in Quorum](tps-monitor/README.md) for more details on how to use it.


### Scenario4 - Testing locally
Use this when you want to run Jmeter test profile and measure TPS in an existing Quorum network running locally.

Refer [running locally](scripts/README.md) for more details on how to use it.


