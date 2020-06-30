## Quorum Test
> Quorum Test is intended to be used for stress testing Quorum

Quorum Test suite can be used to do the following:
* It can be used to bootstrap a Quorum network in AWS cloud and run stress test & measure TPS on it.  
* It has Jmeter test profiles which can be executed manually in a Quorum network.
* It has TPS Monitor which can be used to measure TPS in a Quorum network.


### Prerequisites
Download Terraform runtime to your machine:
* From [HashiCorp website](https://www.terraform.io/downloads.html)
* MacOS: `brew install terraform`

Download Jmeter to your machine:
* From [Jmeter website](https://jmeter.apache.org/download_jmeter.cgi)
* MacOS: `brew install jmeter`

### Getting Started
* [Testing in AWS](stresstest-aws/README.md): 
It deploys (using terraform) a quorum network in AWS and runs Jmeter test & TPS monitor in AWS automatically.
* [JMeter Test Profiles](stresstest-aws/jmeter-test/README.md): It contains Jmeter test profiles that can be executed manually on a quorum network.
* [TPS Monitor](tps-monitor/README.md): TPS monitor process can be used to measure TPS in a quorum network.
