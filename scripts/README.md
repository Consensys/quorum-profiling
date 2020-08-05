## Test Quorum locally
 In this scenario it is assumed that you have a quorum network running locally and want to run Jmeter stress test on it.
 
 ### Usage
 Update `quorum-test/scripts/jmeter/host_acct.csv` and `quorum-test/scripts/jmeter/network.properties` with details of your quorum network and Jmeter test parameters.
 
 To choose correct Jmeter test profile refer to section _Test profiles_ [here](../stresstest-aws/README.md#test-profiles)
 
 #### To Start Test
`cd quorum-test/scripts`

 `./start-test.sh --testProfile <jmeter-test-profile> --consensus <ibft|raft> --endpoint <quorum-rpc-endpoint> --basedir <repo base dir>`
 
 example: `./start-test.sh --testProfile "4node/deploy-contract-public" --consensus "ibft" --endpoint "http://host.docker.internal:22000" --basedir ~/go/src/github.com/QuorumEngineering/quorum-test`
 
 This brings up `influxdb`, `grafana`, `telegraf`, Jmeter test` and `tps-monitor` containers. 
 
 ##### Grafana dashboard 
  It can be accessed at `http://localhost:3000/login`. Enter `admin/admin` and access predefined dashboards `Quorum Profiling Dashboard` & `Quorum Profiling Jmeter Dashboard`
  
 ##### Influxdb 
  It can be access at `http://localhost:8086/`. DB: `telegraf` user/password: `telegraf/test123`
 
 ##### Prometheus metrics  
  * Quorum node cpu/memory usage metrics can be accessed at `http://localhost:9126/metrics`.
  * TPS metrics can be accessed at `http://localhost:2112/metrics`.
 
 #### To Stop Test
 
 `cd quorum-test/scripts`
 
 grep for `jmeter` docker container and stop it.
 
 grep for  `tpsmonitor` docker container and stop it.
 
 run `docker-compose down` . It will stop `grafana`, `telegraf` and `influxdb`
  
 
   ![Quorum Dashboard](../stresstest-aws/quorumDashboard.jpeg) 
   
   ![Jmeter Dashboard](../stresstest-aws/jmeterDashboard.jpeg) 
   