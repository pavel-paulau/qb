package qb

import (
	"context"
	"math/rand"
	"sync"
	"time"
)

type insertFn func(workerID int64, key string, value *Doc) error

type queryFn func(workerID int64, payload *QueryPayload) error

func singleLoad(wg *sync.WaitGroup, workerID int64, w *WorkloadSettings) {
	defer wg.Done()

	for payload := range generatePayload(workerID, w) {
		if err := w.IFn(workerID, payload.key, payload.value); err != nil {
			logFatalln(err)
		}
	}
}

// WorkloadSettings incorporates all possible workload settings.
type WorkloadSettings struct {
	NumWorkers, NumDocs, DocSize int64
	InsertPercentage             int
	Time                         time.Duration
	IFn                          insertFn
	QFn                          queryFn
	Hostname                     string
	Consistency                  string
	QueryType                    int
	SSL                          bool
}

// SetQueryType matches string query type to integer value
func (w *WorkloadSettings) SetQueryType(workload string) {
	switch workload {
	case "Q1":
		w.QueryType = q1query
	case "Q2":
		w.QueryType = q2query
	case "Q3":
		w.QueryType = q3query
	case "Q4":
		w.QueryType = q4query
	case "Q5":
		w.QueryType = q5query
	}
}

// Load executes the load phase - insertion of brand new items.
func Load(w *WorkloadSettings) {
	wg := sync.WaitGroup{}

	w.NumDocs /= w.NumWorkers

	for i := int64(0); i < w.NumWorkers; i++ {
		wg.Add(1)
		go singleLoad(&wg, i, w)
	}

	wg.Wait()
}

func singleRun(wg *sync.WaitGroup, workerID int64, w *WorkloadSettings, ctx context.Context) {
	defer wg.Done()

	ch1, ch2 := generateMixedPayload(w, workerID)

	for {
		select {
		case payload := <-ch1:
			if err := w.IFn(workerID, payload.key, payload.value); err != nil {
				logFatalln(err)
			}
		case payload := <-ch2:
			if err := w.QFn(workerID, payload); err != nil {
				logFatalln(err)
			}
		case <-ctx.Done():
			return
		}
	}
}

// const rampUpDelay = 250 * time.Millisecond

var mu sync.RWMutex

// Run executes mixed workloads - a mix of queries and insert operations.
func Run(w *WorkloadSettings) {
	ctx, cancel := context.WithTimeout(context.Background(), w.Time)
	defer cancel()

	wg := sync.WaitGroup{}

	mu = sync.RWMutex{}
	currDocuments = w.NumDocs

	for i := int64(0); i < w.NumWorkers; i++ {
		wg.Add(1)
		go singleRun(&wg, i, w, ctx)

		delay := time.Duration(rand.Int63n(25)) * time.Millisecond
		time.Sleep(delay)
	}

	wg.Wait()
}
