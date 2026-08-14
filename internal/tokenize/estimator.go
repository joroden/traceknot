package tokenize

import "sync"

const defaultEstimateCacheBytes = 64 * 1024 * 1024

type Estimator struct {
	mu       sync.Mutex
	cache    map[string]int64
	cached   int64
	maxBytes int64
}

func NewEstimator(maxBytes int64) *Estimator {
	if maxBytes <= 0 {
		maxBytes = defaultEstimateCacheBytes
	}
	return &Estimator{
		cache:    make(map[string]int64),
		maxBytes: maxBytes,
	}
}

func (estimator *Estimator) Count(text string) int64 {
	estimator.mu.Lock()
	defer estimator.mu.Unlock()
	if count, ok := estimator.cache[text]; ok {
		return count
	}
	var count int64
	if counter := O200K(); counter != nil {
		count = counter.Count(text)
	} else {
		count = HeuristicCount(text)
	}
	if estimator.cached+int64(len(text)) > estimator.maxBytes {
		estimator.cache = make(map[string]int64)
		estimator.cached = 0
	}
	estimator.cache[text] = count
	estimator.cached += int64(len(text))
	return count
}

func (estimator *Estimator) MethodName() string {
	if counter := O200K(); counter != nil {
		return counter.Name()
	}
	return "heuristic"
}
