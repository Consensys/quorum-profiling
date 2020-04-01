 ## Usage

 * `deploy-contract-public-4node.jmx` fire transactions to deploy new contracts to all 4 nodes in the network in parallel via `json-rpc` requests. The contract being deployed is a `SimpleStorage` contract, with a random value being initialised in the constructor.  
   
   Parameters that can be specified include  
     * `threads` : Number of threads per node used to send transactions (i.e. threads=1 means 4 threads to send to 4 nodes)  
     * `seconds` : Duration of the test run in seconds  
     * `delay` : Startup delay (default to 5 seconds)  
       
   and
   
     * `url$i` : The RPC endpoint url  
     * `port$i` : The RPC endpoint port  
     * `from$i` : Geth account used to send transactions from  
     
   with `$i` being value 1..4 corresponding to node1 to node4 of the network
        
   Sample usage
    ```shell script
    jmeter -n -t deploy-contract-public-4node.jmx 
         -Jurl1=localhost -Jport1=22000 -Jfrom1=0xed9d02e382b34818e88b88a309c7fe71e65f419d 
         -Jurl1=localhost -Jport2=22001 -Jfrom2=0x4204266650c946a56da82dfded6029cd8b1b54cf
         -Jurl2=localhost -Jport3=22002 -Jfrom3=0xca40127ac0880f44bca898fd357557b70a2fcc42
         -Jurl2=localhost -Jport4=22003 -Jfrom4=0x53a52871988c3b3856280181105d0541d78b38ac
         -Jthreads=1 -Jseconds=60
    ```

        