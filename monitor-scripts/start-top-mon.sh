#!/bin/bash
cd /home/ec2-user/node/
ts=`date +"%d%m%Y%H%M%S"`
glf="top-cmd_geth_$ts.log"
jlf="top-cmd_java_$ts.log"
#echo "start cpu mem usage monitoring.." > $lf
p1=`ps -eaf|grep geth|grep -v grep|grep -v attach|tr -s " " " "|cut -d" " -f2`
#echo "geth process id is $p1" >>
p2=`ps -eaf|grep tessera|grep -v grep|tr -s " " " "|cut -d" " -f2`
#echo "java process id is $p2"
while true
do
	timeStamp=`date +'%d/%m/%Y %H:%M:%S'`
#	echo $timeStamp
	gethrec=`top -bcn 1 -p $p1 | tail -1| tr -s " "| awk '{print $1" " $6" " $9" " $10}'`
#	echo "geth rec is ---- $gethrec"
	goutput="$timeStamp $gethrec"
	echo $goutput >> $glf
	javarec=`top -bcn 1 -p $p2 | tail -1| tr -s " "| awk '{print $1" " $6" " $9" " $10}'`
#	echo "java rec is $javarec"
	joutput="$timeStamp $javarec"
#	echo $joutput
	echo $joutput >> $jlf
	sleep 60
done
