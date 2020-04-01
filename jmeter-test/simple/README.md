 ## Usage

 * `deploy-contract-public.jmx` test plan fire transactions to deploy new contracts to the network via `json-rpc` requests. Parameters that can be specified include  
     * `url` : The RPC endpoint url  
     * `port` : The RPC endpoint port  
     * `from` : Geth account used to send transactions from  
     * `threads` : Number of threads use to send transactions  
     * `seconds` : Duration of the test run in seconds  
     * `delay` : Startup delay (default to 5 seconds)  
        
    Sample usage
    ```shell script
    jmeter -n -t deploy-contract-public.jmx -Jurl=localhost -Jport=22000 -Jfrom=0xed9d02e382b34818e88b88a309c7fe71e65f419d -Jthreads=10 -Jseconds=60
    ```

 
 