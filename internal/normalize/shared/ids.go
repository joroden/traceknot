package shared

import "traceknot/internal/stableid"

func SessionID(parts ...string) string {
	return stableid.From("session", parts...)
}

func NodeID(parts ...string) string {
	return stableid.From("node", parts...)
}

func AgentID(parts ...string) string {
	return stableid.From("agent", parts...)
}
