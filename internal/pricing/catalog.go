package pricing

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"time"
)

//go:embed pricings.json
var embeddedCatalogJSON []byte

type entry struct {
	ID                 string   `json:"id"`
	Input              float64  `json:"input"`
	Output             float64  `json:"output"`
	InputCached        *float64 `json:"input_cached"`
	CacheWrite         *float64 `json:"cache_write"`
	CacheWrite1h       *float64 `json:"cache_write_1h"`
	EffectiveFrom      *string  `json:"effective_from"`
	EffectiveUntil     *string  `json:"effective_until"`
	EffectiveFromUnix  *int64
	EffectiveUntilUnix *int64
}

type Catalog struct {
	Entries []entry
	byID    map[string][]*entry
	byToken map[string][]*entry
}

func loadCatalog(data []byte) (*Catalog, error) {
	var raw struct {
		Prices []entry `json:"prices"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("parse pricing catalog: %w", err)
	}

	catalog := &Catalog{
		Entries: raw.Prices,
		byID:    make(map[string][]*entry, len(raw.Prices)),
		byToken: make(map[string][]*entry, len(raw.Prices)),
	}

	for index := range catalog.Entries {
		catalogEntry := &catalog.Entries[index]
		if catalogEntry.Input < 0 || catalogEntry.Output < 0 {
			return nil, fmt.Errorf("pricing entry %q has negative rates", catalogEntry.ID)
		}
		from, err := parseTimestamp(catalogEntry.EffectiveFrom, fmt.Sprintf("entry %d effective_from", index))
		if err != nil {
			return nil, err
		}
		until, err := parseTimestamp(catalogEntry.EffectiveUntil, fmt.Sprintf("entry %d effective_until", index))
		if err != nil {
			return nil, err
		}
		if from != nil && until != nil && *from >= *until {
			return nil, fmt.Errorf("pricing entry %q effective_from is not earlier than effective_until", catalogEntry.ID)
		}
		catalogEntry.EffectiveFromUnix = from
		catalogEntry.EffectiveUntilUnix = until

		normalized := normalizeIdentifier(catalogEntry.ID)
		canonical := normalized
		if stripped := stripTrailingModelReleaseDate(normalized); stripped != "" {
			canonical = stripped
		}
		catalog.byID[normalized] = append(catalog.byID[normalized], catalogEntry)
		if canonical != normalized {
			catalog.byID[canonical] = append(catalog.byID[canonical], catalogEntry)
		}
		token := tokenKey(normalized)
		catalog.byToken[token] = append(catalog.byToken[token], catalogEntry)
		if canonicalToken := tokenKey(canonical); canonicalToken != token {
			catalog.byToken[canonicalToken] = append(catalog.byToken[canonicalToken], catalogEntry)
		}
	}

	return catalog, nil
}

func LoadEmbeddedCatalog() (*Catalog, error) {
	return loadCatalog(embeddedCatalogJSON)
}

func parseTimestamp(value *string, label string) (*int64, error) {
	if value == nil {
		return nil, nil
	}
	parsed, err := time.Parse(time.RFC3339, *value)
	if err != nil {
		parsed, err = time.Parse("2006-01-02", *value)
	}
	if err != nil {
		return nil, fmt.Errorf("invalid %s %q: %w", label, *value, err)
	}
	unix := parsed.UnixMilli()
	return &unix, nil
}
