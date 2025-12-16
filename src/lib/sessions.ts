export interface BackendSession {
  SessionID: number;
  Expires: number;
  Data?: string | null;
}

export async function ensureSession(): Promise<BackendSession> {
  const res = await fetch("/api/sessions", {
    method: "POST",
    credentials: "include",
  });

  if (!res.ok) {
    throw new Error(`Failed to create/reuse session: ${res.status}`);
  }

  const session = await res.json() as BackendSession;

  return session;
}