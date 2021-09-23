package snowflake

import "testing"

func TestNewIDGenerator(t *testing.T) {
	_, err := NewIDGenerator(1, 2)
	if err != nil {
		t.Fatalf("error creating NewNode, %s", err)
	}

	_, err = NewIDGenerator(-1, -2)
	if err == nil {
		t.Fatalf("no error creating NewNode, %s", err)
	}
}

func TestGenerateDuplicateID(t *testing.T) {
	generator, _ := NewIDGenerator(1, 2)
	var x, y int64
	for i := 0; i < 1000000; i++ {
		y, _ = generator.NextId()
		if x == y {
			t.Errorf("x(%d) & y(%d) are the same", x, y)
		}
		x = y
	}
}

func BenchmarkGenerate(b *testing.B) {
	generator, _ := NewIDGenerator(1, 2)
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, _ = generator.NextId()
	}
}
