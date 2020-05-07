# TPS Monitor
Tool to monitor transactions per second in Quorum.
The tool can be used to calculate TPS in one of the two ways given below:
1. calculate TPS on new blocks inserted to chain - it continuosly calculates TPS on new blocks inserted to the chain. 
2. calculate TPS for a given range of blocks - it calculates TPS for a given range of blocks.

##### In both modes it exposes the TPS data (aggregated at minute level) via http endpoint.

## Building the source

Building tpsmonitor requires Go (version 1.13 or later). 
You can install it using your favourite package manager. Once the dependencies are installed, run
`make`

## Running `tpsmonitor`

#### Usage
Run `tpsmonitor --help` to see usage.

WS API must be enabled in `geth` to run tpsmonitor. Default port of HTTP endpoint is `7575`

#### calculate TPS on new blocks
Displays TPS calculated in the console for new blocks as they are inserted to the chain. Calculates TPS, total no of blocks and total no of transactions for every second and saves these results to the report file.

```tpsmonitor --wsendpoint <ws address> --consensus [raft|ibft] --report <report name>```
Example: `tpsmonitor --wsendpoint ws://52.77.226.85:23000/ --consensus raft --report tps-m.csv --port <port no>`

Sample report:

```aidl
head -20 tps-m.csv
```
````
localTime,refTime,TPS,TxnCount,BlockCount
Mar-23 06:03:01,00:00:00:01,2134,128047,241
Mar-23 06:03:02,00:00:00:02,2023,242871,479
Mar-23 06:03:03,00:00:00:03,1958,352613,713
Mar-23 06:03:04,00:00:00:04,1906,457647,937
Mar-23 06:03:05,00:00:00:05,1880,564064,1163
````
#### calculate TPS for a given range of blocks
Displays TPS calculated in the console for given range of blocks as they are read from the chain. Calculates TPS, total no of blocks and total no of transactions for every minute(for the given block range) and saves these results to the report file.

```tpsmonitor --wsendpoint <ws address> --consensus [raft|ibft] --report <report name>```
Example: `tpsmonitor --wsendpoint ws://52.77.226.85:23000/ --consensus raft --from 1 --to 10000 --report tps-m.csv --port 8888`

#### HTTP endpoint for TPS data

`http://<host>:<port>/tpsdata`

sample output:
```aidl
http://localhost:8888/tpsdata
```

````
localTime,refTime,TPS,TxnCount,BlockCount
Mar-23 06:03:01,00:00:00:01,2134,128047,241
Mar-23 06:03:02,00:00:00:02,2023,242871,479
Mar-23 06:03:03,00:00:00:03,1958,352613,713
Mar-23 06:03:04,00:00:00:04,1906,457647,937
Mar-23 06:03:05,00:00:00:05,1880,564064,1163
````
