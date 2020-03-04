package com.quorum.stresstest;

import org.apache.http.StatusLine;
import org.apache.http.client.methods.CloseableHttpResponse;
import org.apache.http.client.methods.HttpPost;
import org.apache.http.entity.ContentType;
import org.apache.http.entity.StringEntity;
import org.apache.http.impl.client.CloseableHttpClient;
import org.apache.http.impl.client.HttpClients;
import org.apache.http.util.EntityUtils;
import org.json.JSONObject;

import java.io.IOException;
import java.util.Map;

class ReqResult {
    public String data;

    public ReqResult(String data, boolean ok, String error) {
        this.data = data;
        this.ok = ok;
        this.error = error;
    }

    public boolean ok;
    public String error;
}

public class SendSignedTransaction implements Runnable {
    private String endPoint;
    private int value;
    private long nonce;
    private String fromAcct;
    private String toAcct;
    private long loopCount;
    private int delayAfter;
    private int delay;
    private boolean debug = false;
    private long id = 0;
    private CloseableHttpClient httpClient;
    private HttpPost postReq;

    public SendSignedTransaction(String endPoint, int value, long nonce, String fromAcct, String toAcct, long loopCount, int delayAfter, int delay) {
        this.endPoint = endPoint;
        this.value = value;
        this.nonce = nonce;
        this.fromAcct = fromAcct;
        this.toAcct = toAcct;
        this.loopCount = loopCount;
        this.delayAfter = delayAfter;
        this.delay = delay;
        this.httpClient = HttpClients.createDefault();
        this.postReq = new HttpPost(endPoint);

    }

    private ReqResult sendPostRequest(String postData, long c, boolean signTx) {
        if (debug)
            System.out.println("postData=" + postData);
        postReq.setEntity(new StringEntity(postData, ContentType.APPLICATION_JSON));

        CloseableHttpResponse response = null;
        try {
            response = httpClient.execute(postReq);

        } catch (IOException ex) {
            System.out.println("ERROR: count=" + c + " exception:" + ex.getMessage());
            return new ReqResult("", false, ex.getMessage());
        }

        if (response != null && response.getStatusLine().getStatusCode() == 200) {
            try {
                String result = EntityUtils.toString(response.getEntity());
                if (debug)
                    System.out.println("Executed post req RESULT:\n" + result);
                JSONObject resultObject = new JSONObject(result);
                Map<String, Object> resMap = resultObject.toMap();
                if (debug)
                    System.out.println("jsonObj:" + resultObject.toMap().toString());
                if (resMap.containsKey("error")) {
                    Map<String, Object> errMap = (Map<String, Object>) resMap.get("error");
                    System.out.println("ERROR: sendTx failed with error code=" + errMap.get("code") + " reason=" + errMap.get("message"));
                    return new ReqResult("", false, (String) errMap.get("message"));
                } else {
                    if (signTx) {
                        Map<String, Object> rMap = (Map<String, Object>) resMap.get("result");
                        String rawData = (String) rMap.get("raw");
                        if(debug)
                            System.out.println("rawData=" + rawData);
                        long endTime = System.currentTimeMillis();
                        return new ReqResult(rawData, true, "");
                    }
                    return new ReqResult(result, true, "");
                }

            } catch (IOException e) {
                System.out.println("ERROR: count=" + c + " reading result - exception:" + e.getMessage());
                return new ReqResult("", false, e.getMessage());
            }
        } else {
            StatusLine sl = response.getStatusLine();
            System.out.println("ERROR: response failed - " + sl.getStatusCode() + " - " + sl.getReasonPhrase());
            return new ReqResult("", false, sl.getReasonPhrase());
        }

    }


    public void run() {
        long startPTime = System.currentTimeMillis();
        long txCnt = 0;
        for (long c = 1; c <= loopCount; c++) {
            long t1 = System.currentTimeMillis();
            String nonceHex = Long.toHexString(nonce);
            String postData = "{\"jsonrpc\":\"2.0\",\"method\":\"eth_signTransaction\",\"params\":[{\"from\":\"" + fromAcct +
                    "\",\"to\":\"" + toAcct + "\", \"value\": \"0x1\", \"gasPrice\":\"0x0\",\"gas\":\"0x47b760\", \"nonce\":\"0x" + nonceHex + "\"}],\"id\":\"" + id + "\"}";
            ReqResult signReq = sendPostRequest(postData, c, true);
            if (signReq.ok) {
                ++id;
                String rawData = signReq.data;
                String sendTxData = "{\"jsonrpc\":\"2.0\",\"method\":\"eth_sendRawTransaction\",\"params\":[\"" + rawData + "\"],\"id\":\"" + id + "\"}";
                ReqResult txRes = sendPostRequest(sendTxData, c, false);
                if (txRes.ok) {
                    if(debug)
                        System.out.println(c+".Tx sent. " + txRes.data);
                    ++nonce;
                    ++txCnt;
                }
            }
            long t2 = System.currentTimeMillis();
            if (debug)
                System.out.println("iter=" + c + " time taken " + (t2 - t1) + " milliseconds");
            if (delayAfter>0 && c % delayAfter == 0) {
                try {
                    if(debug)
                        System.out.println("sleeping for "+ delay + " milliseconds");
                    Thread.sleep(delay);
                } catch (InterruptedException e) {
                    System.out.println("ERROR: waiting c=" + c + " exception:" + e.getMessage());
                }
            }
        }
        long endPTime = System.currentTimeMillis();
        String timeTaken = "Total time taken " + (endPTime - startPTime) + " milliseconds";
        System.out.println(endPoint+" " + txCnt + " txns sent successfully! nonce="+nonce+ ". " + timeTaken);
        System.out.println();
    }
}
