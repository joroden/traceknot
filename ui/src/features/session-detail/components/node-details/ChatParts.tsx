import { ChatBlock } from "./ChatBlock";

export interface ChatPart {
  type: string;
  text: string;
}

export function parseChatParts(raw: string | null): ChatPart[] | null {
  if (!raw) {
    return null;
  }
  try {
    const parsed: unknown = JSON.parse(raw);
    if (!Array.isArray(parsed)) {
      return null;
    }
    const parts: ChatPart[] = [];
    for (const entry of parsed) {
      if (entry !== null && typeof entry === "object") {
        const record = entry as { text?: unknown; type?: unknown };
        if (typeof record.text === "string" && record.text.trim()) {
          parts.push({
            type: typeof record.type === "string" ? record.type : "message",
            text: record.text,
          });
        }
      }
    }
    return parts.length > 0 ? parts : null;
  } catch {
    return null;
  }
}

function labelForType(type: string): string {
  switch (type) {
    case "user":
      return "User message";
    case "system":
      return "System prompt";
    case "assistant":
    case "message":
      return "Output";
    default:
      return "Context";
  }
}

interface ChatPartsProps {
  parts: ChatPart[];
  label?: string;
  allReasoning?: boolean;
}

export function ChatParts({ parts, label, allReasoning }: ChatPartsProps) {
  if (parts.length === 0) {
    return null;
  }
  return (
    <div className="flex min-h-0 flex-1 flex-col gap-3">
      {parts.map((part, index) => {
        const isReasoning = allReasoning || part.type === "reasoning";
        return (
          <ChatBlock
            key={index}
            label={isReasoning ? "Reasoning" : (label ?? labelForType(part.type))}
            text={part.text}
            collapsible={isReasoning}
            muted={isReasoning}
          />
        );
      })}
    </div>
  );
}
