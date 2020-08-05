# Quorum Profiling
Quorum Profiling is a toolset built for the purpose of running stress tests on networks running on Quorum and measure/monitor the TPS and other benchmarking parameters at network level. The tests are fired using `Jmeter`. The tools can be used in different scenarios as described below:


* **Scenario 1 - Spin up a Quorum network in AWS and execute stress test:** 
In this scenario, the tool can be used to spin up a quorum network in AWS and run some tests using `Jmeter` tests and measure/monitor TPS, CPU & Memory usage metrics. A dashboard is also available to monitor the progress of `Jmeter` tests. Refer [testing Quorum in AWS](stresstest-aws/README.md) for more details on how to use it.

* **Scenario 2 - Executing stress test on an existing Quorum network with Jmeter test profile:**
In case there is an existing Quorum network already running and the sole purpose is to execute certain stress test scenarios, the `Jmeter` test profiles available with the in the tool can be executed as required.Refer [testing Quorum with JMeter Test Profiles](stresstest-aws/jmeter-test/README.md) for more details on how to use it.

* **Scenario 3 - Measuring TPS in an existing Quorum network:**
The toolset include a TPS monitoring tool which can be used to measure/monitor the TPS of an existing Quorum network
Refer [measuring TPS in Quorum](tps-monitor/README.md) for more details on how to use it.

* **Scenario 4 - Using the tool for local testing:**
The tool can be used for development purpose as well to execute tests on local Quorum network. Refer [running locally](scripts/README.md) for more details on how to use it.

## Metrics gathering and dashboards
The tool executes the stress test profile selected and then collects the following metrics:
 * CPU/memory usage metrics of `geth` & `tessera` docker containers from all the nodes in the network (using `telegraf`)
 * TPS, total transactions count and total block count metrics from first node(`node0`) in the network
 * `Jmeter` test execution metrics
 
 The above metrics are pushed to `influxdb`. The metrics can be viewed in `grafana` dashboards. Sample dashboards are as shown below:
  

 ### Sample Quorum Profiling dashboard
 
 ![Quorum Dashboard](images/quorumDashboard.jpeg) 
 
 ### Sample JMeter Dashboard
  
 ![Jmeter Dashboard](images/jmeterDashboard.jpeg) 
