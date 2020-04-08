package com.quorum.stresstest;

import java.io.BufferedReader;
import java.io.FileReader;
import java.io.IOException;
import java.sql.Time;
import java.util.ArrayList;
import java.util.Date;
import java.util.List;

class InputData {
    String endPoint;

    public InputData(String endPoint, String fromAcct, String toAcct, long nonce, long loopCnt, int delayAfter, int delay) {
        this.endPoint = endPoint;
        this.fromAcct = fromAcct;
        this.toAcct = toAcct;
        this.nonce = nonce;
        this.loopCnt = loopCnt;
        this.delayAfter = delayAfter;
        this.delay = delay;
    }

    @Override
    public String toString() {
        return "InputData{" +
                "endPoint='" + endPoint + '\'' +
                ", fromAcct='" + fromAcct + '\'' +
                ", toAcct='" + toAcct + '\'' +
                ", nonce=" + nonce +
                ", loopCnt=" + loopCnt +
                ", delayAfter=" + delayAfter +
                ", delay=" + delay +
                '}';
    }

    String fromAcct;
    String toAcct;
    long nonce;
    long loopCnt;
    int delayAfter;
    int delay;
}

public class StartTest {

    public static List<InputData> getDataFromFile(String fileName) {
        BufferedReader reader;
        List<InputData> dataArr = new ArrayList<InputData>();
        try {
            reader = new BufferedReader(new FileReader(fileName));

            while (true) {
                // read next line
                String line = reader.readLine();
                if(line == null)
                    break;
                if(line.trim().startsWith("#"))
                    continue;

                String[] flds = line.split(",");
                if (flds.length != 7) {
                    System.out.println(line);
                    System.out.println("invalid number of fields in the input file");
                    System.out.println("endPoint,fromAcct,toAcct,nonce,loopCount,delayAfter,delay");
                    System.exit(1);
                }
                String ep = flds[0].trim();
                String fa = flds[1].trim();
                String ta = flds[2].trim();
                long n = Long.parseLong(flds[3].trim());
                long lc = Long.parseLong(flds[4].trim());
                int da = Integer.parseInt(flds[5].trim());
                int d = Integer.parseInt(flds[6].trim());
                dataArr.add(new InputData(ep, fa, ta, n, lc, da, d));
            }
            reader.close();
        } catch (IOException e) {
            System.out.println("ERROR: reading file " + fileName + " " + e.getMessage());
        }
        return dataArr;
    }

    public static void main(String args[]) {
        if (args.length == 0) {
            System.out.println("usage: java StartTest <inputFileName>");
            System.exit(1);
        }
        long t1 = System.currentTimeMillis();
        System.out.println(new Date().toString() + " stress test start...");
        List<InputData> dataArr = getDataFromFile(args[0]);
        int dataCnt = dataArr.size();
        Thread[] threadArr = new Thread[dataCnt];
        for (int k = 0; k < threadArr.length; ++k) {
            InputData d = dataArr.get(k);
            SendSignedTransaction tx = new SendSignedTransaction(d.endPoint, 1, d.nonce, d.fromAcct, d.toAcct, d.loopCnt, d.delayAfter, d.delay);
            threadArr[k] = new Thread(tx);
            //System.out.println("thread " + k + " created for input " + d.toString());
        }


        for (int k = 0; k < threadArr.length; ++k) {

            System.out.println("thread " + k + " started for input " + dataArr.get(k).toString());
            threadArr[k].start();
        }

        for (int k = 0; k < threadArr.length; ++k) {

            System.out.println("thread " + k + " waiting for input " + dataArr.get(k).toString() + " to finish...");

            try {
                threadArr[k].join();
            } catch (InterruptedException e) {
                System.out.println("ERROR: thread " + k + " exception:" + e.getMessage());
            }
            System.out.println("thread " + k + " finished for input " + dataArr.get(k).toString());
        }

        long t2 = System.currentTimeMillis();

        System.out.println(new Date().toString() + " stress test end. total time taken " + (t2-t1) + "ms");

    }
}
