import { getAPI } from "./api";

// check health of backend server
export async function checkBackendHealth(): Promise<boolean> {
  try {
    const res = await fetch(`${getAPI()}/health`, {
      method: "GET",
      credentials: "include",
    });

    return res.ok;
  } 
  catch (err) {
    console.error("Network error checking health:", err);
    return false;
  }
}