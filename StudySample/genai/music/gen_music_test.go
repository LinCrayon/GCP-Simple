package music

import (
	"testing"
)

func TestGenerate(t *testing.T) {
	//if testing.Short() {
	//	t.Skip("skipping integration test")
	//}

	t.Run("提出你的第一个请求", func(t *testing.T) {
		GenerateMusicStream()
	})
}
