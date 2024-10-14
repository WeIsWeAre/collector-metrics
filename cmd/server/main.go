package main

import (
	"errors"
	"flag"
	"fmt"
	"handlers"
	"log"
	"middleware"
	"models"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"service"
	"time"

	"github.com/gorilla/mux"
)

const validateAdderPattern string = "^[a-zA-Z/ ]{1,100}:[0-9]{1,10}$"

type Options struct {
	PORT string
}

var memStorage models.MemStorage
var s *http.Server
var serv *service.MetricService
var options Options

func main() {

	err := initialization()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	go func() {
		log.Println("Listening to", ":"+options.PORT)
		err := s.ListenAndServe()
		if err != nil {
			log.Printf("Error starting server: %s\n", err)
			return
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	sig := <-sigs
	log.Println("Quitting after signal:", sig)
	time.Sleep(2 * time.Second)
	s.Shutdown(nil)
}

func initialization() error {

	err := parseFlags()
	if err != nil {
		return err
	}
	err = getENV()
	if err != nil {
		return err
	}
	mux := mux.NewRouter()

	memStorage.Init()
	handler := handlers.NewSetMetricHandler(serv, &memStorage)

	postMux := mux.Methods(http.MethodPost).Subrouter()
	postMux.HandleFunc("/update/{metric_type}/{metric_name}/{metric_value}", handler.SetMetricHandler)
	postMux.Use(middleware.MiddlewareSetMetric)

	getMux := mux.Methods(http.MethodGet).Subrouter()
	getMux.HandleFunc("/value/{metric_type}/{metric_name}", handler.GetMetricHandler)
	getMux.Use(middleware.MiddlewareSetMetric)

	notAllowed := handlers.NotAllowedHandler{}
	mux.MethodNotAllowedHandler = notAllowed

	mux.NotFoundHandler = http.HandlerFunc(handler.DefaultHandler)

	s = &http.Server{
		Addr:         options.PORT,
		Handler:      mux,
		ErrorLog:     nil,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	return nil
}

func parseFlags() error {
	flag.StringVar(&options.PORT, "a", "localhost:8080", "Port to listen on.")
	flag.Parse()

	if !isValidAddr(options.PORT) {
		return errors.New("Addr is not valid")
	}
	return nil
}

func getENV() error {
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		if isValidAddr(envRunAddr) {
			options.PORT = envRunAddr
		} else {
			return errors.New("Addr is not valid")
		}
	}
	return nil
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
