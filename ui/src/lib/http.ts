const API_BASE = "/api/v1";

export async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const response = await fetch(API_BASE + path, init);
  if (!response.ok) {
    let message = response.statusText;
    try {
      const body: unknown = await response.json();
      if (body && typeof body === "object" && "message" in body && typeof body.message === "string") {
        message = body.message;
      }
    } catch {
    }
    throw new Error(message);
  }
  return (await response.json()) as T;
}

export function postJSON(path: string, body: unknown): Promise<unknown> {
  return request(path, {
    method: "POST",
    headers: { "content-type": "application/json" },
    body: JSON.stringify(body),
  });
}

export function getJSON<T>(
  path: string,
  params?: Record<string, string | number | undefined>,
): Promise<T> {
  const entries = Object.entries(params ?? {}).filter(
    (entry): entry is [string, string | number] => entry[1] !== undefined && entry[1] !== "",
  );
  const query = entries.length > 0
    ? "?" + new URLSearchParams(entries.map(([key, value]) => [key, String(value)])).toString()
    : "";
  return request<T>(path + query);
}
