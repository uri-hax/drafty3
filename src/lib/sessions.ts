import { getUsersAPI } from "./api";

export interface BackendProfile {
  IDProfile: number;
  DateCreated: string;
  DateUpdated: string;
}

export interface BackendSession {
  IDSession: number;
  IDProfile: number;
  Start: string;
  End: string;
}

export interface EnsureSessionResponse {
  profile: BackendProfile;
  session: BackendSession;
}

export async function ensureSession(): Promise<EnsureSessionResponse> {
  const res = await fetch(`${getUsersAPI()}/sessions`, {
    method: "POST",
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
    },
  });

  if (!res.ok) {
    throw new Error(`Failed to create/reuse session: ${res.status}`);
  }

  const payload = (await res.json()) as EnsureSessionResponse;
  return payload;
}