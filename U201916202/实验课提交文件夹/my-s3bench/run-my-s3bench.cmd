my-s3bench ^
    -accessKey=hust ^
    -secretKey=hust_obs ^
    -bucket=loadgen ^
    -endpoint=http://127.0.0.1:9090 ^
    -numOfClients=1,4,8 ^
    -numOfSamples=256,512,1024 ^
    -objectNamePrefix=loadgen ^
    -objectSize=1024,4096 ^
    -requestMode=0
pause