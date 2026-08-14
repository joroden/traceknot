import type { NodeDetail } from "../api";
import { parseChatParts } from "./node-details/ChatParts";

export interface Tab {
  key: string;
  label: string;
}

export function metadataTools(metadataJSON: string): string | null {
  try {
    const parsed: unknown = JSON.parse(metadataJSON);
    if (parsed !== null && typeof parsed === "object") {
      const tools = (parsed as { tools?: unknown }).tools;
      if (typeof tools === "string" && tools.trim()) {
        return tools;
      }
    }
  } catch {
    return null;
  }
  return null;
}

export function buildTabs(detail: NodeDetail, isSubagent: boolean, launchDetail: NodeDetail | null): Tab[] {
  const tabs: Tab[] = [];
  if (isSubagent) {
    tabs.push({ key: "summary", label: "Summary" });
    if (launchDetail?.arguments_json || detail.arguments_json) {
      tabs.push({ key: "launch", label: "Launch" });
    }
    if (detail.prompt_text) {
      tabs.push({ key: "prompt", label: "Prompt" });
    }
    if (detail.result_text) {
      tabs.push({ key: "result", label: "Result" });
    }
    return tabs;
  }
  if (detail.kind === "tool_call") {
    if (detail.arguments_json) {
      tabs.push({ key: "arguments", label: "Arguments" });
    }
    if (detail.result_text) {
      tabs.push({ key: "result", label: "Result" });
    }
    if (detail.error_details_json) {
      tabs.push({ key: "error", label: "Error" });
    }
    return tabs;
  }

  const messages = parseChatParts(detail.prompt_text);
  if (messages !== null) {
    if (messages.some((message) => message.type === "system")) {
      tabs.push({ key: "system", label: "System" });
    }
    if (messages.some((message) => message.type === "user")) {
      tabs.push({ key: "user", label: "User" });
    }
    if (messages.some((message) => message.type !== "user" && message.type !== "system")) {
      tabs.push({ key: "context", label: "Context" });
    }
  } else if (detail.prompt_text) {
    tabs.push({ key: "prompt", label: "Prompt" });
  }
  if (detail.output_text) {
    tabs.push({ key: "output", label: "Output" });
  }
  if (detail.reasoning_text) {
    tabs.push({ key: "reasoning", label: "Reasoning" });
  }
  if (detail.kind === "agent") {
    tabs.push({ key: "summary", label: "Summary" });
  }
  if (detail.kind === "chat" && detail.name === "assistant" && !detail.output_text && !detail.reasoning_text) {
    tabs.push({ key: "summary", label: "Summary" });
  }
  if (metadataTools(detail.metadata_json)) {
    tabs.push({ key: "tools", label: "Tools" });
  }
  return tabs;
}
