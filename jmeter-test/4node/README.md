# Usage

 * `deploy-contract-public-4node.jmx` fire transactions to deploy new contracts to all 4 nodes in the network in parallel via `json-rpc` requests. The contract being deployed is a `SimpleStorage` contract, with a random value being initialised in the constructor.  
   
   Parameters that can be specified include  
     * `threads` : Number of threads per node used to send transactions (i.e. threads=1 means 4 threads to send to 4 nodes)  
     * `seconds` : Duration of the test run in seconds  
     * `delay` : Startup delay (default to 5 seconds)  
     * `throughput` : (value is transactions **per minute**). Only specify this if you want to throttle the throughput so transactions can be sent slower 
     * `influxdburl` : influxdb url for capturing metrics. Refer [here](../README.md#Disabling-influxDB) on how to disable influxdb. 
   
   and
   
     * `url$i` : The RPC endpoint url  
     * `port$i` : The RPC endpoint port  
     * `from$i` : Geth account used to send transactions from  
     
   with `$i` being value `1`..`4` corresponding to node1 to node4 of the network
        
   Sample usage
    ```shell script
   #!/bin/bash
   #without influxdb
    jmeter -n -t deploy-contract-public-4node.jmx 
         -Jurl1=localhost -Jport1=22000 -Jfrom1=0xed9d02e382b34818e88b88a309c7fe71e65f419d 
         -Jurl2=localhost -Jport2=22001 -Jfrom2=0x4204266650c946a56da82dfded6029cd8b1b54cf
         -Jurl3=localhost -Jport3=22002 -Jfrom3=0xca40127ac0880f44bca898fd357557b70a2fcc42
         -Jurl4=localhost -Jport4=22003 -Jfrom4=0x53a52871988c3b3856280181105d0541d78b38ac
         -Jthreads=1 -Jseconds=60
   
   #with influxdb
       jmeter -n -t deploy-contract-public-4node.jmx 
            -Jinfluxdburl=http://localhost:8086/write?db=telegraf
            -Jurl1=localhost -Jport1=22000 -Jfrom1=0xed9d02e382b34818e88b88a309c7fe71e65f419d 
            -Jurl2=localhost -Jport2=22001 -Jfrom2=0x4204266650c946a56da82dfded6029cd8b1b54cf
            -Jurl3=localhost -Jport3=22002 -Jfrom3=0xca40127ac0880f44bca898fd357557b70a2fcc42
            -Jurl4=localhost -Jport4=22003 -Jfrom4=0x53a52871988c3b3856280181105d0541d78b38ac
            -Jthreads=1 -Jseconds=60
    ```
   
   
 * `deploy-contract-private-4node.jmx` similarly fire transactions to deploy new private `SimpleStorage` contracts to all 4 nodes in the network in parallel with a `privateFor` recipient. The contracts are also being initialised with random value.  
 
    Parameters are similar to the above, with an extra variable  
      * `privateFor$i` : Public key of the private recipient
      
    with `$i` being value `1`..`4` corresponding to node1 to node4 of the network  
    
    Sample usage
    ```shell script
     #!/bin/bash
     #without influxdb
     jmeter -n -t deploy-contract-private-4node.jmx 
          -Jurl1=localhost -Jport1=22000 -Jfrom1=0xed9d02e382b34818e88b88a309c7fe71e65f419d -JprivateFor1=\"BULeR8JyUWhiuuCMU/HLA0Q5pzkYT+cHII3ZKBey3Bo=\"
          -Jurl2=localhost -Jport2=22001 -Jfrom2=0x4204266650c946a56da82dfded6029cd8b1b54cf -JprivateFor2=\"QfeDAys9MPDs2XHExtc84jKGHxZg/aj52DTh0vtA3Xc=\"
          -Jurl3=localhost -Jport3=22002 -Jfrom3=0xca40127ac0880f44bca898fd357557b70a2fcc42 -JprivateFor3=\"1iTZde/ndBHvzhcl7V68x44Vx7pl8nwx9LqnM/AfJUg=\"
          -Jurl4=localhost -Jport4=22003 -Jfrom4=0x53a52871988c3b3856280181105d0541d78b38ac -JprivateFor4=\"oNspPPgszVUFw0qmGFfWwh1uxVUXgvBxleXORHj07g8=\"
          -Jthreads=1 -Jseconds=60
   
     #with influxdb
     jmeter -n -t deploy-contract-private-4node.jmx 
          -Jinfluxdburl=http://localhost:8086/write?db=telegraf
          -Jurl1=localhost -Jport1=22000 -Jfrom1=0xed9d02e382b34818e88b88a309c7fe71e65f419d -JprivateFor1=\"BULeR8JyUWhiuuCMU/HLA0Q5pzkYT+cHII3ZKBey3Bo=\"
          -Jurl2=localhost -Jport2=22001 -Jfrom2=0x4204266650c946a56da82dfded6029cd8b1b54cf -JprivateFor2=\"QfeDAys9MPDs2XHExtc84jKGHxZg/aj52DTh0vtA3Xc=\"
          -Jurl3=localhost -Jport3=22002 -Jfrom3=0xca40127ac0880f44bca898fd357557b70a2fcc42 -JprivateFor3=\"1iTZde/ndBHvzhcl7V68x44Vx7pl8nwx9LqnM/AfJUg=\"
          -Jurl4=localhost -Jport4=22003 -Jfrom4=0x53a52871988c3b3856280181105d0541d78b38ac -JprivateFor4=\"oNspPPgszVUFw0qmGFfWwh1uxVUXgvBxleXORHj07g8=\"
          -Jthreads=1 -Jseconds=60
      ```


  
  
**Note**: `threads$i` parameter can also be used to customized the number of threads used for a spefic node - or to disable sending transactions for that particular node

Example
 ```shell script
    #!/bin/bash
    #without influxdb
    jmeter -n -t deploy-contract-public-4node.jmx 
    -Jurl1=localhost -Jport1=22000 -Jfrom1=0xed9d02e382b34818e88b88a309c7fe71e65f419d 
    -Jurl2=localhost -Jport2=22001 -Jfrom2=0x4204266650c946a56da82dfded6029cd8b1b54cf
    -Jthreads=1 -Jseconds=60
    -Jthreads3=0 -Jthreads4=0

    #with influxdb
    jmeter -n -t deploy-contract-public-4node.jmx 
    -Jinfluxdburl=http://localhost:8086/write?db=telegraf    
    -Jurl1=localhost -Jport1=22000 -Jfrom1=0xed9d02e382b34818e88b88a309c7fe71e65f419d 
    -Jurl2=localhost -Jport2=22001 -Jfrom2=0x4204266650c946a56da82dfded6029cd8b1b54cf
    -Jthreads=1 -Jseconds=60
    -Jthreads3=0 -Jthreads4=0
```
The command above will start 2 threads to send transactions to node1 and node2 only  


As there are many properties required, these can also be put into a `properties` file and given to the test execution via command line.  

Example property file `vars.properties`

```shell script
#Node 1
from1=0x4204266650c946a56da82dfded6029cd8b1b54cf
url1=3.10.116.15
port1=22000

#Node2
from2=0xca40127ac0880f44bca898fd357557b70a2fcc42
url2=35.176.231.84
port2=22001

#Node3
from3=0x53a52871988c3b3856280181105d0541d78b38ac
url3=3.9.17.190
port3=22002

#Node4
from4=0xa13ced78febcac20a8268b796a9e0208c17d8313
url4=35.178.199.42
port4=22003


threads=1
seconds=10
delay=5

influxdburl=http://localhost:8086/write?db=telegraf

```

Test will then be executed using the following command  

```shell script
jmeter -n -t deploy-contract-public-4node.jmx -q vars.properties
```

Note that `-q` option was used to specify additional property file, instead of `-p` which will replace the whole existing default `jmeter.properties` file