package translate

import (
	"bytes"
	"testing"
)

func TestTranslate(t *testing.T) {
	//t.Helper()
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	t.Run("stdin 到 stdout 翻译成功", func(t *testing.T) {

		in := bytes.NewBufferString("Hello world\n")
		out := &bytes.Buffer{}

		err := translateWithReference(in, out)
		if err != nil {
			t.Fatalf("翻译失败: %v", err)
		}

		if out.String() == "" {
			t.Fatal("期望有翻译输出，但结果为空")
		}

		t.Logf("翻译结果:\n%s", out.String())
	})

}
