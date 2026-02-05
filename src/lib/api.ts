export function getAPI(): string {
  if (typeof window === 'undefined') return '';

  const { hostname, pathname } = window.location;

  if (hostname === 'localhost') return '/api';

  if (
    hostname === 'uri-hax.github.io' &&
    pathname.startsWith(import.meta.env.PUBLIC_BASE_PATH)
  ) {
    return import.meta.env.PUBLIC_API_BASE;
  }

  return import.meta.env.PUBLIC_API_BASE;
}
