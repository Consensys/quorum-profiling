region = "ap-southeast-1"
aws_profile = "amal-sbox"
aws_user = "l1-developers/f050437"
network_name = "run2-raft-g197"
// aws instance type
instance_type = "t2.xlarge"
num_of_nodes_in_network = 4
// disk storage size for each node
volume_size = 100
vpc_id = "vpc-0fdc784976577e73b"

// consensus and block period
blockPeriod = 250
consensus = "raft"

// txpool.accountslots/gloablslots/globalqueue size at geth commandline for all nodes
txpoolSize = 50000
// geth 197 flag for specifying geth197 specific command line args
geth19 = true
// docker image of quorum
quorum_docker_image = "amalrajmani/quorum-test:v2"
// docker image of tessera
tessera_docker_image = "quorumengineering/tessera:0.11"

tps_docker_image = "amalrajmani/tpsmonitor:v7"
jmeter_docker_image = "amalrajmani/jmeter:5.2.1"

// jmeter test profile type you want to run
test_profile = "4node/deploy-contract-public"

// gas limit of the block and min/max at geth commandline for all nodes
gasLimit = 37500000

// no of threads jmeter test profile should run
no_of_threads = 1
// duration of jemetr test profile run in seconds
duration_of_run = 86400

// no of transactions to be sent per minute to quorum - for 1node and 4node jmeter test profiles
throughput = 60000

// no of transactions to be sent per minute to quorum - only applicable for custom mixed contract jmeter test profile
public_throughput = 12000
private_throughput = 2400
