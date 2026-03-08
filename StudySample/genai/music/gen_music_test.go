package main

import (
	"testing"
)

func TestGenerate(t *testing.T) {
	//if testing.Short() {
	//	t.Skip("skipping integration test")
	//}

	t.Run("生成音乐", func(t *testing.T) {
		GenerateMusicStream()
	})
}
