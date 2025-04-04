package main

import (
	"context"
	"encoding/json"
	"hagakun/service"
	"log"
	"net/http"
	"os/exec"

	mcp_golang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport/stdio"
)

func main() {
	// Start the server process
	cmd := exec.Command("go", "run", "./server/main.go")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatalf("Failed to get stdin pipe: %v", err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalf("Failed to get stdout pipe: %v", err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	defer cmd.Process.Kill()

	clientTransport := stdio.NewStdioServerTransportWithIO(stdout, stdin)
	client := mcp_golang.NewClient(clientTransport)

	if _, err := client.Initialize(context.Background()); err != nil {
		log.Fatalf("Failed to initialize client: %v", err)
	}

	// List available tools
	tools, err := client.ListTools(context.Background(), nil)
	if err != nil {
		log.Fatalf("Failed to list tools: %v", err)
	}

	log.Println("Available Tools:")
	for _, tool := range tools.Tools {
		desc := ""
		if tool.Description != nil {
			desc = *tool.Description
		}
		log.Printf("Tool: %s. Description: %s", tool.Name, desc)
	}

	ctx := context.Background()
	s := service.NewMCPService(ctx, client)

	// ツール一覧を取得するエンドポイント
	http.HandleFunc("/tools", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// JSONレスポンスを返す
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"tools": tools.Tools,
		})
	})

	// スタイル取得のエンドポイント
	http.HandleFunc("/style", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// クエリパラメータからtoolNameを取得
		toolName := r.URL.Query().Get("toolName")
		if toolName == "" {
			http.Error(w, "toolName parameter is required", http.StatusBadRequest)
			return
		}

		// GetWritingStyleを呼び出して結果を取得
		styleText, err := s.GetWritingStyle(toolName)
		if err != nil {
			http.Error(w, "Error getting writing style: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// 結果を返す
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte(styleText))
	})

	// サーバーの起動
	log.Println("Starting HTTP server on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
