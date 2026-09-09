import { getJSON } from "../../lib/http";

export interface SettingsResponse {
  require_work_item: boolean;
}

export function getSettings(): Promise<SettingsResponse> {
  return getJSON<SettingsResponse>("/settings");
}
