package main

import (
	"my5G-RANTester/config"
	"my5G-RANTester/internal/templates"
	"os"
	"strconv"
	"time"

	"github.com/davecgh/go-spew/spew"
	log "github.com/sirupsen/logrus"
)

const version = "1.0.2"

func init() {

	cfg, err := config.GetConfig()
	if err != nil {
		//return nil
		log.Fatal("Error in get configuration")
	}

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	if cfg.Logs.Level == 0 {
		log.SetLevel(log.InfoLevel)
	} else {
		log.SetLevel(log.Level(cfg.Logs.Level))
	}

	spew.Config.Indent = "\t"

	log.Info("my5G-RANTester version " + version)
}

func main() {
	log.Warn("Reading Environment Variables...")

	num_requests_start, err := strconv.Atoi(os.Getenv("NUM_REQUESTS_START"))
	if err != nil {
		log.Fatal("Invalid NUM_REQUESTS_START: ", os.Getenv("NUM_REQUESTS_START"),
			" \nAn environment variable named NUM_REQUESTS_START must be set as an integer ",
			"indicating the number of requests per seconds to start with")
	}
	log.Info("    Found NUM_REQUESTS_START: ", os.Getenv("NUM_REQUESTS_START"))

	requests_increment, err := strconv.Atoi(os.Getenv("REQUESTS_INCREMENT"))
	if err != nil {
		log.Fatal("Invalid REQUESTS_INCREMENT: ", os.Getenv("REQUESTS_INCREMENT"),
			" \nAn environment variable named REQUESTS_INCREMENT must be set as an integer ",
			"indicating the number of requests increased after each iteration")
	}
	log.Info("    Found REQUESTS_INCREMENT: ", os.Getenv("REQUESTS_INCREMENT"))

	num_iterations, err := strconv.Atoi(os.Getenv("NUM_ITERATIONS"))
	if err != nil {
		log.Fatal("Invalid NUM_ITERATIONS: ", os.Getenv("NUM_ITERATIONS"),
			" \nAn environment variable named NUM_ITERATIONS must be set as an integer ",
			"indicating the number of iterations to execute")
	}
	log.Info("    Found NUM_ITERATIONS: ", os.Getenv("NUM_ITERATIONS"))

	iteration_duration, err := time.ParseDuration(os.Getenv("ITERATION_DURATION"))
	if err != nil {
		log.Fatal("Invalid ITERATION_DURATION: ", os.Getenv("ITERATION_DURATION"),
			" \nAn environment variable named ITERATION_DURATION must be set as an duration ",
			"indicating the number of iterations to execute.\n",
			"A duration string is a sequence of decimal numbers, each with optional fraction ",
			"and a unit suffix, such as \"300ms\", \"1.5h\" or \"2h45m\". Valid time units are ",
			"\"ns\", \"us\" (or \"µs\"), \"ms\", \"s\", \"m\", \"h\".")
	}
	log.Info("    Found ITERATION_DURATION: ", os.Getenv("ITERATION_DURATION"))

	cfg := config.Data

	log.Info("---------------------------------------")
	log.Warn("[TESTER] Starting Iterative AMF Load Test")
	log.Warn("[TESTER] Expected Duration: ", time.Duration(iteration_duration.Seconds()*float64(num_iterations)).String())
	log.Info("[TESTER][GNB] gNodeB control interface IP/Port: ", cfg.GNodeB.ControlIF.Ip, "/", cfg.GNodeB.ControlIF.Port)
	log.Info("[TESTER][GNB] gNodeB data interface IP/Port: ", cfg.GNodeB.DataIF.Ip, "/", cfg.GNodeB.DataIF.Port)
	log.Info("[TESTER][AMF] AMF IP/Port: ", cfg.AMF.Ip, "/", cfg.AMF.Port)
	log.Info("---------------------------------------")

	for iteration := 0; iteration < num_iterations; iteration++ {
		log.Warn("[TESTER] -------- Starting Iteration ", strconv.Itoa(iteration), " --------")
		request_rate := num_requests_start + (iteration * requests_increment)
		log.Info("[TESTER] Running ", strconv.Itoa(request_rate), " requests for the AMF per second, during ", iteration_duration.String())
		log.Warn("[TESTER][GNB] Total of AMF Responses in the interval:", templates.TestRqsLoop(request_rate, int(iteration_duration.Seconds())))
	}
}
