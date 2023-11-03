package es

import (
	"context"
	"github.com/olivere/elastic/v7"
	"testing"
)

type Person struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Married bool   `json:"married"`
}

func TestES(t *testing.T) {
	esClient, err := elastic.NewClient(elastic.SetURL("http://192.168.30.130:9200"))
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log("connect to elastic successfully")

	p1 := &Person{
		Name:    "jojo",
		Age:     16,
		Married: false,
	}
	esClient.Index().Index("people").BodyJson(p1).Do(context.Background())
	if err != nil {
		panic(err)
	}
}
