package queue

import (
	"context"
	"fmt"
	"github.com/demyanovs/urlcrawler/parser"
	"github.com/demyanovs/urlcrawler/store"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// Queue represents a queue for processing URLs.
type Queue struct {
	Config          ConfigType
	startURL        *url.URL
	report          Reporter
	RobotsData      RobotsData
	parser          parser.Parser
	logger          Logger
	startedAt       time.Time
	sURLsDone       URLStore
	sURLsToDo       URLStore
	sURLsInProgress URLStore
	sURLsToSave     URLStore
}

// ConfigType represents a configuration for the queue.
type ConfigType struct {
	QueueLen   int
	LimitURLs  int
	BulkSize   int
	ReqTimeout time.Duration
	Delay      time.Duration
	Depth      int
	Quiet      bool
}

// URLStore represents a store for URLs.
type URLStore interface {
	Add(key string, value any)
	Get(key string) (any, error)
	Delete(key string)
	List() map[string]any
	Values() []any
	Keys() []string
	Clear()
	Len() int
}

// RobotsData represents a robots.txt file.
type RobotsData interface {
	IsAllowed(userAgent string, URL string) bool
	CrawlDelay(userAgent string) (*int, error)
}

// Reporter represents a reporter.
type Reporter interface {
	SaveBulk(records []parser.PageData) error
}

// Logger represents a logger.
type Logger interface {
	Println(v ...any)
}

// New creates a new queue.
func New(
	config ConfigType,
	startURL string,
	report Reporter,
	logger Logger,
	robotsData RobotsData,
) (*Queue, error) {
	parsedURL, err := url.Parse(startURL)
	if err != nil || parsedURL == nil {
		return nil, err
	}

	sURLsToDo := store.New()
	sURLsToDo.Add(startURL, 0)

	return &Queue{
		Config:          config,
		startURL:        parsedURL,
		report:          report,
		RobotsData:      robotsData,
		parser:          parser.New(),
		logger:          logger,
		sURLsDone:       store.New(),
		sURLsToDo:       sURLsToDo,
		sURLsInProgress: store.New(),
		sURLsToSave:     store.New(),
	}, nil
}

// Start starts the queue.
func (q *Queue) Start() {
	q.startedAt = time.Now()
	active := true
	var wg sync.WaitGroup

	queue := make(chan struct{}, q.Config.QueueLen)

	for ok := true; ok; ok = active {
		URLs := q.sURLsToDo.Keys()
		wg.Add(len(URLs))

		for w := 0; w < len(URLs); w++ {
			if q.sURLsDone.Len() >= q.Config.LimitURLs && q.Config.LimitURLs > 0 {
				q.log(fmt.Sprintf("reached max URLs limit of %d", q.Config.LimitURLs))
				active = false
				q.Stop()
				break
			}

			URL := URLs[w]
			v, err := q.sURLsToDo.Get(URL)
			if err != nil {
				log.Fatalln(err)
			}
			q.process(queue, &wg, URL, v.(int))

			time.Sleep(q.Config.Delay)
		}

		if q.sURLsToDo.Len() == 0 && q.sURLsInProgress.Len() == 0 {
			q.Stop()
			active = false
		}
	}
}

// Stop stops the queue.
func (q *Queue) Stop() {
	err := q.saveResults()
	if err != nil {
		log.Println(err)
		return
	}

	elapsed := time.Since(q.startedAt)
	q.log(fmt.Sprintf("crawling completed. %d of %d URLs processed in %s", q.sURLsDone.Len(), q.sURLsToDo.Len(), elapsed.Round(time.Second)))
}

func (q *Queue) process(queue chan struct{}, wg *sync.WaitGroup, URL string, depth int) {
	queue <- struct{}{}

	go func() {
		defer wg.Done()

		q.sURLsToDo.Delete(URL)
		q.sURLsInProgress.Add(URL, depth)

		q.log(fmt.Sprintf("processing: %s, (found: %d)", URL, q.sURLsToDo.Len()))

		ctx, cancel := context.WithTimeout(context.Background(), q.Config.ReqTimeout)
		defer cancel()

		var pageData parser.PageData
		var linksOnPage []string

		// Start processing
		resp, err := q.readURL(ctx, URL)
		if err != nil {
			pageData = parser.PageData{
				URL: URL,
			}
			fmt.Println(fmt.Errorf("can't send request to url %s. Error: %s", URL, err))
		} else {
			pageData, linksOnPage, err = q.parser.ParseResponse(resp)
			if err != nil {
				log.Println(err)
			}
		}

		q.sURLsDone.Add(URL, pageData)
		q.sURLsToSave.Add(URL, pageData)

		if len(linksOnPage) > 0 && (q.Config.Depth == 0 || depth <= q.Config.Depth) {
			q.addSURLsToDo(linksOnPage, depth)
		}

		q.sURLsInProgress.Delete(URL)

		if q.sURLsToSave.Len() >= q.Config.BulkSize {
			q.log(fmt.Sprintf("store is full: %d", q.sURLsToSave.Len()))

			err = q.saveResults()
			if err != nil {
				log.Println(err)
				return
			}
		}

		<-queue
	}()
}

func (q *Queue) toPagesData(data []any) parser.PagesData {
	var pagesData parser.PagesData
	for _, d := range data {
		switch d.(type) {
		case parser.PageData:
			pagesData = append(pagesData, d.(parser.PageData))
		}
	}

	return pagesData
}

func (q *Queue) saveResults() error {
	q.log("saving to the file...")
	if q.report == nil || q.sURLsToSave.Len() == 0 {
		return nil
	}

	vs := q.sURLsToSave.Values()
	pagesData := q.toPagesData(vs)
	err := q.report.SaveBulk(pagesData)
	if err != nil {
		return err
	}

	q.sURLsToSave.Clear()

	q.log(fmt.Sprintf("saved. Done %d, todo %d", q.sURLsDone.Len(), q.sURLsToDo.Len()))

	return nil
}

func (q *Queue) addSURLsToDo(linksOnPage []string, depth int) {
	for _, l := range linksOnPage {
		// Check if the URL is allowed in robots.txt
		isAllowed := q.RobotsData == nil || q.RobotsData.IsAllowed("*", l)
		if isAllowed == false {
			continue
		}

		fullURL := fmt.Sprintf("%s://%s/%s", q.startURL.Scheme, q.startURL.Host, l)
		_, err := q.sURLsDone.Get(fullURL)
		if err != nil {
			// Do not add the URL if depth is greater than the limit
			nextDepth := depth + 1
			if q.Config.Depth > 0 && nextDepth > q.Config.Depth {
				continue
			}

			q.sURLsToDo.Add(fullURL, nextDepth)
		}
	}
}

func (q *Queue) readURL(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (q *Queue) log(message string) {
	if q.Config.Quiet == true {
		return
	}

	q.logger.Println(message)
}
