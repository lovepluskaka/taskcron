package taskcron

import (
	"log"
	"net/http"
)

func get(url string) bool {
	result := false

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Get url %s error is %s\n", url, err.Error())
		return result
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		result = true
	}

	return result
}
