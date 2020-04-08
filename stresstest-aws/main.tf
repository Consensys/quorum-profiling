variable "network_name" {
  default = "amal-test"
}

variable "instance_type" {
  default = "t2.xlarge"
}

variable "volume_size" {
  type = number
  default = 100
  description = "volume size of each geth node"
}

variable "vpc_id" {

}

variable "region" {

}

variable "consensus" {
  description = "name of consensus supported by quorum. should be raft or ibft"
  default = "raft"
}

variable "test_profile" {
  default = "public-deploy-contract-all-nodes"
}

variable "no_of_threads" {
  type = number
  default = 1
  description = "number of threads to run for each thread group in jmeter test profile"
}

variable "duration_of_run" {
  type = number
  default = 600
  description = "duration of test run"
}

variable "tps_docker_image" {
  default = "amalrajmani/tpsmonitor:v1"
}

variable "jmeter_docker_image" {
  default = "amalrajmani/jmeter:5.2.1"
}

variable "quorum_docker_image" {
  default = "quorumengineering/quorum:2.4.0"
}
variable "tessera_docker_image" {
  default = "quorumengineering/tessera:0.11"
}

variable "gasLimit" {
 type = number
 description = "gasLimit to be used in genesis block and set in geth commandline"
}

variable "blockPeriod" {
  type = number
  description = "block period to be used by consensus which is specified in geth commandline"
}

variable "geth19" {
  type = bool
  description = "specifies geth version of quorum"
  default = false
}

variable "txpoolSize" {
  type = number
  description = "txpool size to be used by geth"
  default = 50000
}

variable "num_of_nodes_in_network" {
  type = number
  default = 4
  description = "number of nodes in the network"
}



locals {
  network_name              = var.network_name == "" ? basename(abspath(path.module)) : var.network_name
  generated_dir             = "build"
  tmkeys_generated_dir      = "${local.generated_dir}/${local.network_name}/tmkeys"
  accountkeys_generated_dir = "${local.generated_dir}/${local.network_name}/accountkeys"
  tm_dir_container_path     = "/data/tm"
  qdata_dir_container_path  = "/data/qdata"
  qdata_dir_vm_path         = "/data/qdata"
  tm_dir_vm_path            = "/data/tm"
  wrk_stresstest_home_path  = "~/stresstest"
  node_monitor_home_path  = "~/monitor"
  stresstest_src_path        = "stresstest"
  wrk_stresstest_gen_dir      = "${local.generated_dir}/${local.network_name}/stresstest"
  node_dir_prefix           = "node-"
  tm_dir_prefix             = "tm-"
  geth_addt_args            = ( var.geth19 ? "--allow-insecure-unlock" : "")
  number_of_nodes           = var.num_of_nodes_in_network
  quorum_docker_image       = var.quorum_docker_image
  tessera_docker_image      = var.tessera_docker_image
  tmNamedKeyAllocation = [for k in data.null_data_source.node_mapping[*].inputs: [format("A%d",k.idx)] ]

  allTmNamedKeys                = flatten(local.tmNamedKeyAllocation)
  network_cidr                  = cidrsubnet("172.16.0.0/16", 8, random_integer.additional_bits.id)
  container_raft_port           = 50400
  container_p2p_port            = 21000
  container_rpc_port            = 8545
  container_ws_port             = 8546
  container_tm_p2p_port         = 9000
  container_tm_third_party_port = 9080
  host_raft_port                = 50400
  host_p2p_port                 = 21000
  host_rpc_port                 = 8545
  host_ws_port                  = 8546
  host_tm_p2p_port              = 9000
  host_tm_third_party_port      = 9080
}

# randomize the docker network cidr
resource "random_integer" "additional_bits" {
  max = 254
  min = 1
}

data "null_data_source" "node_mapping" {
  count = var.num_of_nodes_in_network
  inputs = {
    idx = count.index
  }
}



# this file is read by acceptance test targeting this network setup
resource "local_file" "application-yml" {
  filename = format("%s/application-%s.yml", quorum_bootstrap_network.this.network_dir_abs, local.network_name)
  content  = <<-EOF
quorum:
  nodes:
%{for i in data.null_data_source.meta[*].inputs.idx~}
    ${format("Node%d:", i + 1)}
      named-privacy-address:
%{for k in local.tmNamedKeyAllocation[i]~}
        ${k}: ${element(quorum_transaction_manager_keypair.tm.*.public_key_b64, index(local.allTmNamedKeys, k))}
%{endfor~}
      url: ${format("http://%s:%d", data.null_data_source.meta[i].inputs.nodeVMPublicIP, local.host_rpc_port)}
      third-party-url: ${format("http://%s:%d", data.null_data_source.meta[i].inputs.nodeVMPublicIP, local.host_tm_third_party_port)}
%{endfor~}
EOF
}

output "network_configuration" {
  value = <<EOL

${local_file.application-yml.content}
EOL
}

output "static-nodes" {
  value = jsondecode(local_file.static-nodes[0].content)
}

output "ips" {
  value = formatlist("%s: %s - %s", aws_instance.node[*].tags.Name, aws_instance.node[*].public_dns, aws_instance.node[*].public_ip)
}

output "pvt_ips_eth_accts" {
  value = formatlist("%s,%s", aws_instance.node[*].private_ip, quorum_bootstrap_keystore.accountkeys-generator[*].account[0].address)
}

output "namedKeyAllocation" {
  value = local.tmNamedKeyAllocation
}

output "test" {
  value = aws_instance.wrk.public_dns
}

