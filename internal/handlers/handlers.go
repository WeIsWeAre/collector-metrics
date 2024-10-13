package handlers

import (
	"fmt"
	"models"
	"net/http"
	"regexp"
	"service"
	"strconv"

	"github.com/gorilla/mux"
)

const PORT string = ":8080"
const metrics string = "gauge|counter"
const acceptedContentType string = "text/plain"

var vMetric validMetric

type validMetric struct {
	mtype        string
	mname        string
	mvalue       string
	mvalue_float float64
	mvalue_int   int64
}
type setMetricHandler struct {
	serv     service.Service
	memStore *models.MemStorage
}

type NotAllowedHandler struct{}

func (h NotAllowedHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	MethodNotAllowedHandler(rw, r)
}

func (h *setMetricHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
}

func NewSetMetricHandler(serv service.Service, memStore *models.MemStorage) *setMetricHandler {
	return &setMetricHandler{serv: serv, memStore: memStore}
}

func MethodNotAllowedHandler(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusNotFound)
	Body := "Method not allowed!\n"
	fmt.Fprintf(rw, "%s", Body)
}

func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

func (h *setMetricHandler) SetMetricHandler(w http.ResponseWriter, r *http.Request) {

	value, status := isValid(r)
	if !value {
		w.WriteHeader(status)
		return
	} else {
		addMetricToMemStore(h.memStore)
		w.WriteHeader(status)
		Body := "OK\n"
		fmt.Fprintf(w, "%s", Body)
		return
	}
}

func (h *setMetricHandler) GetMetricHandler(w http.ResponseWriter, r *http.Request) {

	/*use in the future*/

	/*s := &service.MetricService{}
	handler := NewGetMetricHandler(s)
	temp, err := handler.serv.GetMetric("any")
	fmt.Println(temp, err)*/
}

func addMetricToMemStore(store *models.MemStorage) {
	if vMetric.mtype == "gauge" {
		store.AddGauge(&models.Gauge{Name: vMetric.mname, Value: vMetric.mvalue_float})

	} else if vMetric.mtype == "counter" {
		store.AddCounter(&models.Counter{Name: vMetric.mname, Value: vMetric.mvalue_int})
	}
}

func isValid(r *http.Request) (bool, int) {

	if !isValidContentType(r.Header.Get("Content-Type")) {
		return false, http.StatusBadRequest
	}

	if !isMethodPost(r.Method) {
		return false, http.StatusMethodNotAllowed
	}

	vMetric.mtype = mux.Vars(r)["metric_type"]
	vMetric.mname = mux.Vars(r)["metric_name"]
	vMetric.mvalue = mux.Vars(r)["metric_value"]

	if !isValidMetricName(vMetric.mname) {
		return false, http.StatusNotFound
	}

	if !isValidMetricType(vMetric.mtype) || !isValidMeticValue(vMetric.mvalue, vMetric.mtype) {
		return false, http.StatusBadRequest
	}

	return true, http.StatusOK
}

func isMethodPost(method string) bool {
	if method == http.MethodPost {
		return true
	} else {
		return false
	}
}

func isValidContentType(contentType string) bool {

	var pattern string = acceptedContentType

	res, err := MatchString(pattern, contentType)
	if err == nil && res == true {
		return true
	} else {
		return false
	}

}

func MatchString(pattern string, s string) (matched bool, err error) {

	re, err := regexp.Compile(pattern)
	if err == nil {
		return re.MatchString(s), nil
	} else {
		return false, err
	}

}

func isValidMetricType(metricType string) bool {

	var pattern string = "^" + metrics + "$"

	res, err := MatchString(pattern, metricType)
	if err == nil && res == true {
		return true
	} else {
		return false
	}

}

func isValidMetricName(metricName string) bool {

	var pattern string = "^[a-zA-Z/ ]{1,20}$"

	res, err := MatchString(pattern, metricName)
	if err == nil && res == true {
		return true
	} else {
		return false
	}
}

func isValidMeticValue(metricValue string, metricType string) bool {

	if metricType == "gauge" {

		if value, err := strconv.ParseFloat(metricValue, 64); err == nil {
			vMetric.mvalue_float = value
			return true
		} else {
			return false
		}

	} else if metricType == "counter" {

		if value, err := strconv.ParseInt(metricValue, 10, 64); err == nil {
			vMetric.mvalue_int = value
			return true
		} else {
			return false
		}

	} else {
		return false
	}
}
