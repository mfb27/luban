package providerfactory

import (
	"fmt"

	"github.com/mfb27/luban/internal/anthropic"
	"github.com/mfb27/luban/internal/openai"
	"github.com/mfb27/luban/internal/provider"
	"github.com/mfb27/luban/internal/zhipu"
)

// ProviderConfig holds configuration for a provider
type ProviderConfig struct {
	APIKey  string
	BaseURL string
}

// NewProvider creates a new provider based on the provider type
func NewProvider(providerType string, cfg *ProviderConfig) (provider.Provider, error) {
	switch providerType {
	case "openai":
		return openai.NewClient(cfg.APIKey, cfg.BaseURL), nil
	case "anthropic":
		return anthropic.NewClient(cfg.APIKey, cfg.BaseURL), nil
	case "zhipu":
		return zhipu.NewClient(cfg.APIKey, cfg.BaseURL), nil
	default:
		return nil, fmt.Errorf("unknown provider type: %s", providerType)
	}
}