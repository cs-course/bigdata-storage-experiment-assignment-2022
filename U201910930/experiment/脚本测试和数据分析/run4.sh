#!/bin/bash
#test shell-script 1 for NumClient=1 NumSample=[1,1280] ObjectSize=1KB

# Locate s3bench



if [ -n "$GOPATH" ]; then
    s3bench=$GOPATH/bin/s3bench
fi


endpoint="http://127.0.0.1:9001"
bucket="exterminate"
ObjectNamePrefix="loadgen"
AccessKey="admin"
AccessSecret="password"
filepath="minio_server_4.txt"

declare  -a  NumClient
declare  -a  NumSample
declare  -a  ObjectSize

NumClient=(1    1    1    1     1     1     1      1      1)
NumSample=(5  10  20  40   80   160   320    640    1280)
ObjectSize=(1024 1024 1024 1024 1024 1024 1024 1024 1024)

#display run progress
progress=8

for(( i=0;i<${#NumClient[@]};i++)) 
do  
    # run sh
    ./s3bench -accessKey=$AccessKey -accessSecret=$AccessSecret -bucket=$bucket -endpoint=$endpoint \
    -numClients=${NumClient[i]} -numSamples=${NumSample[i]} -objectNamePrefix=$ObjectNamePrefix -objectSize=${ObjectSize[i]} >> $filepath
 
    echo -e "=========================================================> $i/$progress done\n" >> $filepath
    echo -e "=========================================================> $i/$progress done\n"
done
