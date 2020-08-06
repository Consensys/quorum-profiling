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
  echo "    basedir - base dir of repo. eg: /home/bob/quorum-profiling)"
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

if [ ! -d "${homeDir}/scripts" ]; then
    echo "error: ${homeDir}/scripts does not exist"
    exit
fi

cd ${homeDir}/scripts

# create the tpsmon directory to store the test results
mkdir -p tpsmon jmeter

# check that properties file are there in the local directory
HOSTFILE="${homeDir}/scripts/host_acct.csv"
PROPERTIES="${homeDir}/scripts/network.properties"

if [ ! -f "${HOSTFILE}" ]; then
    echo "error: configuration file ${HOSTFILE} not found"
    exit
fi

if [ ! -f "${PROPERTIES}" ]; then
    echo "error: configuration file ${PROPERTIES} not found"
    exit
fi

# copy jmeter test profiles to jmeter for docker volume mapping
JMXFILE=${homeDir}/jmeter-test/${JMXDIR}/${jmeter_test_profile}.jmx

if [ ! -f "${JMXFILE}" ]; then
    echo "error: jmeter test script ${JMXFILE} not found"
    exit
fi

JMXDIR=`echo ${jmeter_test_profile} | cut -f1 -d "/"`
mkdir -p ${homeDir}/scripts/jmeter/${JMXDIR}

cp ${homeDir}/jmeter-test/${jmeter_test_profile}.jmx ${homeDir}/scripts/jmeter/${JMXDIR}
cp $HOSTFILE ${homeDir}/scripts/jmeter
cp $PROPERTIES ${homeDir}/scripts/jmeter

echo "starting grafana, influxdb, prometheus.."
docker-compose up -d
echo "started grafana, influxdb and prometheus"

# TODO - change this to quorumengineering namespace once open source is completed and docker image pushed to quorumengineering
tpsDockerImg="amalrajmani/tpsmonitor:v1"
jmeterDockerImg="amalrajmani/jmeter:5.2.1"

echo "start jmeter profile ${jmeter_test_profile}.."
docker run -d -v ${homeDir}/scripts/jmeter:/stresstest \
 --name jmeter  ${jmeterDockerImg} -n -t /stresstest/${jmeter_test_profile}.jmx -q /stresstest/network.properties -j /stresstest/jmeter.log
echo "jmeter test profile ${jmeter_test_profile} started"


echo "start tps monitor..."
docker run -d -v ${homeDir}/scripts/tpsmon:/tpsmon -p 7777:7777 -p 2112:2112 --name tps-monitor ${tpsDockerImg}  --httpendpoint ${quorumEndpoint} --consensus=${consensus} --report /tpsmon/tps-report.csv --prometheusport 2112 --port 7777 --influxdb --influxdb.endpoint "http://host.docker.internal:8086" --influxdb.token "telegraf:test123"
echo "tps monitor started"
