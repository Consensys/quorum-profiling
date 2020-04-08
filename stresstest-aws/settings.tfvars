region = "ap-southeast-1"
network_name = "aj-dev3-test"
instance_type = "t2.xlarge"
num_of_nodes_in_network = 7
volume_size = 100
vpc_id = "vpc-a3286ec6"
gasLimit = 200000000
blockPeriod = 1
txpoolSize = 50000
#quorum_docker_image = "quorumengineering/quorum:2.5.0"
geth19 = true
quorum_docker_image = "amalrajmani/quorum-test:v2"
tessera_docker_image = "quorumengineering/tessera:0.11"
tps_docker_image = "amalrajmani/tpsmonitor:v7"
jmeter_docker_image = " amalrajmani/jmeter:5.2.1"
consensus = "ibft"

// test profile types

#test_profile = "allnode/deploy-contract-public"
#test_profile = "allnode/deploy-contract-private"

#test_profile = "1node/deploy-contract-public"
#test_profile = "1node/deploy-contract-private"
#test_profile = "1node/update-contract-public"
#test_profile = "1node/update-contract-private"

#test_profile = "4node/deploy-contract-public"
#test_profile = "4node/deploy-contract-private"

test_profile = "4node/deploy-contract-public"

no_of_threads = 1
duration_of_run = 36000