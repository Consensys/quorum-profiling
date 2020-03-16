 ## Introduction
 `start-top-mon.sh` monitors CPU and Memory usage of `geth` and `tessera` processes periodically (every 60seconds) by filtering the output of `top` command and writes the result to a log file with timestamp.
 This should be run locally where the `geth` and `tessera` processes are running in a node.
 The log file can be used to create graph at the end of the test.
 Sample output from log as follows:
 ```
 #date time processid memusage cpu% mem%
21/02/2020 01:00:27 13086 4.1g 100.0 13.2
21/02/2020 01:01:28 13086 4.1g 100.0 13.2
21/02/2020 01:02:28 13086 4.2g 200.0 13.6
21/02/2020 01:03:28 13086 4.3g 100.0 13.9
```
 