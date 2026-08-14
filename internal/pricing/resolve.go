package pricing

import (
	"sort"
	"strconv"
	"strings"
	"time"
)

func normalizeIdentifier(value string) string {
	var builder strings.Builder
	for _, char := range strings.ToLower(strings.TrimSpace(value)) {
		if (char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') {
			builder.WriteRune(char)
		} else {
			builder.WriteByte('-')
		}
	}
	normalized := strings.Trim(builder.String(), "-")
	return collapseDashes(normalized)
}

func collapseDashes(value string) string {
	for strings.Contains(value, "--") {
		value = strings.ReplaceAll(value, "--", "-")
	}
	return value
}

func canonicalizeModelIdentifier(value string) *string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil
	}
	normalized := normalizeIdentifier(trimmed)
	if normalized == "" {
		return nil
	}
	stripped := stripTrailingModelReleaseDate(normalized)
	if stripped != "" {
		normalized = stripped
	}
	return &normalized
}

func stripTrailingModelReleaseDate(normalized string) string {
	parts := strings.Split(normalized, "-")
	if len(parts) < 4 {
		return normalized
	}

	if year, err := strconv.Atoi(parts[len(parts)-3]); err == nil {
		if year >= 1900 && year <= 2099 {
			if month, err := strconv.Atoi(parts[len(parts)-2]); err == nil && month >= 1 && month <= 12 {
				if day, err := strconv.Atoi(parts[len(parts)-1]); err == nil && day >= 1 && day <= 31 {
					return strings.Join(parts[:len(parts)-3], "-")
				}
			}
		}
	}

	if len(parts[len(parts)-1]) == 8 {
		if date, err := time.Parse("20060102", parts[len(parts)-1]); err == nil {
			year := date.Year()
			if year >= 1900 && year <= 2099 {
				return strings.Join(parts[:len(parts)-1], "-")
			}
		}
	}
	return normalized
}

func tokenKey(normalized string) string {
	parts := strings.Split(normalized, "-")
	sort.Strings(parts)
	return strings.Join(parts, "-")
}

func entryMatchesTimestamp(catalogEntry *entry, usageUnixMs int64) bool {
	if catalogEntry.EffectiveFromUnix != nil && usageUnixMs < *catalogEntry.EffectiveFromUnix {
		return false
	}
	if catalogEntry.EffectiveUntilUnix != nil && usageUnixMs >= *catalogEntry.EffectiveUntilUnix {
		return false
	}
	return true
}

func pickEntryForTimestamp(entries []*entry, usageUnixMs int64) *entry {
	var matches []*entry
	for _, catalogEntry := range entries {
		if entryMatchesTimestamp(catalogEntry, usageUnixMs) {
			matches = append(matches, catalogEntry)
		}
	}
	if len(matches) == 0 {
		return nil
	}
	sort.Slice(matches, func(left, right int) bool {
		leftFrom := int64(-1 << 62)
		rightFrom := int64(-1 << 62)
		if matches[left].EffectiveFromUnix != nil {
			leftFrom = *matches[left].EffectiveFromUnix
		}
		if matches[right].EffectiveFromUnix != nil {
			rightFrom = *matches[right].EffectiveFromUnix
		}
		if leftFrom != rightFrom {
			return leftFrom > rightFrom
		}
		leftUntil := int64(1 << 62)
		rightUntil := int64(1 << 62)
		if matches[left].EffectiveUntilUnix != nil {
			leftUntil = *matches[left].EffectiveUntilUnix
		}
		if matches[right].EffectiveUntilUnix != nil {
			rightUntil = *matches[right].EffectiveUntilUnix
		}
		return leftUntil < rightUntil
	})
	return matches[0]
}

func (catalog *Catalog) ResolveEntry(model string, usageTimestampUnixMs *int64) *entry {
	normalized := canonicalizeModelIdentifier(model)
	if normalized == nil {
		return nil
	}

	usageUnix := time.Now().UnixMilli()
	if usageTimestampUnixMs != nil {
		usageUnix = *usageTimestampUnixMs
	}

	token := tokenKey(*normalized)
	var candidates []*entry
	seen := make(map[*entry]struct{})
	for _, catalogEntry := range append(append([]*entry{}, catalog.byID[*normalized]...), catalog.byToken[token]...) {
		if _, ok := seen[catalogEntry]; ok {
			continue
		}
		seen[catalogEntry] = struct{}{}
		candidates = append(candidates, catalogEntry)
	}
	return pickEntryForTimestamp(candidates, usageUnix)
}
