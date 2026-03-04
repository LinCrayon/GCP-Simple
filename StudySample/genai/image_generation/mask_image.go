package image_generation

import (
	"context"
	"fmt"
	"io"
	"os"

	"google.golang.org/genai"
)

/*功能流程：
1. 创建 GenAI Client（Vertex AI 模式）
2. 读取原图（fruit.png）
3. 读取 mask 图（fruit_mask.png）
4. 构造 RawReferenceImage（原图）
5. 构造 MaskReferenceImage（mask，绑定到原图）
6. 调用 EditImage 进行 inpaint 插入生成
7. 将生成结果写入 output-image.png
8. 将结果信息写入 io.Writer（便于测试 / HTTP / CLI 复用）
*/

func MaskImage(w io.Writer) error {
	ctx := context.Background()

	//Imagen 属于 Vertex AI
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
	})
	if err != nil {
		panic(err)
	}

	//读取原图
	imgBytes, err := os.ReadFile("fruit.png")
	if err != nil {
		panic(err)
	}

	rawRef := &genai.RawReferenceImage{
		ReferenceID: 1, //原图的逻辑 ID
		ReferenceImage: &genai.Image{
			ImageBytes: imgBytes,
			MIMEType:   "image/png",
		},
	}

	//读取 mask 图片
	// mask 图通常：
	// - 白色：需要被编辑 / 替换的区域
	// - 黑色：保持不变
	maskBytes, _ := os.ReadFile("fruit_mask.png")
	// mask 膨胀参数（可选）：
	// - 0   ：严格使用 mask 边界
	// - > 0 ：向外扩展一点，防止边缘生硬
	dilation := float32(0.01)
	maskRef := &genai.MaskReferenceImage{
		ReferenceID: 1, //作用于原图
		ReferenceImage: &genai.Image{
			ImageBytes: maskBytes,
			MIMEType:   "image/png",
		},
		Config: &genai.MaskReferenceConfig{
			MaskMode:     genai.MaskReferenceModeMaskModeUserProvided, // 使用用户提供的 mask 图片
			MaskDilation: &dilation,
		},
	}

	resp, err := client.Models.EditImage(
		ctx,
		"imagen-3.0-capability-001",
		"A plate of cookies",
		[]genai.ReferenceImage{
			rawRef,
			maskRef,
		},
		&genai.EditImageConfig{
			EditMode: genai.EditModeInpaintInsertion,
		},
	)
	if err != nil {
		panic(err)
	}

	outputFile := "output-image.png"
	if err := os.WriteFile(outputFile, resp.GeneratedImages[0].Image.ImageBytes, 0644); err != nil {
		panic(err)
	}

	fmt.Printf(
		"Created output image using %d bytes → %s\n",
		len(resp.GeneratedImages[0].Image.ImageBytes),
		outputFile,
	)
	fmt.Fprintln(w, len(resp.GeneratedImages[0].Image.ImageBytes), outputFile)

	return nil
}
