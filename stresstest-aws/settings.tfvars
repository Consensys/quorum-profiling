region = "ap-southeast-1"
network_name = "aj-dev5-raft-g197"
instance_type = "t2.xlarge"
num_of_nodes_in_network = 4
volume_size = 100
vpc_id = "vpc-a3286ec6"

blockPeriod = 250
consensus = "raft"

txpoolSize = 50000
#quorum_docker_image = "quorumengineering/quorum:2.5.0"
geth19 = true
quorum_docker_image = "amalrajmani/quorum-test:v2"
tessera_docker_image = "quorumengineering/tessera:0.11"
tps_docker_image = "amalrajmani/tpsmonitor:v7"
jmeter_docker_image = " amalrajmani/jmeter:5.2.1"

// test profile types

#test_profile = "allnode/deploy-contract-public"
#test_profile = "allnode/deploy-contract-private"

#test_profile = "1node/deploy-contract-public"
#test_profile = "1node/deploy-contract-private"
#test_profile = "1node/update-contract-public"
#test_profile = "1node/update-contract-private"

#test_profile = "4node/deploy-contract-public"
#test_profile = "4node/deploy-contract-private"

test_profile = "custom/deploy-contract-public"
#deploy simple contract costs about 175000,
gasLimit = 70000000

no_of_threads = 4
duration_of_run = 600

#no of transactions to be sent per minute - for 1node and 4node test profiles
throughput = 96000

#no of transactions to be sent per minute - only applicable for custom mixed contract test profile
public_throughput = 12000
private_throughput = 2400
