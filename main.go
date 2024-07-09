package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/demyanovs/robotstxt"
	_ "golang.org/x/lint"

	"github.com/demyanovs/urlcrawler/queue"
	"github.com/demyanovs/urlcrawler/report"
)

const (
	fileNameDefault = "result"
)

const (
	outputCSV  = "csv"
	outputJSON = "json"
)

var supportedOutputs = []string{outputCSV, outputJSON}

func main() {
	startURL := flag.String("u", "", "Start url (required)")
	output := flag.String("output", outputCSV, "Output format (csv, json)")
	outputFile := flag.String("output-file", "", "File path to save r")
	delay := flag.Int("delay", 1000, "Delay between requests in milliseconds")
	depth := flag.Int("depth", 0, "Depth of the crawl (0 - infinite")
	limitURLs := flag.Int("limit", 0, "Limit of URLs to crawl (0 - unlimited")
	reqTimeout := flag.Int("timeout", 5000, "Request timeout in milliseconds")
	bulkSize := flag.Int("bulk-size", 30, "Bulk size for saving to the file")
	queueLen := flag.Int("q-len", 50, "Queue length")
	quietMode := flag.Bool("q", false, "Quiet mode (no logs")
	ignoreRobotsTXT := flag.Bool("ignore-robots", false, "Ignore crawl-delay and disallowed URLs from robots.txt")

	flag.Parse()

	if *startURL == "" {
		log.Fatal("url flag is required")
	}

	if *output != outputCSV && *output != outputJSON {
		log.Fatalf("unsupported output format: %s. Supported formats: %v", *output, supportedOutputs)
	}

	logger := log.New(log.Writer(), "", log.Ldate|log.Ltime)

	r, reportFile := getReport(*output, *outputFile)

	q, err := queue.New(
		queue.ConfigType{
			QueueLen:   *queueLen,
			LimitURLs:  *limitURLs,
			ReqTimeout: time.Duration(*reqTimeout) * time.Millisecond,
			Delay:      time.Duration(*delay) * time.Millisecond,
			BulkSize:   *bulkSize,
			Quiet:      *quietMode,
			Depth:      *depth,
		},
		*startURL,
		r,
		logger,
		nil,
	)

	if *ignoreRobotsTXT == true {
		if *quietMode == false {
			logger.Println("ignoring robots.txt")
		}
	} else {
		if *quietMode == false {
			log.Println("parsing robots.txt")
		}
		robots, err := getRobotsTXT(*startURL)
		if err != nil {
			log.Fatal(err)
		}

		q.RobotsData = robots

		crawlDelay, err := robots.GetCrawlDelay("*")
		if err != nil {
			log.Fatal(err)
		}

		if crawlDelay != nil {
			q.Config.Delay = time.Duration(*crawlDelay) * time.Second
			if *quietMode == false {
				logger.Printf("found crawl-delay in robots.txt: %ds. Ignoring delay from the config\n", *crawlDelay)
			}
		}
	}

	if err != nil {
		log.Fatal(err)
	}

	if *quietMode == false {
		printConfig(q, *output, reportFile, *ignoreRobotsTXT, logger)
	}

	q.Start()
}

func getRobotsTXT(startURL string) (*robotstxt.RobotsData, error) {
	parsedURL, err := url.Parse(startURL)
	if err != nil || parsedURL == nil {
		return nil, err
	}

	robotsPath := fmt.Sprintf("%s://%s/robots.txt", parsedURL.Scheme, parsedURL.Host)
	resp, err := http.Get(robotsPath)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	robots, err := robotstxt.FromResponse(resp)
	if err != nil {
		return nil, err
	}

	return robots, nil
}

func getReport(output string, outputFile string) (queue.Reporter, string) {
	var r queue.Reporter
	if output == outputJSON {
		if outputFile == "" {
			outputFile = fmt.Sprintf("%s.%s", fileNameDefault, outputJSON)
		}
		r = report.NewJSONReport(outputFile)
	} else {
		if outputFile == "" {
			outputFile = fmt.Sprintf("%s.%s", fileNameDefault, outputCSV)
		}
		r = report.NewCSVReport(outputFile)
	}

	return r, outputFile
}

func printConfig(queue *queue.Queue, output string, outputFile string, ignoreRobotsTXT bool, logger *log.Logger) {
	logger.Printf(
		"Starting crawling, "+
			"delay: %dms, "+
			"depth: %d, "+
			"limit: %d, "+
			"reqTimeout: %dms, "+
			"bulk-size: %d, "+
			"output: %s, "+
			"output-file: %s, "+
			"ignore-robots: %t "+
			"\n",
		queue.Config.Delay/time.Millisecond,
		queue.Config.Depth,
		queue.Config.LimitURLs,
		queue.Config.ReqTimeout/time.Millisecond,
		queue.Config.BulkSize,
		output,
		outputFile,
		ignoreRobotsTXT,
	)
}
