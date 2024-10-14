package main

import (
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"models"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"runtime"
	"strconv"
	"sync"
	"syscall"
	"time"
)

var url string = "http://"
var pollInterval int = 2
var reportInterval int = 10

const validateAdderPattern string = "^[a-zA-Z/ ]{1,100}:[0-9]{1,10}$"

const contentTypeSendMetric string = "text/plain"
const minRandomValue float64 = 1.0
const maxRandomValue float64 = 999.0

type Options struct {
	PORT           string
	reportInterval int
	pollInterval   int
}

/*type responseData struct {
	r   *http.Response
	err error
}*/

var wg sync.WaitGroup
var m models.Monitor
var options Options
var httpClient *http.Client

// var dataChannel chan responseData
var gauges []models.Gauge
var counters map[string]models.Counter

func main() {

	err := initialization()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	wg.Add(1)
	go collectMetrics()
	wg.Add(1)
	go sendMetrics()
	wg.Wait()
}

func collectMetrics() {

	defer wg.Done()

	channelCancel := make(chan os.Signal, 1)
	signal.Notify(channelCancel, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	for {
		select {
		case <-channelCancel:
			return
		case <-time.After(time.Duration(pollInterval) * time.Second):
			setValuesMonitor()
		}
	}
}

func sendMetrics() {

	defer wg.Done()
	channelCancel := make(chan os.Signal, 1)
	signal.Notify(channelCancel, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	for {
		select {
		case <-channelCancel:
			return
		case <-time.After(time.Duration(reportInterval) * time.Second):
			go doReqSendMetrics()
			//case answer := <-dataChannel:
			//	parseAnswer(&answer)
		}
	}
}

/*func parseAnswer(answer *responseData) {

	err := answer.err
	resp := answer.r
	if err != nil {
		fmt.Println("Error select:", err)
	}
	defer resp.Body.Close()
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error read data:", err)
	}

	fmt.Printf("Server Response: %s\n", responseBody)
}*/

func doReqSendMetrics() {

	tmp_url := url + "/update/" + "counter/"
	fmt.Println(counters)
	for _, metric := range counters {
		sendMetricEndpoint(tmp_url + metric.Name + "/" + fmt.Sprintf("%v", metric.Value))
	}
	tmp_url = url + "/update/" + "gauge/"
	for _, metric := range gauges {
		sendMetricEndpoint(tmp_url + metric.Name + "/" + strconv.FormatFloat(metric.Value, 'f', -1, 64))
	}
}

func sendMetricEndpoint(endpoint string) {
	req, _ := http.NewRequest("POST", endpoint, nil)
	req.Header.Set("Content-Type", contentTypeSendMetric)
	_, err := httpClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	//response, _ := httpClient.Do(req)
	//fmt.Println(response)
	/*if err != nil {
		dataChannel <- responseData{nil, err}
	} else {
		pack := responseData{response, err}
		dataChannel <- pack
	}*/
}

func initialization() error {

	tr := &http.Transport{}
	httpClient = &http.Client{Transport: tr}
	err := parseFlags()
	if err != nil {
		return err
	}
	err = getENV()
	if err != nil {
		return err
	}
	url = url + options.PORT
	pollInterval = options.pollInterval
	reportInterval = options.reportInterval
	rand.Seed(time.Now().Unix())
	m.Init()
	//dataChannel = make(chan responseData, 1)
	return nil
}

func getENV() error {
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		if isValidAddr(envRunAddr) {
			options.PORT = envRunAddr
		} else {
			return errors.New("addr is not valid")
		}
	}
	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		value, err := strconv.Atoi(envReportInterval)
		if err != nil {
			return err
		} else {
			options.reportInterval = value
		}
	}
	if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
		value, err := strconv.Atoi(envPollInterval)
		if err != nil {
			return err
		} else {
			options.pollInterval = value
		}
	}
	return nil
}

func parseFlags() error {
	flag.StringVar(&options.PORT, "a", "localhost:8080", "Port to listen on.")
	flag.IntVar(&options.pollInterval, "p", pollInterval, "Frequency of sending metrics to the server.")
	flag.IntVar(&options.reportInterval, "r", reportInterval, "Frequency of polling metrics from the runtime package.")
	flag.Parse()

	if !isValidAddr(options.PORT) {
		return errors.New("addr is not valid")
	}
	return nil
}

func randomF64(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func setValuesMonitor() {

	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)

	m.PollCount.Value += 1

	counters = make(map[string]models.Counter, 1)
	counters["PollCount"] = m.PollCount
	fmt.Println(counters["PollCount"])

	gauges = make([]models.Gauge, 0, 27)

	m.Alloc.Value = float64(rtm.Alloc)
	m.BuckHashSys.Value = float64(rtm.BuckHashSys)
	m.Frees.Value = float64(rtm.Frees)
	m.GCCPUFraction.Value = rtm.GCCPUFraction
	m.GCSys.Value = float64(rtm.GCSys)
	m.HeapAlloc.Value = float64(rtm.HeapAlloc)
	m.HeapIdle.Value = float64(rtm.HeapIdle)
	m.HeapInuse.Value = float64(rtm.HeapInuse)
	m.HeapObjects.Value = float64(rtm.HeapObjects)
	m.HeapReleased.Value = float64(rtm.HeapReleased)
	m.HeapSys.Value = float64(rtm.HeapSys)
	m.LastGC.Value = float64(rtm.LastGC)
	m.Lookups.Value = float64(rtm.Lookups)
	m.MCacheInuse.Value = float64(rtm.MCacheInuse)
	m.MCacheSys.Value = float64(rtm.MCacheSys)
	m.MSpanInuse.Value = float64(rtm.MSpanInuse)
	m.MSpanSys.Value = float64(rtm.MSpanSys)
	m.Mallocs.Value = float64(rtm.Mallocs)
	m.NextGC.Value = float64(rtm.NextGC)
	m.NumForcedGC.Value = float64(rtm.NumForcedGC)
	m.NumGC.Value = float64(rtm.NumGC)
	m.OtherSys.Value = float64(rtm.OtherSys)
	m.PauseTotalNs.Value = float64(rtm.PauseTotalNs)
	m.StackInuse.Value = float64(rtm.StackInuse)
	m.StackSys.Value = float64(rtm.StackSys)
	m.Sys.Value = float64(rtm.Sys)
	m.TotalAlloc.Value = float64(rtm.TotalAlloc)
	m.RandomValue.Value = randomF64(minRandomValue, maxRandomValue)

	gauges = append(gauges, m.Alloc)
	gauges = append(gauges, m.BuckHashSys)
	gauges = append(gauges, m.Frees)
	gauges = append(gauges, m.GCCPUFraction)
	gauges = append(gauges, m.GCSys)
	gauges = append(gauges, m.HeapAlloc)
	gauges = append(gauges, m.HeapIdle)
	gauges = append(gauges, m.HeapInuse)
	gauges = append(gauges, m.HeapObjects)
	gauges = append(gauges, m.HeapReleased)
	gauges = append(gauges, m.HeapSys)
	gauges = append(gauges, m.LastGC)
	gauges = append(gauges, m.Lookups)
	gauges = append(gauges, m.MCacheInuse)
	gauges = append(gauges, m.MCacheSys)
	gauges = append(gauges, m.MSpanInuse)
	gauges = append(gauges, m.MSpanSys)
	gauges = append(gauges, m.Mallocs)
	gauges = append(gauges, m.NextGC)
	gauges = append(gauges, m.NumForcedGC)
	gauges = append(gauges, m.NumGC)
	gauges = append(gauges, m.OtherSys)
	gauges = append(gauges, m.PauseTotalNs)
	gauges = append(gauges, m.StackInuse)
	gauges = append(gauges, m.StackSys)
	gauges = append(gauges, m.Sys)
	gauges = append(gauges, m.TotalAlloc)
	gauges = append(gauges, m.RandomValue)
}

func MatchString(pattern string, s string) (matched bool, err error) { //  in a separate package

	re, err := regexp.Compile(pattern)
	if err == nil {
		return re.MatchString(s), nil
	} else {
		return false, err
	}

}

func isValidAddr(addr string) bool { //  in a separate package

	var pattern string = validateAdderPattern

	res, err := MatchString(pattern, addr)
	if err == nil && res == true {
		return true
	} else {
		return false
	}

}
