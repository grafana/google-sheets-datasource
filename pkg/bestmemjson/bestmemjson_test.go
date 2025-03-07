package bestmemjson

import (
	"encoding/json"
	"testing"

	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	String  string `json:"string"`
	Number  int    `json:"number"`
	Boolean bool   `json:"boolean"`
	Array   []int  `json:"array"`
}

var testJSON = []byte(`{
	"string": "test string",
	"number": 42,
	"boolean": true,
	"array": [1, 2, 3, 4, 5]
}`)

var testData = &testStruct{
	String:  "test string",
	Number:  42,
	Boolean: true,
	Array:   []int{1, 2, 3, 4, 5},
}

func TestBestMemJSONEfficiency(t *testing.T) {
	t.Run("Marshal memory efficiency", func(t *testing.T) {
		// Create sub-benchmarks to measure memory
		standardResult := testing.Benchmark(func(b *testing.B) {
			b.ResetTimer()
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_, _ = json.Marshal(testData)
			}
		})

		jsoniterResult := testing.Benchmark(func(b *testing.B) {
			b.ResetTimer()
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_, _ = jsoniter.Marshal(testData)
			}
		})

		bestmemResult := testing.Benchmark(func(b *testing.B) {
			b.ResetTimer()
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_, _ = Marshal(testData)
			}
		})

		// Print the actual memory usage for verification
		t.Logf("Memory usage per operation:\n"+
			"encoding/json: %d bytes/op, %d allocs/op\n"+
			"jsoniter: %d bytes/op, %d allocs/op\n"+
			"bestmemjson: %d bytes/op, %d allocs/op",
			standardResult.AllocedBytesPerOp(),
			standardResult.AllocsPerOp(),
			jsoniterResult.AllocedBytesPerOp(),
			jsoniterResult.AllocsPerOp(),
			bestmemResult.AllocedBytesPerOp(),
			bestmemResult.AllocsPerOp())

		// Verify bestmemjson is using the most memory-efficient implementation
		assert.LessOrEqual(t, bestmemResult.AllocedBytesPerOp(), standardResult.AllocedBytesPerOp(),
			"bestmemjson.Marshal should not use more memory than encoding/json.Marshal")
		assert.LessOrEqual(t, bestmemResult.AllocedBytesPerOp(), jsoniterResult.AllocedBytesPerOp(),
			"bestmemjson.Marshal should not use more memory than jsoniter.Marshal")
	})

	t.Run("Unmarshal memory efficiency", func(t *testing.T) {
		var standardModel, jsoniterModel, bestmemModel testStruct

		standardResult := testing.Benchmark(func(b *testing.B) {
			b.ResetTimer()
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = json.Unmarshal(testJSON, &standardModel)
			}
		})

		jsoniterResult := testing.Benchmark(func(b *testing.B) {
			b.ResetTimer()
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = jsoniter.Unmarshal(testJSON, &jsoniterModel)
			}
		})

		bestmemResult := testing.Benchmark(func(b *testing.B) {
			b.ResetTimer()
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = Unmarshal(testJSON, &bestmemModel)
			}
		})

		// Print the actual memory usage for verification
		t.Logf("Memory usage per operation:\n"+
			"encoding/json: %d bytes/op, %d allocs/op\n"+
			"jsoniter: %d bytes/op, %d allocs/op\n"+
			"bestmemjson: %d bytes/op, %d allocs/op",
			standardResult.AllocedBytesPerOp(),
			standardResult.AllocsPerOp(),
			jsoniterResult.AllocedBytesPerOp(),
			jsoniterResult.AllocsPerOp(),
			bestmemResult.AllocedBytesPerOp(),
			bestmemResult.AllocsPerOp())

		// Verify bestmemjson is using the most memory-efficient implementation
		assert.LessOrEqual(t, bestmemResult.AllocedBytesPerOp(), standardResult.AllocedBytesPerOp(),
			"bestmemjson.Unmarshal should not use more memory than encoding/json.Unmarshal")
		assert.LessOrEqual(t, bestmemResult.AllocedBytesPerOp(), jsoniterResult.AllocedBytesPerOp(),
			"bestmemjson.Unmarshal should not use more memory than jsoniter.Unmarshal")

		// Verify all implementations produce the same result
		assert.Equal(t, standardModel, bestmemModel, "bestmemjson.Unmarshal should produce same result as encoding/json")
		assert.Equal(t, jsoniterModel, bestmemModel, "bestmemjson.Unmarshal should produce same result as jsoniter")
	})
}
