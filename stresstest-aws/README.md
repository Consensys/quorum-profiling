 ## Quorum benchmark tool
 Tool to benchmark Quorum.
 This tool deploys a quorum network (based on inputs from `setting.tfvars` ) in AWS and runs a given jmeter stress test profile automatically.
 It profiles CPU & Memory usage of `geth` & `tessera` processes of the first node(`node0`) in the network.
 It also profiles TPS, total transactions count and total block count of the first node(`node0`) in the network.
 These profiles can be viewed under AWS cloudwatch > custom namespaces with namespace `<network_name>-<pulbicIp Of node0>`. 
The metric names are self explanatory.
 It creates number of nodes specified in the config and an additional node(for running jmeter test and tps monitor).
 The logs of `geth`, `tessera`, `jmeter` and `tpsmonitor` can be viewed under cloudwatch > Log groups > `/quorum/<network_name>`
 
 
 ## Configuration details (settings.tfvars)
 - `region` = aws region
 - `network_name` = network name prefix. All aws resource names of this network is prefixed with this name.
 - `instance_type` = aws instance type
 - `num_of_nodes_in_network` = number of nodes required in the network
 - `volume_size` = disk storage size of each node in the network
 - `vpc_id` = aws vpc id 
 - `gasLimit` = gasLimit of genesis block and max/min gas limit passed in geth commandline for each node
 - `blockPeriod` = block period of the consensus. units: for raft treated as milliseconds and for ibft treated as seconds
 - `txpoolSize` = initialise `geth`'s `txpool.accountqueue`,`txpool.globalslots` and `txpool.globalqueue` arguments with this txpool size for each node
 -  `geth19` = specifies if quorum is based on geth1.9.x version. This is used to specify `geth`'s commandline arguments like `--allow-insecure-unlock` that is specific to `geth1.9.x`
 - `quorum_docker_image` = quorum docker image
 - `tessera_docker_image` = tessera docker image
 - `tps_docker_image` = tpsmonitor docker image
 - `jmeter_docker_image` = jmeter docker image
 - `consensus` = consensus to be used. It should be raft or ibft
 - `test_profile` = name of the test profile to be executed
 - `no_of_threads` = number of threads per node to be created by jmeter for the specified test profile
 - `duration_of_run` = duration of run for the specified test profile
 - `throughput` = specifies the number of transactions to be sent to quorum per minute. This is used to throttle the input.
 - `private.throughput` = specifies the number of private transactions to be sent to quorum per minute. This is used to throttle the input. It is used by custom mixed test profile described below.
 - `public.throughput` = specifies the number of public transactions to be sent to quorum per minute. This is used to throttle the input. It is used by custom mixed test profile described below.
 #### Sample config:
 ```
region = "ap-southeast-1"
 network_name = "aj-dev2-test"
 instance_type = "t2.xlarge"
 num_of_nodes_in_network = 6
 volume_size = 100
 vpc_id = "vpc-a3286ec6"
 gasLimit = 200000000
 blockPeriod = 250
 txpoolSize = 50000
 #quorum_docker_image = "quorumengineering/quorum:2.5.0"
 geth19 = true
 quorum_docker_image = "amalrajmani/quorum-test:v2"
 tessera_docker_image = "quorumengineering/tessera:0.11"
 tps_docker_image = "amalrajmani/tpsmonitor:v7"
 jmeter_docker_image = " amalrajmani/jmeter:5.2.1"
 consensus = "raft"
test_profile = "4node/deploy-contract-public"
no_of_threads = 1
duration_of_run = 1200
#no of transactions to be sent per minute - for 1node and 4node test profiles
throughput = 96000

#no of transactions to be sent per minute - only applicable for custom mixed contract test profile
public_throughput = 12000
private_throughput = 2400
```
 ## Test Profiles
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

 ## Profiling - Cloudwatch metrics
 It can be viewed under AWS cloudwatch > custom namespaces with namespace `<network_name>-<pulbicIp Of node0>`. 
 The metric details are as follows:
 - `system=CpuMemMonitor`
 
 | Metric name | Description |
 | :----------- |:----------- |
 | geth-MEM% | geth memory usage |
 | geth-CPU% | geth cpu usage |
 | tm-CPU%   | cpu usage    |
 | tm-MEM%   | tessera memory usage |


 
 - `System=TpsMonitor`
 
 | Metric name | Description |
  | :----------- | :----------- |
  | TPS | transactions per second |
  | TxnCount  | total transactions count   |
  | BlockCount   | total block count |
 
 ## Usage
 - Run `terraform init` to initialize
 - To start the stress test, update `setting.tfvars` with preferred config.
 Run `terraform apply -var-file settings.tfvars`. 
 - Once testing is done, destroy the environment by running `terraform destroy -var-file settings.tfvars`.
 
 NOTE: `terraform-provider-quorum_v0.1.0` plugin is not available in terraform registry yet. You can build it from [here](https://github.com/QuorumEngineering/terraform-provider-quorum) and place under `stresstest-aws/.terraform/plugins/darwin_amd64` 
 
     