#!/bin/bash
for host in `cat host_acct.csv |cut -d"," -f1 | grep -v url`
do
  counter=0
  while [ $counter -eq 0 ]
  do
    echo "checking connectivity for node $host ..."
    counter=$( curl -s -X POST --data '{"jsonrpc":"2.0","method":"eth_accounts","params":[],"id":1}' -H "Content-Type: application/json" $host:8545|grep result |wc -l )
    echo $counter
    sleep 1
  done
  echo "node $host is up"
done
echo "all nodes are up"

./start-jmeter-test.sh

./start-tps.sh
