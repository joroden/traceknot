export interface ProviderCopy {
  label: string;
  description: string;
}

export const HOOK_COPY: Record<string, ProviderCopy> = {
  claude: {
    label: "Claude Code",
    description: "Collects telemetry from every Claude Code session. Takes effect the next time you start one.",
  },
  codex: {
    label: "Codex CLI",
    description:
      "Collects telemetry from every Codex CLI session. After enabling, run /hooks inside the Codex CLI and trust the traceknot hook.",
  },
  copilot: {
    label: "VS Code Copilot Chat / Copilot CLI",
    description:
      "Collects telemetry from every Copilot Chat and Copilot CLI session. Also turns on Copilot's experimental setting and pre-approves the traceknot extension's permission prompt, so you won't be asked.",
  },
};

export const SKILL_COPY: Record<string, ProviderCopy> = {
  claude: {
    label: "Claude Code",
    description: "Lets Claude Code export and read its own session telemetry when you ask it to.",
  },
  codex: {
    label: "Codex CLI",
    description: "Lets Codex CLI export and read its own session telemetry when you ask it to.",
  },
  copilot: {
    label: "Copilot CLI",
    description: "Lets Copilot CLI export and read its own session telemetry when you ask it to.",
  },
};

export function providerCopy(table: Record<string, ProviderCopy>, binary: string): ProviderCopy {
  return table[binary] ?? { label: binary, description: "" };
}
