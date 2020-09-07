package camera

import "testing"

func BenchmarkListCameras(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := Cameras()
		if err != nil {
			b.Error("Failed to get cameras:", err)
		}
	}
}
