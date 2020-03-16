#!/bin/bash

cd /home/ubuntu/stresstest/jmeter
nohup ./apache-jmeter-5.2.1/bin/jmeter.sh -n -t private-contract-creation-4node.jmx >stress-test.log 2>&1 &
