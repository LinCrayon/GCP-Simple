package think_signature

import "testing"

func TestGenerate(t *testing.T) {

	t.Run("思考签名获取天气与新闻", func(t *testing.T) {
		thinkSignature()
	})
}
