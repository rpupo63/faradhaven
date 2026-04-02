export const getBaseUrl = () => import.meta.env?.VITE_API_BASE_URL ?? '';

export async function handleResponse<T>(res: Response, defaultError: string): Promise<T> {
  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: res.statusText }));
    throw new Error((err as { error?: string }).error ?? defaultError);
  }
  return res.json();
}

let isRefreshing = false;
let refreshSubscribers: ((token: string) => void)[] = [];

function onRefreshed(token: string) {
  refreshSubscribers.forEach((cb) => cb(token));
  refreshSubscribers = [];
}

function addRefreshSubscriber(cb: (token: string) => void) {
  refreshSubscribers.push(cb);
}

/**
 * Wrapper around fetch that handles auth token refreshing.
 */
export async function apiFetch(input: RequestInfo | URL, init?: RequestInit): Promise<Response> {
  const options = init || {};
  options.credentials = 'include'; // Ensure cookies are sent for refresh token

  const res = await fetch(input, options);

  if (res.status === 401) {
    const url = typeof input === 'string' ? input : input instanceof URL ? input.toString() : input.url;
    // Don't intercept refresh or login requests
    if (url.includes('/api/auth/refresh') || url.includes('/api/auth/login')) {
      return res;
    }

    if (!isRefreshing) {
      isRefreshing = true;
      try {
        const base = getBaseUrl();
        const refreshRes = await fetch(`${base}/api/auth/refresh`, {
          method: 'POST',
          credentials: 'include' // need cookie to refresh
        });

        if (refreshRes.ok) {
          const data = await refreshRes.json();
          const newToken = data.token;
          const newUserId = data.user_id;

          // Update local storage so other tabs/contexts get the new token
          localStorage.setItem('faradhaven-auth-token', newToken);
          localStorage.setItem('faradhaven-user-id', newUserId);

          isRefreshing = false;
          onRefreshed(newToken);

          // Retry the original request
          const newHeaders = new Headers(options.headers);
          newHeaders.set('Authorization', `Bearer ${newToken}`);
          return fetch(input, { ...options, headers: newHeaders });
        } else {
          // Refresh failed, clear auth and redirect to login
          localStorage.removeItem('faradhaven-auth-token');
          localStorage.removeItem('faradhaven-user-id');
          window.location.href = '/login';
          throw new Error('Session expired');
        }
      } catch (err) {
        isRefreshing = false;
        localStorage.removeItem('faradhaven-auth-token');
        localStorage.removeItem('faradhaven-user-id');
        window.location.href = '/login';
        throw err;
      }
    } else {
      // Wait for the refresh to complete
      return new Promise((resolve) => {
        addRefreshSubscriber((newToken: string) => {
          const newHeaders = new Headers(options.headers);
          newHeaders.set('Authorization', `Bearer ${newToken}`);
          resolve(fetch(input, { ...options, headers: newHeaders }));
        });
      });
    }
  }

  return res;
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

  const res = await apiFetch(url, options);
  return handleResponse<T>(res, defaultError);
}
