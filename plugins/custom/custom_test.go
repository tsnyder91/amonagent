package custom

import (
	"path"
	"runtime"
	"testing"

	"github.com/amonapp/amonagent/internal/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCustomCollect(t *testing.T) {

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("testdata directory not found")
	}

	var pythonScript = path.Join("python ", path.Dir(filename), "testdata", "connections.py")

	config := Config{}
	configLine := util.Command{Name: "connections", Command: pythonScript}

	config.Commands = append(config.Commands, configLine)

	c := Custom{}
	c.Config = config

	result, err := c.Collect()
	require.NoError(t, err)

	fields := map[string]interface{}{
		"connections.active": float64(100),
		"connections.error":  float64(500),
	}

	expectedResults := make(PerformanceStructBlock, 0)
	p := PerformanceStruct{Gauges: fields}
	expectedResults["connections"] = p

	assert.Equal(t, result, expectedResults)

}

func TestCustomParseLine(t *testing.T) {

	line := "connections.active:100|gauge"
	result, err := ParseLine(line)
	require.NoError(t, err)

	assert.Equal(t, Metric{Name: "connections.active", Value: 100, Type: "gauge"}, result)

	lineTwo := "ping.amoncx.lookup_time:300|gauge"

	resultTwo, errTwo := ParseLine(lineTwo)
	require.NoError(t, errTwo)

	assert.Equal(t, Metric{Name: "ping.amoncx.lookup_time", Value: 300, Type: "gauge"}, resultTwo)

}
