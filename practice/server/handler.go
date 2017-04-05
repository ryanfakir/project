package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"strings"

	"regexp"

	"github.com/julienschmidt/httprouter"
)

// NumberOfWorker for goroutines
const NumberOfWorker = 4

// Handler is a collection of all the service handlers.
type Handler struct {
	*httprouter.Router
	Client *Client
	Logger *log.Logger
}

// Result for single query
type Result struct {
	Hostname  string    `json:"hostname"`
	Ports     []int     `json:"ports"`
	Err       error     `json:"err,omitempty"`
	QueryTime time.Time `json:"querytime, omitempty"`
}

// Results is list of Result
type Results struct {
	Prev []Result `json:"prev"`
	Now  []Result `json:"now"`
}

// RequestData for parsing from body
type RequestData struct {
	Querykeys []string `json:"querykeys"`
}

// NewHandler returns a new instance of Handler.
func NewHandler() *Handler {
	h := &Handler{
		Router: httprouter.New(),
		Client: NewClient(),
		Logger: log.New(os.Stderr, "", log.LstdFlags),
	}
	// h.GET("/", defaultHandler)
	h.POST("/", h.postHandler)
	// h.ServeFiles("/static/*filepath", http.Dir("static"))
	return h
}

// func defaultHandler(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
// 	index := template.Must(template.ParseFiles("static/index.html"))
// 	index.Execute(rw, nil)
// }

func (h *Handler) postHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var data RequestData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		Error(w, err, http.StatusBadRequest, h.Logger)
		return
	}
	if !verifyRequestData(data) {
		Error(w, nil, http.StatusBadRequest, h.Logger)
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	in := startChan(ctx, data.Querykeys)
	var wg sync.WaitGroup
	wg.Add(NumberOfWorker)
	resultChan := make(chan Result)
	for i := 0; i < NumberOfWorker; i++ {
		go func() {
			for el := range in {
				h.Logger.Printf("starting process... %v", el)
				start := time.Now()
				res := h.Client.CallService(ctx, el)
				elapsed := time.Since(start)
				h.Logger.Printf("%v cost time: %v", el, elapsed)
				select {
				case <-ctx.Done():
					return
				case resultChan <- res:
				}
			}
			wg.Done()
		}()
	}
	go func() {
		wg.Wait()
		close(resultChan)
	}()
	var d Results
	for res := range resultChan {
		if res.Err != nil {
			Error(w, res.Err, http.StatusInternalServerError, h.Logger)
			return
		}

		d.Now = append(d.Now, res)
		if h.Client.history[res.Hostname] != nil {
			d.Prev = append(d.Prev, h.Client.history[res.Hostname].res)
		}

	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	encodeJSON(w, d, h.Logger)
}

func startChan(ctx context.Context, querys []string) <-chan string {
	in := make(chan string)
	go func() {
		defer close(in)
		for _, el := range querys {
			select {
			case <-ctx.Done():
				return
			case in <- el:
			}
		}
	}()
	return in
}

func verifyRequestData(data RequestData) bool {
	for _, v := range data.Querykeys {
		if !verifyIP(v) && !verifyHostName(v) {
			return false
		}
	}
	return true
}

func verifyIP(input string) bool {
	if len(input) < 7 {
		return false
	}
	if input[0] == '.' || input[len(input)-1] == '.' {
		return false
	}
	tokens := strings.Split(input, ".")
	if len(tokens) != 4 {
		return false
	}
	for _, v := range tokens {
		if len(v) > 1 && v[0] == '0' {
			return false
		}
		output, err := strconv.Atoi(v)
		if err != nil {
			return false
		}
		if output < 0 || output > 255 {
			return false
		}
	}

	return true
}

func verifyHostName(input string) bool {
	// Check HostName validation
	r, err := regexp.Compile("^(([a-zA-Z]|[a-zA-Z][a-zA-Z\\-]*[a-zA-Z])\\.)*([A-Za-z]|[A-Za-z][A-Za-z\\-]*[A-Za-z])$")
	if err != nil {
		return false
	}
	if r.MatchString(input) {
		return true
	}
	return false
}

// encodeJSON encodes v to w in JSON format
func encodeJSON(w http.ResponseWriter, v interface{}, logger *log.Logger) {
	if err := json.NewEncoder(w).Encode(v); err != nil {
		Error(w, err, http.StatusInternalServerError, logger)
	}
}

// Error writes an API error message to the response and logger.
func Error(w http.ResponseWriter, err error, code int, logger *log.Logger) {
	// Log error.
	logger.Printf("http error: %s (code=%d)", err, code)
	http.Error(w, err.Error(), code)
}
