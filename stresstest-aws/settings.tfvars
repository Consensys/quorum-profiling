aws_region = "<aws-region>"
aws_profile = "<aws-profile>"
aws_user = "<aws-user-name>"
aws_network_name = "<aws-network-name>"
// aws instance type, example: t3a.2xlarge
aws_instance_type = "<aws-instance-type>"
// no of nodes in Quorum network
aws_num_of_nodes_in_network = 4
// disk storage size for each node in GB
aws_volume_size = 100
aws_vpc_id = "<aws-vpc-id>"

// consensus and block period
blockPeriod = 250
// consensus - raft or ibft or clique
consensus = "raft"
// flag to enable/disable tessera
enable_tessera = true
// flag that indicates upstream geth or quorum geth
is_quorum = false

// txpool.accountslots/gloablslots/globalqueue size at geth commandline for all nodes
txpoolSize = 50000
// geth 197 flag for specifying geth197 specific command line args
geth19 = true
// docker image of Quorum
quorum_docker_image = "quorumengineering/quorum:latest"
// docker image of tessera
tessera_docker_image = "quorumengineering/tessera:0.11"

tps_docker_image = "quorumengineering/tpsmonitor:v1"
jmeter_docker_image = "quorumengineering/jmeter:5.2.1"

// jmeter test profile type you want to run
jmeter_test_profile = "4node/deploy-contract-public"

enable_tessera = false

// gas limit of the block and min/max at geth commandline for all nodes
gasLimit = 37500000

// no of threads jmeter test profile should run
jmeter_no_of_threads = 1
// duration of jemetr test profile run (in seconds)
jmeter_duration_of_run = 1800

// no of transactions to be sent per minute to quorum - for 1node and 4node jmeter test profiles
jmeter_throughput = 60000

// no of transactions to be sent per minute to quorum - only applicable for custom mixed contract jmeter test profile
jmeter_public_throughput = 12000
jmeter_private_throughput = 2400
