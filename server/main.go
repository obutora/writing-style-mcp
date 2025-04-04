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

//go:embed styles/*
var content embed.FS

func main() {
	// Create MCP server
	s := server.NewMCPServer(
		"writing_style",
		"1.0.0",
	)

	// Add techBlogWritingStyleTool
	techBlogWritingStyleTool := mcp.NewTool("tech_blog_style",
		mcp.WithDescription("テックブログを書く際に参考にすべきライティングスタイルです。ユーザーがテックブログを書くことを望んでいるときはこのツールを使うこと。"),
	)
	s.AddTool(techBlogWritingStyleTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return getWritingStyleHandler(ctx, request, "styles/tech-blog-style.md")
	})

	dazaiStyleTool := mcp.NewTool("dazai_style",
		mcp.WithDescription("太宰治の文体を模倣するツールです。ユーザーが太宰治の文体を望んでいるときはこのツールを使うこと。"),
	)
	s.AddTool(dazaiStyleTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return getWritingStyleHandler(ctx, request, "styles/dazai-style.md")
	})

	lpStyleTool := mcp.NewTool("lp_style",
		mcp.WithDescription("LPスタイルのライティングを模倣するツールです。ユーザーがLPスタイルのライティングを望んでいるときはこのツールを使うこと。"),
	)
	s.AddTool(lpStyleTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return getWritingStyleHandler(ctx, request, "styles/lp-style.md")
	})

	offTopicStyleTool := mcp.NewTool("off_topic_style",
		mcp.WithDescription("オフトピックなライティングスタイルを模倣するツールです。ユーザーがオフトピックなライティングスタイルを望んでいるときはこのツールを使うこと。"),
	)
	s.AddTool(offTopicStyleTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return getWritingStyleHandler(ctx, request, "styles/off-topic-style.md")
	})

	// Start the stdio server
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}

func getWritingStyleHandler(ctx context.Context, request mcp.CallToolRequest, filePath string) (*mcp.CallToolResult, error) {
	// Read the embedded file
	file, err := content.Open(filePath)
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
