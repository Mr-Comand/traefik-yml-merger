package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Sources   []string `yaml:"sources"`
	MAX_DEPTH uint     `yaml:"max_depth"`
}

var config Config
var client = &http.Client{
	Timeout: 10 * time.Second,
}

func fetch(url string) map[string]any {
	resp, err := client.Get(url)
	if err != nil {
		log.Printf("failed to fetch %s: %v", url, err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Printf("failed to fetch %s: HTTP %d", url, resp.StatusCode)
		return nil
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("failed to read %s: %v", url, err)
		return nil
	}

	var result map[string]any

	if err := yaml.Unmarshal(data, &result); err != nil {
		log.Printf("failed to parse %s: %v", url, err)
		return nil
	}

	return result
}

func merge(dst, src map[string]any, depth uint) {
	for key, value := range src {
		srcMap, srcOK := value.(map[string]any)
		dstMap, dstOK := dst[key].(map[string]any)

		if depth < config.MAX_DEPTH && srcOK && dstOK {
			merge(dstMap, srcMap, depth+1)
			continue
		}

		// Later files override earlier values.
		dst[key] = value
		// fmt.Printf("merged key %s at depth %d\n", key, depth)
	}
}

func load(sources []string) map[string]any {
	results := make([]map[string]any, len(sources))

	var wg sync.WaitGroup

	for i, url := range sources {
		wg.Add(1)

		go func(i int, url string) {
			defer wg.Done()
			results[i] = fetch(url)
		}(i, url)
	}

	wg.Wait()

	merged := make(map[string]any)

	// Merge in config order even though fetching was parallel.
	for _, result := range results {
		if result != nil {
			merge(merged, result, 0)
		}
	}

	return merged
}

func main() {
	configData, err := os.ReadFile("config.yml")
	if err != nil {
		log.Fatal(err)
	}

	if err := yaml.Unmarshal(configData, &config); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/merged.yml", func(w http.ResponseWriter, r *http.Request) {
		merged := load(config.Sources)

		data, err := yaml.Marshal(merged)
		if err != nil {
			http.Error(w, "failed to generate YAML", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/yaml")
		w.Write(data)
	})
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK\n"))
	})
	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
