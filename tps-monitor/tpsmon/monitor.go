package tpsmon

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// TPSRecord represents a data point of TPS at a specific time
type TPSRecord struct {
	rtime string // reference time starts from time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)
	ltime string // block time in local time
	tps   uint32 // no of transactions per second
	blks  uint64 // total block count
	txns  uint64 // total transaction count
}

func (t TPSRecord) String() string {
	return fmt.Sprintf("TPSRecord: ltime:%v rtime:%v tps:%v txns:%v blks:%v", t.ltime, t.rtime, t.tps, t.txns, t.blks)
}

func (t TPSRecord) ReportString() string {
	return fmt.Sprintf("%v,%v,%v,%v,%v\n", t.ltime, t.rtime, t.tps, t.txns, t.blks)
}

// TPSMonitor implements a monitor service
type TPSMonitor struct {
	isRaft       bool                  // represents consensus
	sub          ethereum.Subscription // ethereum subscriptions
	headers      chan *types.Header    // block header from chain
	client       *ethclient.Client     // ethereum client
	tpsRecs      []TPSRecord           // list of TPS data points recorded
	report       string                // report name to store TPS data points
	fromBlk      uint64                // from block number
	toBlk        uint64                // to block number
	stopc        chan struct{}         // stop channel
	firstBlkTime *time.Time            // first block's time
	refTime      time.Time             // reference time
	refTimeNext  time.Time             // next expected reference time
	blkTimeNext  time.Time             // next expected block time
	blkCnt       uint64
	txnsCnt      uint64   // total transaction count
	rptFile      *os.File // report file
	awsCfg       *AwsCloudwatchService
}

// Date format to show only hour and minute
const (
	dateFmtMin = "Jan-02 15:04"
)

func NewTPSMonitor(awsCfg *AwsCloudwatchService, isRaft bool, report string, frmBlk uint64, toBlk uint64, sub ethereum.Subscription, headers chan *types.Header, client *ethclient.Client) *TPSMonitor {

	tm := &TPSMonitor{
		isRaft:  isRaft,
		report:  report,
		sub:     sub,
		headers: headers,
		client:  client,
		fromBlk: frmBlk,
		toBlk:   toBlk,
		stopc:   make(chan struct{}),
		awsCfg:  awsCfg,
	}

	if tm.report != "" {
		var err error
		if tm.rptFile, err = os.Create(tm.report); err != nil {
			log.Fatalf("error creating report file %s\n", tm.report)
		}
		if _, err := tm.rptFile.WriteString("localTime,refTime,TPS,TxnCount,BlockCount\n"); err != nil {
			log.Printf("ERROR: writing to report failed err:%v", err)
		}
		tm.rptFile.Sync()
	}
	return tm
}

// starts service to calculate tps
func (tm *TPSMonitor) Start() {
	tm.init()
	if tm.fromBlk > 0 && tm.toBlk > 0 {
		go tm.calcTpsFromBlockRange()
		log.Printf("tps calc started - fromBlock:%v toBlock:%v\n", tm.fromBlk, tm.toBlk)
	} else {
		go tm.calcTpsFromNewBlocks()
		log.Print("tps monitor started")
	}
}

// stops service
func (tm *TPSMonitor) Stop() {
	close(tm.stopc)
	if tm.rptFile != nil {
		tm.rptFile.Close()
	}
}

// waits for stop signal to end service
func (tm *TPSMonitor) Wait() {
	log.Print("tps monitor waiting to stop")
	<-tm.stopc
	log.Print("tps monitor wait over stopping")
}

// initializes service
func (tm *TPSMonitor) init() {
	if tm.isRaft {
		log.Print("consensus is raft")
	} else {
		log.Print("consensus is ibft")
	}
	tm.refTime = time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)
	tm.refTimeNext = time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)
	tm.blkCnt = 0
	tm.txnsCnt = 0
}

// read data from block and calculated TPS
func (tm *TPSMonitor) readBlock(block *types.Block) {

	var blkTime time.Time
	if tm.isRaft {
		r := block.Time() % 1e9
		blkTime = time.Unix(int64(block.Time()/1e9), int64(r))
	} else {
		blkTime = time.Unix(int64(block.Time()), 0)
	}

	if tm.firstBlkTime != nil {
		totSecs := blkTime.Sub(*tm.firstBlkTime).Milliseconds() / 1000
		if totSecs > 0 {
			tps := tm.txnsCnt / uint64(totSecs)
			log.Printf("TPS:%v txnsCnt:%v blkCnt:%v\n", tps, tm.txnsCnt, tm.blkCnt)
		}
	}

	if tm.firstBlkTime == nil {
		tm.firstBlkTime = &blkTime
		tm.refTimeNext = tm.refTimeNext.Add(time.Minute)
		tm.blkTimeNext = blkTime.Add(time.Minute)
	}

	txns := len(block.Transactions())

	if blkTime.After(tm.blkTimeNext) {

		// this loop fills tps for missing time points when block mining is delayed beyond 1minute
		for blkTime.After(tm.blkTimeNext) {
			ltime := tm.blkTimeNext.Format(dateFmtMin)
			yd := tm.refTimeNext.YearDay() - 1
			hh := tm.refTimeNext.Hour()
			mm := tm.refTimeNext.Minute()
			rtime := fmt.Sprintf("%02d:%02d:%02d", yd, hh, mm)
			totSecs := tm.refTimeNext.Sub(tm.refTime).Milliseconds() / 1000
			tps := tm.txnsCnt / uint64(totSecs)
			tr := TPSRecord{rtime: rtime, ltime: ltime, tps: uint32(tps), blks: tm.blkCnt, txns: tm.txnsCnt}
			log.Print(tr.String() + "\n")
			if tm.rptFile != nil {
				if _, err := tm.rptFile.WriteString(tr.ReportString()); err != nil {
					log.Printf("ERROR: writing to report failed %v", err)
				}
				tm.rptFile.Sync()
			}
			tm.tpsRecs = append(tm.tpsRecs, tr)
			//publish metrics to aws cloudwatch
			go tm.putMetricsInAws(tm.blkTimeNext, fmt.Sprintf("%v", tps), fmt.Sprintf("%v", tm.txnsCnt), fmt.Sprintf("%v", tm.blkCnt))
			tm.refTimeNext = tm.refTimeNext.Add(time.Minute)
			tm.blkTimeNext = tm.blkTimeNext.Add(time.Minute)
		}
	}

	tm.blkCnt++
	tm.txnsCnt += uint64(txns)
}

func (tm *TPSMonitor) putMetricsInAws(lt time.Time, tps string, txnCnt string, blkCnt string) {
	if tm.awsCfg != nil {
		tm.awsCfg.PutMetrics("TPS", tps, lt)
		tm.awsCfg.PutMetrics("TxnCount", txnCnt, lt)
		tm.awsCfg.PutMetrics("BlockCount", blkCnt, lt)
	}
}

// calculates TPS for new block added to the chain
func (tm *TPSMonitor) calcTpsFromNewBlocks() {
	for {
		select {
		case err := <-tm.sub.Err():
			log.Fatal(err)
		case header := <-tm.headers:
			block, err := tm.client.BlockByNumber(context.Background(), header.Number)
			if err != nil {
				log.Fatal(err)
			}
			tm.readBlock(block)

		case <-tm.stopc:
			log.Print("tps monitor stopped - exit loop")
			return
		}
	}
}

// calculates TPS for a given block range
func (tm *TPSMonitor) calcTpsFromBlockRange() {
	stBlk := tm.fromBlk
	toBlk := tm.toBlk
	for stBlk <= toBlk {
		block, err := tm.client.BlockByNumber(context.Background(), big.NewInt(int64(stBlk)))
		if err != nil {
			log.Fatal(err)
		}
		tm.readBlock(block)
		stBlk++
	}
}

func (tm *TPSMonitor) printTPS() {
	trl := len(tm.tpsRecs)
	log.Printf("Total tps records %d\n", trl)
	for i, v := range tm.tpsRecs {
		log.Printf("%d. %v\n", i, v.String())
	}
}
