package workloads

import (
	"sync"

	"github.com/pavel-paulau/blurr/databases"
)

type Config struct {
	Type                    string
	CreatePercentage        int
	ReadPercentage          int
	UpdatePercentage        int
	DeletePercentage        int
	QueryPercentage         int
	Records                 int64
	Operations              int64
	ValueSize               int
	Workers                 int
	Throughput              int
	HotDataPercentage       int64
	HotSpotAccessPercentage int
	Views                   []string
}

type Workload interface {
	SetImplementation(i Workload)

	GenerateNewKey(currentRecords int64) string

	GenerateExistingKey(currentRecords int64) string

	GenerateKeyForRemoval() string

	GenerateValue(key string, size int) map[string]interface{}

	GenerateQueryArgs(key string) []interface{}

	PrepareBatch() []string

	PrepareSeq(size int64) chan string

	DoBatch(database databases.Database, state *State, seq chan string)

	RunWorkload(database databases.Database, state *State, wg *sync.WaitGroup)
}
