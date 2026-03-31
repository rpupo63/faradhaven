export const getBaseUrl = () => import.meta.env?.VITE_API_BASE_URL ?? '';

export async function handleResponse<T>(res: Response, defaultError: string): Promise<T> {
  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: res.statusText }));
    throw new Error((err as { error?: string }).error ?? defaultError);
  }
  return res.json();
}

/**
 * Generic API request function.
 * @param method HTTP method (GET, POST, PUT, DELETE)
 * @param path API endpoint path (e.g., /api/monsters)
 * @param data Request body (for POST/PUT)
 * @param token Authorization token
 * @param defaultError Custom error message
 */
export async function makeApiRequest<T>(
  method: 'GET' | 'POST' | 'PUT' | 'DELETE',
  path: string,
  data: object | null = null,
  token?: string,
  defaultError: string = `Failed to ${method} ${path}`
): Promise<T> {
  const base = getBaseUrl();
  const url = `${base}${path}`;

  const headers: Record<string, string> = {};
  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }
  if (data && (method === 'POST' || method === 'PUT')) {
    headers['Content-Type'] = 'application/json';
  }

  const options: RequestInit = {
    method,
    headers,
  };

  if (data && (method === 'POST' || method === 'PUT')) {
    options.body = JSON.stringify(data);
  }

  const res = await fetch(url, options);
  return handleResponse<T>(res, defaultError);
}

