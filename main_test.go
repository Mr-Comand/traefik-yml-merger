package main

import "testing"

func TestMergeDeeplyMergesSourceMaps(t *testing.T) {
	config.MAX_DEPTH = 10

	first := map[string]any{
		"http": map[string]any{
			"routers": map[string]any{
				"router1": map[string]any{
					"rule":    "Host(`one.example.com`)",
					"service": "svc-one",
				},
			},
		},
	}

	second := map[string]any{
		"http": map[string]any{
			"routers": map[string]any{
				"router2": map[string]any{
					"rule":    "Host(`two.example.com`)",
					"service": "svc-two",
				},
			},
		},
	}

	merged := make(map[string]any)
	merge(merged, first, 0)
	merge(merged, second, 0)

	httpMap, ok := merged["http"].(map[string]any)
	if !ok {
		t.Fatal("expected http map to exist")
	}

	routers, ok := httpMap["routers"].(map[string]any)
	if !ok {
		t.Fatal("expected routers map to exist")
	}

	if _, ok := routers["router1"]; !ok {
		t.Fatal("expected router1 from first source to remain after merge")
	}
	if _, ok := routers["router2"]; !ok {
		t.Fatal("expected router2 from second source to be added after merge")
	}
}

func TestMergeLaterValuesOverrideEarlierValues(t *testing.T) {
	config.MAX_DEPTH = 10

	first := map[string]any{
		"middlewares": map[string]any{
			"auth": map[string]any{"header": "old"},
		},
	}
	second := map[string]any{
		"middlewares": map[string]any{
			"auth": map[string]any{"header": "new"},
		},
	}

	merged := make(map[string]any)
	merge(merged, first, 0)
	merge(merged, second, 0)

	middlewares, ok := merged["middlewares"].(map[string]any)
	if !ok {
		t.Fatal("expected middlewares map to exist")
	}

	auth, ok := middlewares["auth"].(map[string]any)
	if !ok {
		t.Fatal("expected auth map to exist")
	}

	if got := auth["header"]; got != "new" {
		t.Fatalf("expected later config to override earlier value, got %v", got)
	}
}
