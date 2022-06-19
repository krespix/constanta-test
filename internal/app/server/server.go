package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	reqs "github.com/krespix/constanta-test/internal/services/request"
)

type Server interface {
	Start(ctx context.Context) error
}
type server struct {
	reqService reqs.Service

	addr        string
	maxRequests int

	mux *http.ServeMux

	//requests in process counter
	mu            *sync.Mutex
	reqsInProcess int
}

func (s *server) Start(ctx context.Context) error {
	s.configure(ctx)
	return http.ListenAndServe(s.addr, s.limit(s.mux))
}

func (s *server) configure(ctx context.Context) {
	s.mux.HandleFunc("/v1/collect-data", s.collectData(ctx))
}

func (s *server) collectData(ctx context.Context) http.HandlerFunc {
	type request struct {
		URLs []string `json:"urls"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			s.incReqsCounter()
			defer s.decReqsCounter()
			req := &request{}
			if err := json.NewDecoder(r.Body).Decode(req); err != nil {
				s.error(w, r, http.StatusBadRequest, err)
				return
			}
			if len(req.URLs) > 20 {
				s.error(w, r, http.StatusBadRequest, fmt.Errorf("length of urls array must be less or equal 20"))
				return
			}
			res, err := s.reqService.CollectData(ctx, req.URLs)
			if err != nil {
				s.error(w, r, http.StatusBadRequest, err)
				return
			}
			s.respond(w, r, http.StatusOK, res)
		default:
			s.error(w, r, http.StatusBadRequest, fmt.Errorf("only supported post method by this uri"))
			return
		}
	}
}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
		}
	}
}

func (s *server) limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if s.reqsInProcess >= s.maxRequests {
			http.Error(w, http.StatusText(429), http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *server) decReqsCounter() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.reqsInProcess--
}

func (s *server) incReqsCounter() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.reqsInProcess++
}

func New(addr string, reqService reqs.Service, maxRequests int) *server {
	return &server{
		addr:        addr,
		reqService:  reqService,
		mux:         http.NewServeMux(),
		maxRequests: maxRequests,
		mu:          &sync.Mutex{},
	}
}
