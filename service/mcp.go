package service

import (
	"context"
	"fmt"
	"log"

	mcp_golang "github.com/metoro-io/mcp-golang"
)

type MCPService struct {
	ctx    context.Context
	client *mcp_golang.Client
}

func NewMCPService(ctx context.Context, client *mcp_golang.Client) *MCPService {
	return &MCPService{
		ctx:    ctx,
		client: client,
	}
}

func (s *MCPService) GetWritingStyle(toolName string) (string, error) {
	res, err := s.client.CallTool(context.Background(), "dazai_style", map[string]interface{}{})
	if err != nil {
		log.Printf("Failed to call tool: %v", err)
	}

	fmt.Println("Response from tool:")
	if res != nil && len(res.Content) > 0 && res.Content[0].TextContent != nil {
		return res.Content[0].TextContent.Text, nil
	}

	return "", fmt.Errorf("no content returned from tool")
}
