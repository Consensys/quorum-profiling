#!/bin/bash
jmeter_test_profile=$1
consensus=$2
quorumEndpoint=$3
homeDir=/Users/maniam/go/src/github.com/QuorumEngineering/quorum-test/
cd ${homeDir}/scripts

echo "starting grafana, influxdb, prometheus.."
docker-compose up -d
echo "started grafana, influxdb and prometheus"

echo "start tps monitor..."
docker run -d -v ${homeDir}/scripts:/tpsmon -p 7777:7777 -p 2112:2112 --name tps-monitor amalrajmani/tpsmonitor:v1  --httpendpoint ${quorumEndpoint} --consensus=${consensus} --report /tpsmon/tps-report.csv --prometheusport 2112 --port 7777 --influxdb --influxdb.endpoint "http://host.docker.internal:8086" --influxdb.token "telegraf:test123"
echo "tps monitor started"

echo "start jmeter profile ${jmeter_test_profile}.."
docker run -d -v ${homeDir}/stresstest-aws/jmeter-test:/stresstest -v ${homeDir}/scripts/network.properties:/stresstest/network.properties \
-v ${homeDir}/scripts/host_acct.csv:/stresstest/host_acct.csv --name jmeter  amalrajmani/jmeter:5.2.1 -n -t /stresstest/${jmeter_test_profile}.jmx -q /stresstest/network.properties -j /stresstest/jmeter.log
echo "jmeter test profile ${jmeter_test_profile} started"
