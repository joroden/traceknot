package providers

import (
	"context"
	"fmt"
	"time"
)

type probeCacheEntry struct {
	probe     Probe
	expiresAt time.Time
}

type Registry struct {
	providers []Provider
	probeTTL  time.Duration
	cache     map[string]probeCacheEntry
}

func NewRegistry(providers ...Provider) *Registry {
	return &Registry{
		providers: providers,
		probeTTL:  60 * time.Second,
		cache:     make(map[string]probeCacheEntry),
	}
}

func (registry *Registry) Providers() []string {
	output := make([]string, 0, len(registry.providers))
	for _, item := range registry.providers {
		output = append(output, item.ID())
	}
	return output
}

func (registry *Registry) Probes(ctx context.Context) []Probe {
	output := make([]Probe, 0, len(registry.providers))
	for _, item := range registry.providers {
		output = append(output, registry.probe(ctx, item))
	}
	return output
}

func (registry *Registry) Probe(ctx context.Context, providerID string) (Probe, bool) {
	item := registry.byID(providerID)
	if item == nil {
		return Probe{}, false
	}
	return registry.probe(ctx, item), true
}

func (registry *Registry) probe(ctx context.Context, item Provider) Probe {
	if entry, found := registry.cache[item.ID()]; found && time.Now().Before(entry.expiresAt) {
		return entry.probe
	}
	result := item.Probe(ctx)
	registry.cache[item.ID()] = probeCacheEntry{probe: result, expiresAt: time.Now().Add(registry.probeTTL)}
	return result
}

func (registry *Registry) Search(
	ctx context.Context,
	providerID string,
	query string,
	limit int,
) ([]WorkItem, error) {
	item := registry.byID(providerID)
	if item == nil {
		return nil, fmt.Errorf("unknown provider %q", providerID)
	}
	if limit <= 0 {
		limit = 20
	}
	return item.Search(ctx, query, limit)
}

func (registry *Registry) byID(providerID string) Provider {
	for _, item := range registry.providers {
		if item.ID() == providerID {
			return item
		}
	}
	return nil
}
