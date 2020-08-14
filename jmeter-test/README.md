# `Jmeter` test profiles
 Jmeter is used to run the test profiles. 
 The objective is to test performance and robustness of quorum under high volume of transactions over a longer duration.
To achieve this, the tests are designed in such a way to simulate high volume of transactions flowing into quorum from different nodes in the network concurrently. 
Each test profile creates one or more threads (param: `no_of_threads`) per node and runs those threads concurrently. Each thread sends transactions(create/update SimpleStorage contract as public/private) to a specified node(from default account of that node) continuously for a given period of time (param: `duration_of_run`). 
Private transactions have only one participants in `privateFor` by default.
 The private transactions are executed in each node in such a way that:
  - transaction submitted in `node0` is between `node0` & `node1`
  - transaction submitted in `node1` is between `node1` & `node2`
  - ... 
  - transaction submitted in `node[N]` is between `node[N]` and `node0`
  
 Private transactions have only one participant in `privateFor` field.
 Once the test has finished running, charts on profiles(like cpu & memory usage, tps, etc) can be saved down from AWS cloudwatch.
 
 |Profile No | Test profile name | Transaction | Description |
 | --------- | ----------------- | ----------- | ----------- |
 |1| `1node/deploy-contract-public` | create simpleStorage public contract (with constructor initialised to random number) | creates specified no of threads for first node. sends transactions to first node only. |
 |2| `1node/deploy-contract-private` | create simpleStorage private contract (with constructor initialised to random number) | same as profile `1` |
 |3| `1node/update-contract-public` | update simpleStorage public contract (with setter initialised to random number) | same as profile `1` |
 |4| `1node/update-contract-private` | update simpleStorage private contract (with setter initialised to random number) | same as profile `1` |
 |5| `4node/deploy-contract-public` | create simpleStorage public contract (with constructor - initialised to random number)| creates specified no of threads for each node (first 4 nodes only). sends transactions to first 4 nodes only. |
 |6| `4node/deploy-contract-private` | create simpleStorage private contract (with constructor - initialised to random number)| same as profile `5` |
 |7| `custom/deploy-contract-public` | create simpleStorage public contract (with constructor - initialised to random number)| creates specified no of threads and each thread will work on one of the nodes specified in the `.csv` input file. |
 |8| `custom/deploy-contract-private` | create simpleStorage private contract (with constructor - initialised to random number)| same as profile `7` |
 |9| `custom/deploy-mixed-contract` | create simpleStorage private & public contract (with constructor - initialised to random number)| creates specified no of thread pairs and each thread pair will work on one of the nodes specified in the `.csv` input file sending private and public transactions concurrently.  |

## Prerequisites
Refer to **scenario 2** [here](../README.md#prerequisites-for-test-execution) for all prerequisites. 
 
## Executing tests
 For further details of the the profiles and how these can be executed, please refer to:
 
 * profiles for executing the stress test from single node, refer [here](1node/)
 * profiles for executing the stress test from four nodes in the network, refer [here](4node/)
 * If you want to customize and execute the stress test from any number of nodes in the network, refer [here](custom/)
 

## Publish jmeter test results
 All jmeter test profiles have a backend listener which can be configured so that test results would be published to an `influxdb` instance as given in input.
 
  
## Disabling `influxDB` 
 In order to run tests without metrics being pushed to `influxDB`, disable the listener by replacing the below line in the test plan jmx files
 
``<BackendListener guiclass="BackendListenerGui" testclass="BackendListener" testname="Backend Listener" enabled="false">``

 
