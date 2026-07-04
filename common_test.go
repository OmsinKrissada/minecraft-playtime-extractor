package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetWorldRunTime(t *testing.T) {
	config.WorldDir = "test_server_dir/world"
	wt := getWorldRunTime()
	assert.EqualValues(t, 119518758, wt)
}

func TestGetPlaytime(t *testing.T) {
	config.WorldDir = "test_server_dir/world"
	pt := getAllPlaytime()
	expectedResult := map[string]int{"0e3715e0-e902-4738-ba8a-f58f70c23e64": 49989, "2a2a2c18-c1e8-4ca6-b50e-2c6c7d509659": 11063603, "6bce8c3c-7c91-4993-865f-e9d28d557d1f": 203690, "d89b17bc-00d8-4b33-ba3a-28f18a58cb51": 21727052}
	assert.EqualValues(t, expectedResult, pt)
}
