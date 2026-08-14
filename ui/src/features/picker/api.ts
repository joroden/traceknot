import { request } from "../../lib/http";

export async function getContext(sessionID: string | null): Promise<string> {
  if (!sessionID) {
    return "";
  }
  const payload = await request<{ prompt: string }>(`/picker/context?session_id=${encodeURIComponent(sessionID)}`);
  return payload.prompt;
}
