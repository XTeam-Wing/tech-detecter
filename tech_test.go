package tech_detecter

import (
	"fmt"
	"log"
	"net/http"
	"testing"
)

func TestName(t *testing.T) {
	resp, err := http.DefaultClient.Get("https://stg-data-in.ads.heytapmobi.com")
	if err != nil {
		log.Fatal(err)
	}
	tech := TechDetecter{}
	err = tech.Init("/Users/wing/PycharmProjects/pythonProject/rules/")
	if err != nil {
		log.Fatal(err)
	}
	result, err := tech.Detect(resp)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)
}
