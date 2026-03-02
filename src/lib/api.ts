export function getAPI(): string {
  if (typeof window === 'undefined') return '';

  const { hostname, pathname } = window.location;

  const base = import.meta.env.BASE_URL;
  const dataset = window.location.pathname
    .replace(base, "")
    .split("/")
    .filter(Boolean)[0];

  if (hostname === "localhost") return dataset ? `/api/${dataset}` : "/api/csprofs"; // temp default
  const apiBase = (import.meta.env.PUBLIC_API_BASE || "").replace(/\/$/, "");

  if (
    hostname === 'uri-hax.github.io' &&
    pathname.startsWith(import.meta.env.PUBLIC_BASE_PATH)
  ) {
    return `${apiBase}${dataset}`;
  }

  return `${apiBase}${dataset}`;
}
