package models

import (
	"errors"
	"fmt"
	"strconv"
)

type Gauge struct {
	Name  string
	Value float64
}

type Counter struct {
	Name  string
	Value int64
}

type MemStorage struct {
	gauges   map[string]Gauge
	counters map[string]Counter
}

func (m *MemStorage) Init() {
	m.gauges = make(map[string]Gauge)
	m.counters = make(map[string]Counter)
}

func (m *MemStorage) GetStringValueGaugeMetric(name string) (string, error) {

	val, ok := m.gauges[name]
	if ok {
		return strconv.FormatFloat(val.Value, 'f', -1, 64), nil
	} else {
		return "", errors.New("metric not found")
	}

}

func (m *MemStorage) GetStringValueCounterMetric(name string) (string, error) {
	val, ok := m.counters[name]
	if ok {
		return fmt.Sprintf("%v", val.Value), nil
	} else {
		return "", errors.New("metric not found")
	}
}

func (m *MemStorage) GetMapStringsAllMetrics() *map[string]string {
	mapMetrics := make(map[string]string)

	for key, value := range m.counters {
		mapMetrics[key] = fmt.Sprintf("%v", value.Value)
	}

	for key, value := range m.gauges {
		mapMetrics[key] = strconv.FormatFloat(value.Value, 'f', -1, 64)
	}

	return &mapMetrics
}

func (m *MemStorage) AddGauge(gauge *Gauge) {
	m.gauges[gauge.Name] = *gauge
}

func (m *MemStorage) AddCounter(counter *Counter) {

	val, ok := m.counters[counter.Name]
	if ok {
		var temp *Counter = new(Counter)
		temp.Name = val.Name
		temp.Value = val.Value + counter.Value
		m.counters[counter.Name] = *temp
	} else {
		m.counters[counter.Name] = *counter
	}
}

type Monitor struct {
	Alloc         Gauge
	TotalAlloc    Gauge
	BuckHashSys   Gauge
	Frees         Gauge
	Mallocs       Gauge
	Sys           Gauge
	GCCPUFraction Gauge
	GCSys         Gauge
	HeapAlloc     Gauge
	HeapIdle      Gauge
	HeapInuse     Gauge
	HeapObjects   Gauge
	HeapReleased  Gauge
	HeapSys       Gauge
	LastGC        Gauge
	Lookups       Gauge
	MCacheInuse   Gauge
	MCacheSys     Gauge
	MSpanInuse    Gauge
	MSpanSys      Gauge
	NextGC        Gauge
	NumForcedGC   Gauge
	NumGC         Gauge
	OtherSys      Gauge
	PauseTotalNs  Gauge
	StackInuse    Gauge
	StackSys      Gauge

	PollCount   Counter
	RandomValue Gauge
}

func (m *Monitor) Init() {

	//var users = [3]string{"Alloc", "TotalAlloc", "BuckHashSys"}
	//for index, value := range users{
	// f := reflect.Indirect(m).FieldByName(name)
	//}

	m.Alloc = Gauge{Name: "Alloc", Value: 0}
	m.BuckHashSys = Gauge{Name: "BuckHashSys", Value: 0}
	m.Frees = Gauge{Name: "Frees", Value: 0}
	m.GCCPUFraction = Gauge{Name: "GCCPUFraction", Value: 0}
	m.GCSys = Gauge{Name: "GCSys", Value: 0}
	m.HeapAlloc = Gauge{Name: "HeapAlloc", Value: 0}
	m.HeapIdle = Gauge{Name: "HeapIdle", Value: 0}
	m.HeapInuse = Gauge{Name: "HeapInuse", Value: 0}
	m.HeapObjects = Gauge{Name: "HeapObjects", Value: 0}
	m.HeapReleased = Gauge{Name: "HeapReleased", Value: 0}
	m.HeapSys = Gauge{Name: "HeapSys", Value: 0}
	m.LastGC = Gauge{Name: "LastGC", Value: 0}
	m.Lookups = Gauge{Name: "Lookups", Value: 0}
	m.MCacheInuse = Gauge{Name: "MCacheInuse", Value: 0}
	m.MCacheSys = Gauge{Name: "MCacheSys", Value: 0}
	m.MSpanInuse = Gauge{Name: "MSpanInuse", Value: 0}
	m.MSpanSys = Gauge{Name: "MSpanSys", Value: 0}
	m.Mallocs = Gauge{Name: "Mallocs", Value: 0}
	m.NextGC = Gauge{Name: "NextGC", Value: 0}
	m.NumForcedGC = Gauge{Name: "NumForcedGC", Value: 0}
	m.NumGC = Gauge{Name: "NumGC", Value: 0}
	m.OtherSys = Gauge{Name: "OtherSys", Value: 0}
	m.PauseTotalNs = Gauge{Name: "PauseTotalNs", Value: 0}
	m.StackInuse = Gauge{Name: "StackInuse", Value: 0}
	m.StackSys = Gauge{Name: "StackSys", Value: 0}
	m.Sys = Gauge{Name: "Sys", Value: 0}
	m.TotalAlloc = Gauge{Name: "TotalAlloc", Value: 0}

	m.PollCount = Counter{Name: "PollCount", Value: 0}
	m.RandomValue = Gauge{Name: "RandomValue", Value: 0}
}
