package workloads

import (
	"crypto/md5"
	"encoding/hex"
	"sync"

	"github.com/pavel-paulau/blurr/databases"
)

type Config struct {
	Type                    string
	CreatePercentage        int // shorthand "c"
	ReadPercentage          int // shorthand "r"
	UpdatePercentage        int // shorthand "u"
	DeletePercentage        int // shorthand "d"
	Records                 int64
	Operations              int64
	ValueSize               int
	Workers                 int
	TargetThroughput        int
	HotDataPercentage       int64
	HotSpotAccessPercentage int
}

type Workload interface {
	SetImplementation(i Workload)

	GenerateNewKey(currentRecords int64) string

	GenerateExistingKey(currentRecords int64) string

	GenerateKeyForRemoval() string

	GenerateValue(key string, size int) map[string]interface{}

	PrepareBatch() []string

	DoBatch(database databases.Database, state *State)

	RunWorkload(database databases.Database, state *State, wg *sync.WaitGroup)
}

func Hash(inString string) string {
	h := md5.New()
	h.Write([]byte(inString))
	return hex.EncodeToString(h.Sum(nil))
}

func RandString(key string, expectedLength int) string {
	var randString string
	if expectedLength > 64 {
		baseString := RandString(key, expectedLength/2)
		randString = baseString + baseString
	} else {
		randString = (Hash(key) + Hash(key[:len(key)-1]))[:expectedLength]
	}
	return randString
}
