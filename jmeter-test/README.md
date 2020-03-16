 ## Introduction
 You need to install `jmeter` to run the following tests:
 * `public-contract-creation-4node.jmx` contains test case to deploy simple storage smart contract as a public contract. It runs 4 threads one for each node in a 4node network.
 Update it with correct path to `public_accts_host.csv`. If you want to enable delay after every request you can enable the constant timer.
 Update `public_accts_host.csv` with your environment details.
 
* `private-contract-creation-4node.jmx` contains test case to deploy simple storage smart contract as a private contract( node1 to node2, node2 to node3, node3 to node4 and node4 to node1 ). It runs 4 threads one for each node in a 4node network.
 Update it with correct path to `private_accts_host.csv`. If you want to enable delay after every request you can enable the constant timer.
 Update `private_accts_host.csv` with your environment details.
 