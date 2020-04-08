package tpsmon

import (
	"bufio"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
)

func getNewTxn() *types.Transaction {
	emptyTx := types.NewTransaction(
		0,
		common.HexToAddress("095e7baea6a6c7c4c2dfeb977efac326af552d87"),
		big.NewInt(0), 0, big.NewInt(0),
		nil,
	)
	return emptyTx
}

func getNewHeader(n int64, isRaft bool, t1 time.Time) *types.Header {
	fakeParentHash := common.HexToHash("0xc2c1dc1be8054808c69e06137429899d")

	var tm int64
	if isRaft {
		tm = t1.UnixNano()
	} else {
		tm = t1.Unix()
	}
	header := &types.Header{
		ParentHash: fakeParentHash,
		Number:     big.NewInt(n),
		Difficulty: big.NewInt(1),
		GasLimit:   uint64(0),
		GasUsed:    uint64(0),
		Coinbase:   common.HexToAddress("095e7baea6a6c7c4c2dfeb977efac326af552d87"),
		Time:       uint64(tm),
	}
	return header
}

func txnsList(n int) []*types.Transaction {
	i := 0
	var txns []*types.Transaction
	for i < n {
		txns = append(txns, getNewTxn())
		i++
	}
	return txns
}

func getNewBlock(h *types.Header, txns []*types.Transaction) *types.Block {
	return types.NewBlock(h, txns, nil, nil)
}

func testTPS(isRaft bool, t *testing.T) {
	assert := assert.New(t)
	var tempFile string
	if f, err := ioutil.TempFile("", "_tps_test"); err != nil {
		t.Fatalf("error creating temp file for test")
	} else {
		tempFile = f.Name()
	}

	defer os.Remove(tempFile)

	tm := NewTPSMonitor(nil, isRaft, tempFile, 1, 10, nil, nil, nil)

	assert.NotNil(tm, "creating tps monitor failed")

	tm.init()

	var blksArr []*types.Block
	var c int
	var t1, t2 time.Time
	if isRaft {
		unsec := time.Now().UnixNano()
		nsec := unsec % 1e9
		t1 = time.Unix(unsec/1e9, nsec)
	} else {
		t1 = time.Unix(time.Now().Unix(), 0)
	}
	for c = 1; c <= 20; c++ {
		t2 = t1.Add(time.Minute).Add(time.Second)
		blksArr = append(blksArr, getNewBlock(getNewHeader(int64(c), isRaft, t1), txnsList(c*1000)))
		t1 = t2
	}

	for _, b := range blksArr {
		tm.readBlock(b)
	}
	tm.Stop()

	tm.printTPS()
	var expTxnCnt uint64 = 190000
	var expBlkCnt uint64 = 19
	var expTps uint32 = 166
	expTpsRecs := 19

	lr := len(tm.tpsRecs)

	assert.Equal(len(tm.tpsRecs), expTpsRecs, "tps record count mismatch")

	txnCnt := tm.tpsRecs[lr-1].txns
	blkCnt := tm.tpsRecs[lr-1].blks
	tps := tm.tpsRecs[lr-1].tps

	assert.Equal(txnCnt, expTxnCnt, "total txn count mismatch")
	assert.Equal(blkCnt, expBlkCnt, "block count mismatch")
	assert.Equal(tps, expTps, "tps mismatch")

	if tpsFile, err := os.Open(tempFile); err != nil {
		t.Errorf("opening tps file %s failed", tempFile)
	} else {
		lineCnt := 0
		expLineCnt := 20
		firstLine := "localTime,refTime,TPS,TxnCount,BlockCount"
		lineStr := ""
		var lineStrArr []string
		scanner := bufio.NewScanner(tpsFile)
		for scanner.Scan() {
			lineCnt++
			lineStr = scanner.Text()
			lineStrArr = strings.Split(lineStr, ",")
			if lineCnt == 1 {
				assert.Equal(len(lineStrArr), 5, "tps report file header fields mismatch")
				assert.Equal(lineStr, firstLine, "tps report file header data mismatch")
			}
		}

		assert.Equal(lineCnt, expLineCnt, "tps report file lines")

		if ftps, err := strconv.ParseInt(lineStrArr[2], 10, 32); err != nil {
			t.Errorf("tps data in file is wrong - tps is not a valid number")
		} else {
			assert.Equal(uint32(ftps), expTps, "tps report file - tps")
		}

		if ftxn, err := strconv.ParseInt(lineStrArr[3], 10, 64); err != nil {
			t.Errorf("tps data in file is wrong - tps is not a valid number")
		} else {
			assert.Equal(uint64(ftxn), expTxnCnt, "tps report file - tps")
		}

		if fblk, err := strconv.ParseInt(lineStrArr[4], 10, 64); err != nil {
			t.Errorf("tps data in file is wrong - tps is not a valid number")
		} else {
			assert.Equal(uint64(fblk), expBlkCnt, "tps report file - tps")
		}

	}

}
func TestTPSForIbft(t *testing.T) {
	f, _ := ioutil.TempFile("", "_tps_test")
	fn := f.Name()
	t.Log(fn)
	testTPS(false, t)
}

func TestTPSForRaft(t *testing.T) {
	testTPS(true, t)
}

func TestTime(t *testing.T) {
	t1 := time.Now()
	t.Log(t1.Format("2006-01-02T15:04:05-0700"))
}

func TestEexec(t *testing.T) {
	args := []string{"-ltr", "/Users/amalraj.manigmail.com/"}
	cmd := exec.Command("ls", args...)
	r, err := cmd.Output()
	if err != nil {
		log.Printf("failed to exec ls %v", err)
	}
	log.Printf("result: %s\n", string(r))

}

func TestExecAws(t *testing.T) {
	mySession := session.Must(session.NewSession())
	// Create a CloudWatch client with additional configuration
	svc := cloudwatch.New(mySession, aws.NewConfig().WithRegion("ap-southeast-1"))
	t.Log(svc)
	var pmd *cloudwatch.PutMetricDataInput
	var mdn *cloudwatch.MetricDatum
	dname := "Instance"
	dvalue := "NodeX"
	nspace := "stX-q24-X.Y.Z"
	mname := "TPS"
	var tps uint64 = 1000
	var value float64 = float64(tps)
	ts := time.Now()
	dimension := &cloudwatch.Dimension{Name: &dname, Value: &dvalue}
	mdn = &cloudwatch.MetricDatum{
		Dimensions: []*cloudwatch.Dimension{dimension},
		MetricName: &mname,
		Timestamp:  &ts,
		Value:      &value,
	}
	pmd = &cloudwatch.PutMetricDataInput{
		MetricData: []*cloudwatch.MetricDatum{mdn},
		Namespace:  &nspace,
	}
	svc.PutMetricData(pmd)
}
