package claude

const (
	eventUserPrompt       = "user_prompt"
	eventAPIRequest       = "api_request"
	eventAssistantResp    = "assistant_response"
	eventAPIRequestBody   = "api_request_body"
	eventAPIResponseBody  = "api_response_body"
	eventToolDecision     = "tool_decision"
	eventToolResult       = "tool_result"
	eventSubagentComplete = "subagent_completed"

	toolNameAgent = "Agent"

	querySourceMain = "repl_main_thread"

	querySourceAgentPrefix = "agent:"

	querySourceWebSearch = "web_search_tool"

	chatNameMeta = "meta"
)
