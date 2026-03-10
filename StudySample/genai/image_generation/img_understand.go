package image_generation

import (
	"context"
	"fmt"
	"google.golang.org/genai"
	"io"
	"os"
)

/**
 * @Description
 * @Author linshengqian
 * @Date 2026/2/28 15:29
 **/
func understandImage(w io.Writer) error {
	ctx := context.Background()

	getenv := os.Getenv("GEMINI_API_KEY")

	fmt.Println(getenv)

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: os.Getenv("GEMINI_API_KEY"),
	})
	if err != nil {
		return err
	}

	uploadedFile, err := client.Files.UploadFromPath(
		ctx,
		"bw-example-image.png",
		nil,
	)
	if err != nil {
		return err
	}

	parts := []*genai.Part{
		genai.NewPartFromText("Caption this image."),
		genai.NewPartFromURI(uploadedFile.URI, uploadedFile.MIMEType),
	}

	contents := []*genai.Content{
		genai.NewContentFromParts(parts, genai.RoleUser),
	}

	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-3-flash-preview",
		contents,
		nil,
	)
	if err != nil {
		return err
	}

	fmt.Fprintln(w, result.Text())
	return nil
}
