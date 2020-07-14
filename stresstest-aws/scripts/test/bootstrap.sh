#!/bin/bash
notUp=0
for host in `cat host_acct.csv |cut -d"," -f1 | grep -v url`
do
  res=0
  tries=1
  triesCnt=60
  while [ $tries -lt $triesCnt ]
  do
    echo "checking connectivity for node $host ... tries $tries"
    res=$( curl -s -X POST --data '{"jsonrpc":"2.0","method":"eth_accounts","params":[],"id":1}' -H "Content-Type: application/json" $host:8545|grep result |wc -l )
    echo $res
    if [ $res -gt 0 ]
    then
        echo "node $host is up"
        break
    fi
    sleep 1
    tries=$( expr $tries + 1 )
  done

  if [ $res -eq 0 ]
  then
    echo "node $host is not up"
    notUp=1
    break
  fi

done

if [ $notUp -eq 1 ]
then
    echo "All nodes are not up. Fix the nodes and then start jmeter test(./start-jmeter-test.sh) and tps monitor(./start-tps.sh) manually."
else
    echo "All nodes are up"
    echo "starting grafana influxdb and prometheus..."
    sudo /usr/local/bin/docker-compose -f docker-compose-graf-inflx.yml up -d
    ./start-jmeter-test.sh
    ./start-tps.sh
fi

