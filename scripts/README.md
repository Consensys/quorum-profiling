# Test Quorum locally
 In this scenario it is assumed that there is a running Quorum network and the goal is to run Jmeter stress test on it.
 
 To choose correct Jmeter test profile refer to section _Test profiles_ [here](../jmeter-test)
 
## Prerequisites to run the test
  Refer to **scenario 4** [here](../README.md#prerequisites-for-test-execution) for all prerequisites.
 
## To start test
`cd quorum-profiling/scripts`

 `./start-test.sh --testProfile <jmeter-test-profile> --consensus <ibft|raft> --endpoint <quorum-rpc-endpoint> --basedir <repo base dir>`
 
 example: `./start-test.sh --testProfile "4node/deploy-contract-public" --consensus "ibft" --endpoint "http://host.docker.internal:22000" --basedir ~/go/src/github.com/jpmorganchase/quorum-profiling`
 
 This brings up `influxdb`, `grafana`, `telegraf`, `Jmeter test` and `tps-monitor` containers. 
 
## Grafana dashboard 
  It can be accessed at `http://localhost:3000/login`. Enter `admin/admin` as user id and password to access the predefined dashboards `Quorum Profiling Dashboard` & `Quorum Profiling Jmeter Dashboard`. Sample dashboard are shown below.
 
## Influxdb 
  It can be access at `http://localhost:8086/`. The database name is `telegraf` and user/password is `telegraf/test123`
  > if you wish to change the port, default user id/password, please edit the [telegraf.conf](telegraf/telegraf.conf) file
 
 
 **Note!!!** The endpoint for influxdb instance can be configured by setting up the `influxdburl` in the properties file as shown below.
   
   ```
  #to write jmeter test summary to influxdb
  influxdburl=http://host.docker.internal:8086/write?db=telegraf
  ```

## Prometheus metrics  
  * Quorum node cpu/memory usage metrics can be accessed at `http://localhost:9126/metrics`.
  * TPS metrics can be accessed at `http://localhost:2112/metrics`.
  > if you wish to change the port for `prometheus`, please edit the [telegraf.conf](telegraf/telegraf.conf) file
 
## To Stop Test
 
> `cd quorum-profiling/scripts`

> grep for `jmeter` docker container and stop it.
 
> grep for  `tpsmonitor` docker container and stop it.
 
> run `docker-compose down` . It will stop `grafana`, `telegraf` and `influxdb`
     
  
   