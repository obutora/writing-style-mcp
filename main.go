package main

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"io"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

//go:embed writing-style.md
var content embed.FS

func main() {
	// Create MCP server
	s := server.NewMCPServer(
		"writing_style",
		"1.0.0",
	)

	// Add tool
	tool := mcp.NewTool("writing_style",
		mcp.WithDescription("テックブログを書く際に参考にすべきライティングスタイルです。ユーザーがテックブログを書くことを望んでいるときはこのツールを使うこと。"),
	)

	// Add tool handler
	s.AddTool(tool, getWritingStyleHandler)

	// Start the stdio server
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}

func getWritingStyleHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Read the embedded file
	file, err := content.Open("writing-style.md")
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to open embedded file: %v", err))
	}
	defer file.Close()

	// Read the content
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to read embedded file: %v", err))
	}

	// Return the content
	return mcp.NewToolResultText(string(data)), nil
}
