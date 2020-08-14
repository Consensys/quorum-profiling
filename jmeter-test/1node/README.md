# Usage

 * `deploy-contract-public.jmx` test plan fire transactions to deploy new contracts to the network via `json-rpc` requests. The contract being deployed is a `SimpleStorage` contract, with a random value being initialised in the constructor.  
   
   Parameters that can be specified include  
     * `url` : The RPC endpoint url
     * `influxdburl` : influxdb url for capturing metrics. Refer [here](../README.md#Disabling-influxDB) on how to disable influxdb.
     * `port` : The RPC endpoint port  
     * `from` : Geth account used to send transactions from  
     * `threads` : Number of threads use to send transactions  
     * `seconds` : Duration of the test run in seconds  
     * `delay` : Startup delay (default to 5 seconds)  
     * `throughput` : (value is transactions **per minute**). Only specify this if you want to throttle the throughput so transactions can be sent slower
        
   Sample usage
    ```shell script
   #!/bin/bash
   #without influxdb
    jmeter -n -t deploy-contract-public.jmx -Jurl=localhost -Jport=22000 -Jfrom=0xed9d02e382b34818e88b88a309c7fe71e65f419d -Jthreads=10 -Jseconds=60
   
   #with influxdb
   jmeter -n -t deploy-contract-public.jmx -Jinfluxdburl=http://localhost:8086/write?db=telegraf -Jurl=localhost -Jport=22000 -Jfrom=0xed9d02e382b34818e88b88a309c7fe71e65f419d -Jthreads=10 -Jseconds=60
   ```

  
 * `deploy-contract-private.jmx` test plan similarly fire transactions to deploy new private `SimpleStorage` contracts to the network with a `privateFor` recipient. The contracts are also being initialised with random value.
 
    Parameters are similar to the above, with an extra variable  
      * `privateFor` : Public key of the private recipient
      
    Sample usage
    ```shell script
   #!/bin/bash
   #without influxdb
    jmeter -n -t deploy-contract-private.jmx -Jurl=localhost -Jport=22000 -Jfrom=0xed9d02e382b34818e88b88a309c7fe71e65f419d -JprivateFor=\"ROAZBWtSacxXQrOe3FGAqJDyJjFePR5ce4TSIzmJ0Bc=\" -Jthreads=10 -Jseconds=60
   
   #with influxdb
   jmeter -n -t deploy-contract-private.jmx -Jinfluxdburl=http://localhost:8086/write?db=telegraf -Jurl=localhost -Jport=22000 -Jfrom=0xed9d02e382b34818e88b88a309c7fe71e65f419d -JprivateFor=\"ROAZBWtSacxXQrOe3FGAqJDyJjFePR5ce4TSIzmJ0Bc=\" -Jthreads=10 -Jseconds=60
    ```
   
   
 * `update-contract-public.jmx` fire transactions to update the state of an existing contract in the network. The test plan will carry out the following step
    1. Deploy a new public contract to the network
    2. Query the contract address after the contract is mined by calling `eth_getTransactionReceipt` with transaction hash retrieved from response of (i)
    3. Extract the contract address from response of (ii) and send multiple transactions to the contract deployed with a different randomly generated value.
    
    Sample usage
    ```shell script
   #!/bin/bash
   #without influxdb
    jmeter -n -t update-contract-public.jmx -Jurl=localhost -Jport=22000 -Jfrom=0xed9d02e382b34818e88b88a309c7fe71e65f419d -Jthreads=10 -Jseconds=60
   
   #with influxdb
    jmeter -n -t update-contract-public.jmx -Jinfluxdburl=http://localhost:8086/write?db=telegraf -Jurl=localhost -Jport=22000 -Jfrom=0xed9d02e382b34818e88b88a309c7fe71e65f419d -Jthreads=10 -Jseconds=60
    ```
   
   
 * `update-contract-private.jmx` fire transactions to update the state of an existing private contract. The test plan will carry out the follow steps:  
     1. Deploy a new private contract to the network (with the specified `privateFor` key)
     2. Query the contract address after the contract is mined by calling `eth_getTransactionReceipt` with transaction hash retrieved from response of (i)
     3. Extract the contract address from response of (ii) and send multiple transactions to the contract deployed with a different randomly generated value.
     
     Sample usage
     ```shell script
   #!/bin/bash
   #without influxdb
     jmeter -n -t update-contract-private.jmx -Jurl=localhost -Jport=22000 -Jfrom=0xed9d02e382b34818e88b88a309c7fe71e65f419d -JprivateFor=\"ROAZBWtSacxXQrOe3FGAqJDyJjFePR5ce4TSIzmJ0Bc=\" -Jthreads=10 -Jseconds=60
   
   #with influxdb
     jmeter -n -t update-contract-private.jmx -Jinfluxdburl=http://localhost:8086/write?db=telegraf -Jurl=localhost -Jport=22000 -Jfrom=0xed9d02e382b34818e88b88a309c7fe71e65f419d -JprivateFor=\"ROAZBWtSacxXQrOe3FGAqJDyJjFePR5ce4TSIzmJ0Bc=\" -Jthreads=10 -Jseconds=60
     ```
        