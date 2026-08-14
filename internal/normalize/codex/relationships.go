package codex

import (
	"encoding/json"
	"sort"
)

const waitAgentToolName = "multi_agent_v1wait_agent"

type subagentLink struct {
	conversationID string
	nickname       string
	spawnCallID    string
	spawnPrompt    string
}

type approvalInfo struct {
	decision string
	source   string
}

func rootConversation(all map[string][]Event) string {
	spawned := make(map[string]struct{})
	for _, events := range all {
		for _, event := range events {
			if link := parseSpawnLink(event); link != nil {
				spawned[link.conversationID] = struct{}{}
			}
		}
	}
	var candidates []string
	for id := range all {

		if id == "" {
			continue
		}
		if _, ok := spawned[id]; !ok {
			candidates = append(candidates, id)
		}
	}
	if len(candidates) == 1 {
		return candidates[0]
	}
	if len(candidates) == 0 {
		return ""
	}
	sort.Strings(candidates)
	return candidates[0]
}

func subagentMap(all map[string][]Event) map[string]*subagentLink {
	output := make(map[string]*subagentLink)
	for _, events := range all {
		for _, event := range events {
			if event.Name != eventToolResult || toolName(event) != spawnToolName {
				continue
			}
			if link := parseSpawnLink(event); link != nil {
				output[link.conversationID] = link
			}
		}
	}
	return output
}

func approvalsByCallID(all map[string][]Event) map[string]approvalInfo {
	output := make(map[string]approvalInfo)
	for _, events := range all {
		for _, event := range events {
			if event.Name != eventToolDecision {
				continue
			}
			callID, ok := attributeString(event.Attributes, "call_id")
			if !ok || callID == "" {
				continue
			}
			decision, _ := attributeString(event.Attributes, "decision")
			source, _ := attributeString(event.Attributes, "source")
			output[callID] = approvalInfo{decision: decision, source: source}
		}
	}
	return output
}

func waitResultsByAgentID(all map[string][]Event) map[string]string {
	output := make(map[string]string)
	for _, events := range all {
		for _, event := range events {
			if event.Name != eventToolResult || toolName(event) != waitAgentToolName {
				continue
			}
			text, ok := event.Attributes["output"].(string)
			if !ok || text == "" {
				continue
			}
			var parsed struct {
				Status map[string]struct {
					Completed string `json:"completed"`
				} `json:"status"`
			}
			if err := json.Unmarshal([]byte(text), &parsed); err != nil {
				continue
			}
			for agentID, entry := range parsed.Status {
				if entry.Completed != "" {
					output[agentID] = entry.Completed
				}
			}
		}
	}
	return output
}

func findSpawnLink(subagents map[string]*subagentLink, callID string) *subagentLink {
	for _, link := range subagents {
		if link.spawnCallID == callID {
			return link
		}
	}
	return nil
}

func parseSpawnLink(event Event) *subagentLink {
	if event.Name != eventToolResult || toolName(event) != spawnToolName {
		return nil
	}
	output := event.Attributes["output"]
	text, ok := output.(string)
	if !ok || text == "" {
		return nil
	}
	var parsed struct {
		AgentID  string `json:"agent_id"`
		Nickname string `json:"nickname"`
	}
	if err := json.Unmarshal([]byte(text), &parsed); err != nil {
		return nil
	}
	if parsed.AgentID == "" {
		return nil
	}
	callID, _ := attributeString(event.Attributes, "call_id")
	prompt, _ := attributeString(event.Attributes, "arguments")
	return &subagentLink{
		conversationID: parsed.AgentID,
		nickname:       parsed.Nickname,
		spawnCallID:    callID,
		spawnPrompt:    spawnPrompt(prompt),
	}
}

func spawnPrompt(arguments string) string {
	if arguments == "" {
		return ""
	}
	var parsed struct {
		Message string `json:"message"`
	}
	if err := json.Unmarshal([]byte(arguments), &parsed); err != nil {
		return ""
	}
	return parsed.Message
}
