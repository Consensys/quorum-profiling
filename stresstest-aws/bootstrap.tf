data "null_data_source" "meta" {
  count = local.number_of_nodes
  inputs = {
    peerJson             = "{ \"url\": \"http://${aws_instance.node[count.index].private_ip}:${local.host_tm_p2p_port}\" }"
    idx                  = count.index
    tmKeys               = join(",", [for k in local.tmNamedKeyAllocation[count.index] : element(quorum_transaction_manager_keypair.tm.*.key_data, index(local.allTmNamedKeys, k))])
    nodeContainerIP      = cidrhost(local.network_cidr, count.index + 1 + 10)
    nodeContainerDNS     = format("node%d", count.index)
    txManagerContainerIP = cidrhost(local.network_cidr, count.index + 1 + 100)
    nodeVMPrivateIP      = aws_instance.node[count.index].private_ip
    nodeVMPublicIP       = aws_instance.node[count.index].public_ip
  }
}

resource "random_integer" "network_id" {
  max = 3000
  min = 1400
}

resource "quorum_bootstrap_network" "this" {
  name       = local.network_name
  target_dir = local.generated_dir
}

resource "quorum_bootstrap_keystore" "accountkeys-generator" {
  count        = local.number_of_nodes
 keystore_dir = format("%s/%s%s/keystore", quorum_bootstrap_network.this.network_dir_abs, local.node_dir_prefix, count.index)

  dynamic "account" {
    for_each = list("") // 1 account with empty passphrase
    content {
      passphrase = account.value
      balance    = account.key + 1
    }
  }
}

resource "quorum_bootstrap_node_key" "nodekeys-generator" {
  count = local.number_of_nodes
}

resource "quorum_transaction_manager_keypair" "tm" {
  count = length(local.allTmNamedKeys)
}

resource "local_file" "tm" {
  count    = length(local.allTmNamedKeys)
  filename = format("%s/%s", local.tmkeys_generated_dir, element(local.allTmNamedKeys, count.index))
  content  = quorum_transaction_manager_keypair.tm[count.index].key_data
}


locals {
  clique_extra_data = format("0x0000000000000000000000000000000000000000000000000000000000000000%s0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000", join("", [ for a in flatten(quorum_bootstrap_keystore.accountkeys-generator.*.account) : replace(a.address,"0x","") ]))
  accounts_alloc  = <<-EOT
  %{for a in flatten(quorum_bootstrap_keystore.accountkeys-generator.*.account)~}
    "${a.address}" : {
      "balance": "1000000000000000000000000000"
    },
  %{endfor~}
  EOT

  consensus_extra_data = {
    "ibft" = "${quorum_bootstrap_istanbul_extradata.this.extradata}"
    "clique" = "${local.clique_extra_data}"
    "raft"   = "0x0000000000000000000000000000000000000000000000000000000000000000"
  }

  consensus_genesis_config = {
    "ibft" = <<-EOT
     "istanbul": {
        "epoch": 30000,
        "policy": 0,
        "ceil2Nby3Block": 0
      },
    EOT
    "clique" = <<-EOT
     "clique": {
        "epoch": 30000,
        "period": ${var.blockPeriod}
      },
    EOT
    "raft"   = ""
  }

}

resource "local_file" "genesis-file" {
  filename = format("%s/genesis.json", quorum_bootstrap_network.this.network_dir_abs)
  content  = <<-EOF
{
    "alloc": {
      ${substr(local.accounts_alloc, 0, length(local.accounts_alloc) - 2)}
    },
    "coinbase": "0x0000000000000000000000000000000000000000",
    "config": {
      "homesteadBlock": 0,
      "byzantiumBlock": 0,
      "constantinopleBlock":0,
      "chainId": ${random_integer.network_id.result},
      "eip150Block": 0,
      "eip155Block": 0,
      "eip150Hash": "0x0000000000000000000000000000000000000000000000000000000000000000",
      "eip158Block": 0,
%{if var.is_quorum == true~}
      "isQuorum": true,
%{endif~}
      ${lookup(local.consensus_genesis_config, var.consensus, "" )}
      "maxCodeSize": 50
    },
    "difficulty": "${var.consensus == "ibft" ? "0x1" : "0x0"}",
    "extraData": "${lookup(local.consensus_extra_data, var.consensus, "0x0000000000000000000000000000000000000000000000000000000000000000")}",

    "gasLimit": "${var.gasLimit}",
    "mixhash": "${var.consensus == "ibft" ? data.quorum_bootstrap_genesis_mixhash.this.istanbul : "0x00000000000000000000000000000000000000647572616c65787365646c6578"}",
    "nonce": "0x0",
    "parentHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
    "timestamp": "0x00"
}
EOF
}

data "quorum_bootstrap_genesis_mixhash" "this" {

}

resource "quorum_bootstrap_istanbul_extradata" "this" {
  istanbul_addresses = quorum_bootstrap_node_key.nodekeys-generator[*].istanbul_address
}

resource "quorum_bootstrap_data_dir" "datadirs-generator" {
  count    = local.number_of_nodes
  data_dir = format("%s/%s%s", quorum_bootstrap_network.this.network_dir_abs, local.node_dir_prefix, count.index)
  genesis  = local_file.genesis-file.content
}

resource "local_file" "static-nodes" {
  count    = local.number_of_nodes
  filename = format("%s/static-nodes.json", quorum_bootstrap_data_dir.datadirs-generator[count.index].data_dir_abs)
  content  = "[${join(",", formatlist("\"enode://%s@%s:%d?discport=0&raftport=%d\"", [for m in data.null_data_source.meta[*].inputs : quorum_bootstrap_node_key.nodekeys-generator[m.idx].hex_node_id],  [for m in data.null_data_source.meta[*].inputs : m.nodeVMPrivateIP ] , local.host_p2p_port, local.host_raft_port))}]"
}

resource "local_file" "permissioned-nodes" {
  count    = local.number_of_nodes
  filename = format("%s/permissioned-nodes.json", quorum_bootstrap_data_dir.datadirs-generator[count.index].data_dir_abs)
  content  = local_file.static-nodes[count.index].content
}

resource "local_file" "passwords" {
  count    = local.number_of_nodes
  filename = format("%s/passwords.txt", quorum_bootstrap_data_dir.datadirs-generator[count.index].data_dir_abs)
  content  = ""
}

resource "local_file" "tmconfigs-generator" {
  count    = local.number_of_nodes
  filename = format("%s/%s%s/config.json", quorum_bootstrap_network.this.network_dir_abs, local.tm_dir_prefix, count.index)
  content  = <<-JSON
{
    "useWhiteList": false,
    "jdbc": {
        "username": "sa",
        "password": "",
        "url": "jdbc:h2:${local.tm_dir_container_path}/db;MODE=Oracle;TRACE_LEVEL_SYSTEM_OUT=0",
        "autoCreateTables": true
    },
    "serverConfigs":[
        {
            "app":"ThirdParty",
            "enabled": true,
            "serverAddress": "http://localhost:${local.container_tm_third_party_port}",
            "communicationType" : "REST"
        },
        {
            "app":"Q2T",
            "enabled": true,
             "serverAddress":"unix:${local.tm_dir_container_path}/tm.ipc",
            "communicationType" : "REST"
        },
        {
            "app":"P2P",
            "enabled": true,
            "serverAddress":"http://${aws_instance.node[count.index].private_ip}:${local.container_tm_p2p_port}",
            "sslConfig": {
              "tls": "OFF",
              "generateKeyStoreIfNotExisted": true,
              "serverKeyStore": "${local.tm_dir_container_path}/server-keystore",
              "serverKeyStorePassword": "quorum",
              "serverTrustStore": "${local.tm_dir_container_path}/server-truststore",
              "serverTrustStorePassword": "quorum",
              "serverTrustMode": "TOFU",
              "knownClientsFile": "${local.tm_dir_container_path}/knownClients",
              "clientKeyStore": "${local.tm_dir_container_path}/client-keystore",
              "clientKeyStorePassword": "quorum",
              "clientTrustStore": "${local.tm_dir_container_path}/client-truststore",
              "clientTrustStorePassword": "quorum",
              "clientTrustMode": "TOFU",
              "knownServersFile": "${local.tm_dir_container_path}/knownServers"
            },
            "communicationType" : "REST"
        }
    ],
    "peer": [${join(",", data.null_data_source.meta[*].inputs.peerJson)}],
    "keys": {
      "passwords": [],
      "keyData": [${data.null_data_source.meta[count.index].inputs.tmKeys}]
    },
    "alwaysSendTo": []
}
JSON
}
