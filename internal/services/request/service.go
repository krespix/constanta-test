package request

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	netURL "net/url"
	"strings"
	"sync"
	"time"
)

type Service interface {
	CollectData(ctx context.Context, urls []string) ([]*Response, error)
}

type service struct {
	httpClient *http.Client

	maxWorkers int
}

func (s *service) CollectData(ctx context.Context, urls []string) ([]*Response, error) {
	err := validateURLList(urls)
	if err != nil {
		return nil, err
	}

	responseList := make([]*Response, 0, len(urls))
	chunk := make([]string, 0, s.maxWorkers)
	resChan := make(chan *result, s.maxWorkers)
	for i := 0; i < len(urls)/s.maxWorkers+1; i++ {
		//divide urls to chunks by s.maxWorkers
		wg := &sync.WaitGroup{}
		start := i * s.maxWorkers
		end := (i + 1) * s.maxWorkers
		if i == len(urls)/s.maxWorkers {
			end = i*s.maxWorkers + len(urls)%s.maxWorkers
		}
		chunk = urls[start:end]

		//run s.maxWorkers to do requests
		for _, url := range chunk {
			wg.Add(1)
			s.doRequest(url, resChan, wg)
		}
		//waiting until processing of requests end
		wg.Wait()
		for j := 0; j < len(chunk); j++ {
			res := <-resChan
			//end of processing whole req if one err
			if res.Err != nil {
				return nil, res.Err
			}
			responseList = append(responseList, res.Response)
		}
	}
	close(resChan)
	return responseList, nil
}

// async request by url
// return result in resChan
func (s *service) doRequest(url string, resChan chan<- *result, wg *sync.WaitGroup) {
	go func() {
		defer wg.Done()
		resp, err := s.httpClient.Get(url)
		if err != nil {
			resChan <- &result{
				Err: err,
			}
			return
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			resChan <- &result{
				Err: err,
			}
			return
		}
		resChan <- &result{
			Response: &Response{
				Resource: url,
				Data:     string(body),
			},
		}
	}()
}

func validateURLList(urls []string) error {
	errList := make([]string, 0)
	for _, url := range urls {
		err := validateURL(url)
		if err != nil {
			errList = append(errList, fmt.Sprintf("url %s not valid: %v", url, err))
		}
	}
	if len(errList) > 0 {
		return fmt.Errorf(strings.Join(errList[:], "; "))
	}
	return nil
}

//validates url - return true if url valid, false if not
func validateURL(url string) error {
	_, err := netURL.ParseRequestURI(url)
	return err
}

// New creates new instance of request service
// timeout - timeout for http client
// maxWorkers - max outgoing requests per 1 incoming request
func New(timeout time.Duration, maxWorkers int) *service {
	return &service{
		httpClient: &http.Client{
			Timeout: timeout,
		},
		maxWorkers: maxWorkers,
	}
}
