import type { NodeDetail } from "../api";
import { metadataTools } from "./detailTabs";
import { ChatParts, parseChatParts } from "./node-details/ChatParts";
import { CodeBlock } from "./node-details/CodeBlock";
import { DiffBlock, looksLikeDiff } from "./node-details/DiffBlock";
import { JsonBlock } from "./node-details/JsonBlock";

export interface DetailTabContentProps {
  detail: NodeDetail;
  active: string;
  launchDetail: NodeDetail | null;
  isSubagent: boolean;
}

export function DetailTabContent({ detail, active, launchDetail, isSubagent }: DetailTabContentProps) {
  switch (active) {
    case "arguments":
      return <JsonBlock raw={detail.arguments_json} label="Arguments" fill />;
    case "result":
      return detail.result_text ? (
        looksLikeDiff(detail.result_text) ? (
          <DiffBlock text={detail.result_text} label="Result" fill />
        ) : (
          <CodeBlock text={detail.result_text} label="Result" fill />
        )
      ) : null;
    case "error":
      return <JsonBlock raw={detail.error_details_json} label="Error details" fill />;
    case "prompt":
      return detail.prompt_text ? <CodeBlock text={detail.prompt_text} label="Prompt" fill /> : null;
    case "user":
    case "system":
    case "context": {
      const messages = parseChatParts(detail.prompt_text) ?? [];
      const filtered = messages.filter((message) =>
        active === "user"
          ? message.type === "user"
          : active === "system"
            ? message.type === "system"
            : message.type !== "user" && message.type !== "system",
      );
      if (filtered.length === 0) {
        return null;
      }
      return (
        <ChatParts
          parts={filtered}
          label={active === "user" ? "User message" : active === "system" ? "System prompt" : "Context"}
        />
      );
    }
    case "output": {
      const parts = parseChatParts(detail.output_text);
      return parts ? (
        <ChatParts parts={parts} />
      ) : detail.output_text ? (
        <CodeBlock text={detail.output_text} label="Output" fill />
      ) : null;
    }
    case "reasoning": {
      const parts = parseChatParts(detail.reasoning_text);
      if (parts) {
        const reasoning = parts.filter((part) => part.type === "reasoning" || part.type === "message");
        if (reasoning.length > 0) {
          return <ChatParts parts={reasoning} allReasoning />;
        }
      }
      return detail.reasoning_text ? <CodeBlock text={detail.reasoning_text} label="Reasoning" fill /> : null;
    }
    case "tools": {
      const tools = metadataTools(detail.metadata_json);
      return tools ? <JsonBlock raw={tools} label="Tool definitions" fill /> : null;
    }
    case "launch":
      return (
        <div className="space-y-4">
          {launchDetail?.arguments_json || detail.arguments_json ? (
            <JsonBlock raw={launchDetail?.arguments_json ?? detail.arguments_json} label="Launch arguments" />
          ) : null}
          {launchDetail?.result_text ? <CodeBlock text={launchDetail.result_text} label="Launch result" /> : null}
        </div>
      );
    case "summary":
      return (
        <div className="space-y-3">
          {isSubagent && (
            <p className="text-xs text-zinc-400 light:text-zinc-600">
              Subagent spawned via {launchDetail?.tool_name ?? detail.tool_name ?? "task"} — its tools are in
              this subtree.
            </p>
          )}
          {detail.kind === "agent" && (
            <p className="text-xs text-zinc-400 light:text-zinc-600">
              This agent's task, tool calls, and results are in its subtree below.
            </p>
          )}
          {detail.kind === "chat" && detail.name === "assistant" && (
            <p className="text-xs text-zinc-400 light:text-zinc-600">
              No reply text was captured for this turn. If it made a tool call, that call is nested below it
              in the tree.
            </p>
          )}
        </div>
      );
    default:
      return null;
  }
}
