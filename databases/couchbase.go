package databases

import (
	"log"

	"github.com/couchbaselabs/go-couchbase"
)

type Couchbase struct {
	Bucket *couchbase.Bucket
}

func (cb *Couchbase) Init(config Config) {
	bucket, err := couchbase.GetBucket(config.Addresses[0], config.Name, config.Table)
	if err != nil {
		log.Fatal(err)
	}
	cb.Bucket = bucket
}

func (cb *Couchbase) Shutdown() {
	cb.Bucket.Close()
}

func (cb *Couchbase) Create(key string, value map[string]interface{}) error {
	err := cb.Bucket.Set(key, 0, value)
	return err
}

func (cb *Couchbase) Read(key string) error {
	result := map[string]interface{}{}
	err := cb.Bucket.Get(key, &result)
	return err
}

func (cb *Couchbase) Update(key string, value map[string]interface{}) error {
	err := cb.Bucket.Set(key, 0, value)
	return err
}

func (cb *Couchbase) Delete(key string) error {
	err := cb.Bucket.Delete(key)
	return err
}

func (cb *Couchbase) Query(fieldName, fieldValue string, limit int) error {
	return nil //TODO: implement
}
