export function getAPI(): string {
  if (typeof window === "undefined") return "";

  const { hostname, pathname } = window.location;

  const base = import.meta.env.BASE_URL;
  const dataset = pathname
    .replace(base, "")
    .split("/")
    .filter(Boolean)[0];

  const localApiBase = (import.meta.env.VITE_PUBLIC_LOCAL_DEV_API || "").replace(/\/$/, "");
  const apiBase = (import.meta.env.VITE_PUBLIC_API_BASE || "").replace(/\/$/, "");

  if (hostname === "localhost") {
    return `${localApiBase}/${dataset}`;
  }

  if (
    hostname === "uri-hax.github.io" &&
    pathname.startsWith(import.meta.env.VITE_PUBLIC_BASE_PATH)
  ) {
    return `${apiBase}/${dataset}`;
  }

  return `${apiBase}/${dataset}`;
}

export function getUsersAPI(): string {
  if (typeof window === "undefined") return "";

  const { hostname, pathname } = window.location;

  const localApiBase = (import.meta.env.VITE_PUBLIC_LOCAL_DEV_API || "").replace(/\/$/, "");
  const apiBase = (import.meta.env.VITE_PUBLIC_API_BASE || "").replace(/\/$/, "");

  if (hostname === "localhost") {
    return `${localApiBase}/users`;
  }

  if (
    hostname === "uri-hax.github.io" &&
    pathname.startsWith(import.meta.env.VITE_PUBLIC_BASE_PATH)
  ) {
    return `${apiBase}/users`;
  }

  return `${apiBase}/users`;
}