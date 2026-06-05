import type { ApiResponse } from '@/types/api';

const BASE_URL = '/api/v1';

class ApiClient {
  private async request<T>(path: string, options: RequestInit = {}): Promise<T> {
    var response = await fetch(`${BASE_URL}${path}`, {
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json',
        ...options.headers,
      },
      ...options,
    });

    var json: ApiResponse<T> = await response.json();

    if (!response.ok || json.error) {
      throw new Error(json.error || `Request failed: ${response.status}`);
    }

    return json.data as T;
  }

  get<T>(path: string) {
    return this.request<T>(path);
  }

  post<T>(path: string, body?: unknown) {
    return this.request<T>(path, {
      method: 'POST',
      body: body ? JSON.stringify(body) : undefined,
    });
  }

  put<T>(path: string, body?: unknown) {
    return this.request<T>(path, {
      method: 'PUT',
      body: body ? JSON.stringify(body) : undefined,
    });
  }

  delete<T>(path: string) {
    return this.request<T>(path, { method: 'DELETE' });
  }
}

export const api = new ApiClient();
