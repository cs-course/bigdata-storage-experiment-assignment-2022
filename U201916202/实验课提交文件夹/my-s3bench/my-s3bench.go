package main

import (
	"bytes"
	"crypto/rand"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Params struct {
	putReqChan       chan *s3.PutObjectInput
	getReqChan       chan *s3.GetObjectInput
	putRespChan      chan Response
	getRespChan      chan Response
	accessKey        string
	secretKey        string
	bucket           string
	endpoint         string
	numOfClients     int
	numOfSamples     int
	objectNamePrefix string
	objectSize       int
	sampleData       []byte

	reqMode int
}

type Response struct {
	err      error
	duration time.Duration
}

type Result struct {
	opcode          string
	totalDuration   time.Duration
	totalTransBytes int64
	numOfErrors     int
	durationList    []float64

	numOfClients int
	numOfSamples int
	objectSize   int
}

func strToIntegerList(str string) []int {
	res := make([]int, 0)
	for _, s := range strings.Split(str, ",") {
		i, err := strconv.Atoi(s)
		if err != nil {
			fmt.Println("Please check the endpoint you input!")
			os.Exit(-1)
		}
		res = append(res, i)
	}
	return res
}

func main() {
	fmt.Println("this is my s3bench")
	accessKey := flag.String("accessKey", "", "The access-key of the s3 server.")
	secretKey := flag.String("secretKey", "", "The secret-key of the s3 server.")
	bucket := flag.String("bucket", "", "The bucket where you run your test in.")
	endpoint := flag.String("endpoint", "", "The endpoint of the s3 server - http://IP:PORT")
	numOfClientsStr := flag.String("numOfClients", "", "The numbers of clients set up during the test.")
	numOfSamplesStr := flag.String("numOfSamples", "", "The numbers of test samples you want to put/get during write/read test.")
	objectNamePrefix := flag.String("objectNamePrefix", "", "The prefix of the object name.")
	objectSizeStr := flag.String("objectSize", "", "The size of object. (Bytes)")
	reqMode := flag.Int("requestMode", 0, "The mode of request. 0:Normal, 1:Hedged, 2:Tied")

	flag.Parse()

	numOfClientsList := strToIntegerList(*numOfClientsStr)
	numOfSamplesList := strToIntegerList(*numOfSamplesStr)
	objectSizeList := strToIntegerList(*objectSizeStr)

	if *endpoint == "" {
		fmt.Println("Please check the endpoint you input!")
		os.Exit(-1)
	}

	//fmt.Println(*accessKey, *secretKey, *bucket, *endpoint, *numOfClients, *numOfSamples, *objectNamePrefix, *objectSize)

	writeTestResults := make([]*Result, 0)
	readTestResults := make([]*Result, 0)

	for _, numOfClients := range numOfClientsList {
		for _, numOfSamples := range numOfSamplesList {
			for _, objectSize := range objectSizeList {
				params := Params{
					putReqChan:  make(chan *s3.PutObjectInput, numOfSamples),
					getReqChan:  make(chan *s3.GetObjectInput, numOfSamples),
					putRespChan: make(chan Response, numOfSamples),
					getRespChan: make(chan Response, numOfSamples),

					accessKey:        *accessKey,
					secretKey:        *secretKey,
					bucket:           *bucket,
					endpoint:         *endpoint,
					numOfClients:     numOfClients,
					numOfSamples:     numOfSamples,
					objectNamePrefix: *objectNamePrefix,
					objectSize:       objectSize,
					reqMode:          *reqMode,
				}

				fmt.Printf("params:{ accessKey:%s  secretKey:%s  bucket:%s  endpoint:%s  "+
					"numOfClients:%d  numOfSamples:%d  objectNamePrefix:%s  objectSize:%d }",
					params.accessKey, params.secretKey, params.bucket, params.endpoint,
					params.numOfClients, params.numOfSamples, params.objectNamePrefix, params.objectSize)

				fmt.Printf("Generating the sample data in memory whose size is %d...\n", params.objectSize)
				timeGenBegin := time.Now()
				params.sampleData = make([]byte, params.objectSize)
				_, errorMsg := rand.Read(params.sampleData)
				if errorMsg != nil {
					fmt.Println("Error! The errorMsg is : ", errorMsg)
					os.Exit(-1)
				}
				fmt.Printf("Done. Time cost : %s\n", time.Since(timeGenBegin))

				params.startClients()

				writeTestResult := params.runTest("Write")

				readTestResult := params.runTest("Read")

				writeTestResult.printResult()
				readTestResult.printResult()

				//writeTestResult.showLineChart(5, 5)
				//readTestResult.showLineChart(5, 5)

				writeTestResult.numOfClients = numOfClients
				writeTestResult.numOfSamples = numOfSamples
				writeTestResult.objectSize = objectSize
				readTestResult.numOfClients = numOfClients
				readTestResult.numOfSamples = numOfSamples
				readTestResult.objectSize = objectSize

				writeTestResults = append(writeTestResults, &writeTestResult)
				readTestResults = append(readTestResults, &readTestResult)
			}
		}
	}

	showMultiLineChart(writeTestResults, 5, 5)
	showMultiLineChart(readTestResults, 5, 5)
}

func (result *Result) getPercentData(percent int) float64 {
	index := int(0.01*(float32(percent*len(result.durationList)))) - 1
	if index < 0 {
		index = 0
	}
	return result.durationList[index]
}

func (result *Result) printResult() {
	fmt.Println()
	fmt.Printf("Result Summary for %s Operations\n", result.opcode)
	fmt.Printf("Total Transferred:\t%.3f MB\n", float64(result.totalTransBytes)/(1<<20))
	fmt.Printf("Total Throughput:\t%.3f MB/s \n", float64(result.totalTransBytes)/(1<<20)*result.totalDuration.Seconds())
	fmt.Printf("Total Duration:\t\t%.5f s\n", result.totalDuration.Seconds())
	fmt.Printf("Number of Errors: %d\t\n", result.numOfErrors)
	fmt.Println("----------------------------------------------------")
	fmt.Printf("%s times Max  :\t%.5f s\n", result.opcode, result.getPercentData(100))
	fmt.Printf("%s times 99%%th:\t%.5f s\n", result.opcode, result.getPercentData(99))
	fmt.Printf("%s times 90%%th:\t%.5f s\n", result.opcode, result.getPercentData(90))
	fmt.Printf("%s times 75%%th:\t%.5f s\n", result.opcode, result.getPercentData(75))
	fmt.Printf("%s times 50%%th:\t%.5f s\n", result.opcode, result.getPercentData(50))
	fmt.Printf("%s times 25%%th:\t%.5f s\n", result.opcode, result.getPercentData(25))
	fmt.Printf("%s times Min  :\t%.5f s\n", result.opcode, result.durationList[0])
	fmt.Println()
}

func showMultiLineChart(results []*Result, first int, grad int) {
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithYAxisOpts(opts.YAxis{
			Name: "Latency(ms)",
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Name: "Percent(%)",
		}),
		charts.WithParallelComponentOpts(opts.ParallelComponent{
			Left:   "15%",
			Right:  "13%",
			Bottom: "10%",
			Top:    "20%",
		}),
		charts.WithLegendOpts(opts.Legend{Show: true}),
	)

	xAxis := make([]string, 0)
	for i := first; i <= 100; i += grad {
		xAxis = append(xAxis, fmt.Sprintf("%d", i))
	}

	line.SetXAxis(xAxis).
		SetSeriesOptions(
			charts.WithLineChartOpts(opts.LineChart{
				Smooth: true,
			}),
		)

	for _, result := range results {
		line.AddSeries(fmt.Sprintf("%d-%d-%d", result.numOfClients, result.numOfSamples, result.objectSize), result.getLineData(first, grad))
	}

	f, _ := os.Create(fmt.Sprintf("%s.html", results[0].opcode))
	line.Render(f)
	fmt.Printf("The line chart of the test result is saved as %s.html\n", results[0].opcode)
}

func (result *Result) showLineChart(first int, grad int) {
	line := charts.NewLine()

	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: fmt.Sprintf("This chart show the tail latency in my-s3bench (%s test).", result.opcode),
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Name: "Latency(ms)",
			SplitLine: &opts.SplitLine{
				Show: false,
			},
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Name: "Percent(%)",
		}),
	)

	xAxis := make([]string, 0)
	for i := first; i <= 100; i += grad {
		xAxis = append(xAxis, fmt.Sprintf("%d", i))
	}

	line.SetXAxis(xAxis).
		AddSeries("slice", result.getLineData(first, grad),
			charts.WithLabelOpts(opts.Label{Show: true, Position: "top"})).
		SetSeriesOptions(
			charts.WithLineChartOpts(opts.LineChart{
				Smooth: true,
			}),
		)

	f, _ := os.Create(fmt.Sprintf("%s.html", result.opcode))
	line.Render(f)
	fmt.Printf("The line chart of the %s test result is saved as %s.html", result.opcode)
}

func (result *Result) getLineData(first int, grad int) []opts.LineData {
	lineData := make([]opts.LineData, 0)
	for i := first; i <= 100; i += grad {
		lineData = append(lineData, opts.LineData{Value: 1000 * result.getPercentData(i)})
	}
	return lineData
}

func (params *Params) runTest(opcode string) Result {
	fmt.Printf("Runnig %s Test...\n", opcode)
	timeTestBegin := time.Now()
	go params.produceRequest(opcode)
	result := Result{
		opcode:          opcode,
		totalTransBytes: 0,
		numOfErrors:     0,
		durationList:    make([]float64, 0, params.numOfSamples),
	}
	if opcode == "Write" {
		for i := 0; i < params.numOfSamples; i++ {
			resp := <-params.putRespChan
			if resp.err != nil {
				result.numOfErrors++
			} else {
				result.totalTransBytes += int64(params.objectSize)
				result.durationList = append(result.durationList, resp.duration.Seconds())
			}
		}
	} else if opcode == "Read" {
		for i := 0; i < params.numOfSamples; i++ {
			resp := <-params.getRespChan
			if resp.err != nil {
				result.numOfErrors++
			} else {
				result.totalTransBytes += int64(params.objectSize)
				result.durationList = append(result.durationList, resp.duration.Seconds())
			}
		}
	}

	result.totalDuration = time.Since(timeTestBegin)
	sort.Float64s(result.durationList)

	fmt.Println("Done.")

	return result
}

func (params *Params) produceRequest(opcode string) {

	if opcode == "Write" {
		for i := 0; i < params.numOfSamples; i++ {
			params.putReqChan <- &s3.PutObjectInput{
				Body:   bytes.NewReader(params.sampleData),
				Bucket: aws.String(params.bucket),
				Key:    aws.String(fmt.Sprintf("%s_%d", params.objectNamePrefix, i)),
			}
		}
		close(params.putReqChan)
	} else if opcode == "Read" {
		for i := 0; i < params.numOfSamples; i++ {
			params.getReqChan <- &s3.GetObjectInput{
				Bucket: aws.String(params.bucket),
				Key:    aws.String(fmt.Sprintf("%s_%d", params.objectNamePrefix, i)),
			}
		}
		close(params.getReqChan)
	} else {
		fmt.Println("Check the opcode!")
	}

}

func (params *Params) startClients() {
	fmt.Printf("Setting up %d clients...\n", params.numOfClients)
	cred := credentials.NewStaticCredentials(params.accessKey, params.secretKey, "")
	config := &aws.Config{
		Endpoint:         aws.String(params.endpoint),
		Credentials:      cred,
		Region:           aws.String("igneous-test"),
		S3ForcePathStyle: aws.Bool(true),
	}
	for i := 0; i < params.numOfClients; i++ {
		go params.clientRoutine(config)
	}
}

func (params *Params) clientRoutine(config *aws.Config) {
	svc := s3.New(session.New(config))
	var timeT time.Duration = 30 * time.Millisecond
	for req := range params.putReqChan {
		putTimeBegin := time.Now()
		_, err := svc.PutObject(req)
		duration := time.Since(putTimeBegin)

		if params.reqMode == 1 && duration > timeT { //对冲请求模式
			putTimeBegin2 := time.Now()
			_, err2 := svc.PutObject(req)
			// 由于第二个是在第一个请求超时后才发送
			// 因此第二个请求的相对响应时间是绝对时间+超时阈值
			duration2 := time.Since(putTimeBegin2) + timeT
			// 若第二个请求先到达且响应无报错，则将第二个请求的响应作为本次服务请求的响应
			if duration2 < duration && err2 == nil {
				duration = duration2
				err = err2
			}
		} else if params.reqMode == 2 { //关联请求模式
			//time.Sleep(10 * time.Millisecond)
			putTimeBegin2 := time.Now()
			_, err2 := svc.PutObject(req)
			// 关联请求同时发送两个请求，当其中之一到达时，立即取消另外一个请求
			duration2 := time.Since(putTimeBegin2)
			// 若第二个请求先到达且响应无报错，则将第二个请求的响应作为本次服务请求的响应
			if duration2 < duration && err2 == nil {
				duration = duration2
				err = err2
			}
		}
		params.putRespChan <- Response{err: err, duration: duration}
	}

	for req := range params.getReqChan {
		getTimeBegin := time.Now()
		_, err := svc.GetObject(req)
		duration := time.Since(getTimeBegin)

		if params.reqMode == 1 && duration > timeT { //对冲请求模式
			getTimeBegin2 := time.Now()
			_, err2 := svc.GetObject(req)
			// 由于第二个是在第一个请求超时后才发送
			// 因此第二个请求的相对响应时间是绝对时间+超时阈值
			duration2 := time.Since(getTimeBegin2) + timeT
			// 若第二个请求先到达且响应无报错，则将第二个请求的响应作为本次服务请求的响应
			if duration2 < duration && err2 == nil {
				duration = duration2
				err = err2
			}
		} else if params.reqMode == 2 { //关联请求模式
			//time.Sleep(10 * time.Millisecond)
			getTimeBegin2 := time.Now()
			_, err2 := svc.GetObject(req)
			// 关联请求同时发送两个请求，当其中之一到达时，立即取消另外一个请求
			duration2 := time.Since(getTimeBegin2)
			// 若第二个请求先到达且响应无报错，则将第二个请求的响应作为本次服务请求的响应
			if duration2 < duration && err2 == nil {
				duration = duration2
				err = err2
			}
		}

		params.getRespChan <- Response{err: err, duration: duration}
	}
}
