 ## Introduction
 This java program is a testing tool to test sending value transfer transactions.
 This program does the following:
 * Reads the inputs from a input file
 * For each input (a line in the input file) create a separate thread
 * Each thread sends given number of signed value transfer (`1` ether) transactions in sequence. Make sure there is enough balance in the `fromAccount`
 * All the threads run concurrently
 * if a transaction fails it prints the error and that transactions is skipped
 * Once all the threads finish running it prints out the total time taken to process and the total number
 of transactions sent by each thread.
 
 ## Format of input file
 * each line should have comma separated unquoted strings in the following format:
 `endpoint,fromAccount,toAccount,nonce,numberOfTransactions,delayAfter,delay`
 * `endpoint` the RPC endpoint of ethereum client
 * `fromAccount` the from account from which ether will be transferred. It should be unique across all inputs in the file.
 * `toAccount` the to account recieving the the ether transfer
 * `nonce` the current value of nonce of the from account
 * `numberOfTransactions` number of transactions to be executed
 * `delayAfter` the number of transactions after which the thread should be paused. set it to `0` if no pause is required
 * `delay` the duration in milliseconds to be paused after every `delayAfter` number of transactions
 * sample input: 
 ```
 http://localhost:22001,0xca843569e3427144cead5e4d5999a3d0ccf92b8e,0xd8f63ab1bd6057933a177a7ea5809ee3d4a7f5a6,40003,1,0,1000
 http://localhost:22000,0xed9d02e382b34818e88b88a309c7fe71e65f419d,0xd8f63ab1bd6057933a177a7ea5809ee3d4a7f5a6,64004,1,0,1000
```
 ## how to build and run?
 * build `mvn clean compile assembly:single`
 * This would create  `target/quorum-test-1.0-SNAPSHOT-jar-with-dependencies.jar`
 * run `java -jar quorum-test-1.0-SNAPSHOT-jar-with-dependencies.jar <inputfile>`
 
## TEST AUTOMATION - TODO Items


![test flowchart](QuorumTestFlowchart.jpg) 


![test architecture](QuorumTestArch.png) 

## List of tasks

#### 1. Jmeter - create test profile
Create test plan in Jmeter that can run different transaction types in a loop i.e, in each iteration a combination of different transaction types should be sent to geth client. The following transaction types should be covered.
* Public contract creation
* Public contract - state change
* Private contract creation
* Private contract - state change
* value transfer

#### 2. Terraform - spin up 4 node Raft/istanbul network
* Create Terraform code to spin up 4 node network with Raft/istanbul consensus. 
Each node should be EC2 instance `t3a.2xlarge` with at least 250G storage.
Each node should have a `geth` and `tessera` instance. 
terraform should generate a config file with public IP of all the 4 nodes should be generated and used as input to the test node to run the test plans. 


* Create Terraform code to spin up a node to run the Jmeter test plan. The node should be EC2 instance `t3a.2xlarge` with at least 50G storage.
The node should have Jmeter and chainhammer available. 

#### 3. Chart generation tool
This tool should generate the following charts once Jmeter test plan has finished running.
* TPS over time
* Total transactions over time
* Total blocks over time
* CPU and memory usage of geth and java processes

The tool should transform time input to start from `00:00` and allow flexibility to aggregate statistics at second/minute level.
If data from two test environments are provided the tool should generate single chart with data from one environment overlaying on the other.