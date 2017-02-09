package cb

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/couchbase/go-couchbase"
)

var (
	cbb                       *couchbase.Bucket
	n1ql                      *http.Client
	queryURL, scanConsistency string
)

// InitDatabase initializes Couchbase Server client and HTTP client for N1QL
// queries.
func InitDatabase(hostname string, consistency string) error {
	baseURL := fmt.Sprintf("http://%s:8091/", hostname)

	c, err := couchbase.ConnectWithAuthCreds(baseURL, "default", "")
	if err != nil {
		return err
	}

	pool, err := c.GetPool("default")
	if err != nil {
		return err
	}

	cbb, err = pool.GetBucket("default")
	if err != nil {
		return err
	}

	t := &http.Transport{MaxIdleConnsPerHost: 10240}
	n1ql = &http.Client{Transport: t}

	queryURL = fmt.Sprintf("http://%s:8093/query", hostname)

	scanConsistency = consistency

	return nil
}

// Insert adds new documents to Couchbase Server bucket using SET operation.
func Insert(_ int64, key string, value interface{}) error {
	return cbb.Set(key, 0, value)
}

type n1qlQuery struct {
	Prepared        string        `json:"prepared"`
	Args            []interface{} `json:"args"`
	ScanConsistency string        `json:"scan_consistency"`
}

func executeQuery(q *n1qlQuery) error {
	b, err := json.Marshal(q)
	if err != nil {
		return err
	}
	j := bytes.NewReader(b)

	resp, err := n1ql.Post(queryURL, "application/json", j)
	defer resp.Body.Close()

	if err != nil {
		return err
	}

	if _, err = ioutil.ReadAll(resp.Body); err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New("bad response")
	}

	return nil
}

// Query finds matching documents using prepared N1QL statements.
func Query(_ int64, field string, arg interface{}) error {
	// FIXME: support different queries and multiple arguments
	q := n1qlQuery{"by_" + field, []interface{}{arg}, scanConsistency}

	return executeQuery(&q)
}