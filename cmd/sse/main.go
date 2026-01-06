package main

import (
	"bufio"
	"fmt"
	"log"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	spec := &openapi3.T{
		OpenAPI: "3.0.3",
		Info: &openapi3.Info{
			Title:   "Mini SSE Service",
			Version: "1.0.0",
		},
		Paths: openapi3.NewPaths(),
	}
	op := &openapi3.Operation{
		Summary:     "SSE Stream",
		Description: "Server-Sent Events stream",
		Responses:   openapi3.NewResponses(),
	}
	responseDesc := "SSE event stream"
	schema := &openapi3.Schema{
		Type: &openapi3.Types{"string"},
	}
	op.Responses.Set("200", &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: &responseDesc,
			Content: openapi3.Content{
				"text/event-stream": &openapi3.MediaType{
					Schema: &openapi3.SchemaRef{
						Value: schema,
					},
				},
			},
		},
	})
	spec.Paths.Set("/sse/stream", &openapi3.PathItem{
		Get: op,
	})
	app.Get("/sse/stream", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/event-stream")
		c.Set("Cache-Control", "no-cache")
		c.Set("Connection", "keep-alive")
		c.Set("Transfer-Encoding", "chunked")
		c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
			ticker := time.NewTicker(2 * time.Second)
			defer ticker.Stop()

			for {
				select {
				case t := <-ticker.C:
					_, err := fmt.Fprintf(w, "data: %s\n\n", t.Format(time.RFC3339))
					if err != nil {
						return
					}
					err = w.Flush()
					if err != nil {
						return
					}
				}
			}
		})
		return nil
	})
	app.Get("/docs/openapi.json", func(c *fiber.Ctx) error {
		return c.JSON(spec)
	})
	log.Println("Listening on :8081")
	log.Fatal(app.Listen(":8081"))
}
