package rag

import "testing"

func TestGenerate(t *testing.T) {
	//if testing.Short() {
	//	t.Skip("skipping integration test")
	//}

	t.Run("检索RAG", func(t *testing.T) {
		ragSearch()
	})
	//t.Run("插入RAG", func(t *testing.T) {
	//	ragInsert()
	//})
}
