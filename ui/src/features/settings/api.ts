import { getJSON, patchJSON, request } from "../../lib/http";
import type { SettingsPatch, SettingsState } from "./types";

export function getSettings(): Promise<SettingsState> {
  return getJSON<SettingsState>("/settings");
}

export function patchSettings(patch: SettingsPatch): Promise<SettingsState> {
  return patchJSON<SettingsState>("/settings", patch);
}

export function restartSessions(): Promise<{ closed: number }> {
  return request<{ closed: number }>("/settings/restart-sessions", { method: "POST" });
}
