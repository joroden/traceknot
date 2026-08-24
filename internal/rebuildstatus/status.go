package rebuildstatus

import (
	"sort"
	"sync"
)

type Status struct {
	mu      sync.Mutex
	pending map[string]bool
}

func New() *Status {
	return &Status{pending: map[string]bool{}}
}

func (status *Status) Start(provider string) {
	status.mu.Lock()
	defer status.mu.Unlock()
	status.pending[provider] = true
}

func (status *Status) Done(provider string) {
	status.mu.Lock()
	defer status.mu.Unlock()
	delete(status.pending, provider)
}

func (status *Status) InProgress() []string {
	status.mu.Lock()
	defer status.mu.Unlock()
	providers := make([]string, 0, len(status.pending))
	for provider := range status.pending {
		providers = append(providers, provider)
	}
	sort.Strings(providers)
	return providers
}
