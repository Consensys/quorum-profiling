#!/bin/bash
jmeter_test_profile=""
consensus=""
quorumEndpoint=""
homeDir=""

function usage() {
  echo ""
  echo "Usage:"
  echo "    $0 --testProfile '<jmeter test profile>' --consensus <raft|ibft> --endpoint <quorum RPC endpoint> --basedir <repo base dir>"
  echo ""
  echo "Where:"
  echo "    testProfile - name of jmeter test profile. eg: 4node/deploy-contract-public)"
  echo "    consensus - name of consensus - raft or ibft. eg: raft)"
  echo "    endpoint - quorum rpc endpoint. eg: http://localhost:22000)"
  echo "    basedir - base dir of repo. eg: /home/bob/quorum-test)"
  echo ""
  exit -1
}

while (( "$#" )); do
    case "$1" in
        --testProfile)
            jmeter_test_profile=$2
            shift 2
            ;;
        --consensus)
            consensus=$2
            if [ $consensus != 'raft' ] && [ $consensus != 'ibft' ]
            then
            echo "consensus must be raft or ibft"
            usage
            fi
            shift 2
            ;;
        --endpoint)
            quorumEndpoint=$2
            shift 2
            ;;
        --basedir)
            homeDir=$2
            shift 2
            ;;
        --help)
            shift
            usage
            ;;
        *)
            echo "Error: Unsupported command line parameter $1"
            usage
            ;;
    esac
done

echo "homeDir=$homeDir"
echo "consensus=$consensus"
echo "endpoint=$quorumEndpoint"
echo "testProfile=$jmeter_test_profile"

cd ${homeDir}/scripts

# copy jmeter test profiles to jmeter for docker volume mapping
cp -pR ${homeDir}/stresstest-aws/jmeter-test/* ${homeDir}/scripts/jmeter

echo "starting grafana, influxdb, prometheus.."
docker-compose up -d
echo "started grafana, influxdb and prometheus"

# TODO - change this to quorumengineering namespace once open source is completed and docker image pushed to quorumengineering
tpsDockerImg="docker.pkg.github.com/quorumengineering/quorum-test/tpsmonitor:latest"
jmeterDockerImg="docker.pkg.github.com/quorumengineering/quorum-test/jmeter:5.2.1"

echo "start jmeter profile ${jmeter_test_profile}.."
docker run -d -v ${homeDir}/scripts/jmeter:/stresstest -v ${homeDir}/scripts/jmeter/network.properties:/stresstest/network.properties \
-v ${homeDir}/scripts/jmeter/host_acct.csv:/stresstest/host_acct.csv --name jmeter  ${jmeterDockerImg} -n -t /stresstest/${jmeter_test_profile}.jmx -q /stresstest/network.properties -j /stresstest/jmeter.log
echo "jmeter test profile ${jmeter_test_profile} started"


echo "start tps monitor..."
docker run -d -v ${homeDir}/scripts/tpsmon:/tpsmon -p 7777:7777 -p 2112:2112 --name tps-monitor ${tpsDockerImg}  --httpendpoint ${quorumEndpoint} --consensus=${consensus} --report /tpsmon/tps-report.csv --prometheusport 2112 --port 7777 --influxdb --influxdb.endpoint "http://host.docker.internal:8086" --influxdb.token "telegraf:test123"
echo "tps monitor started"
