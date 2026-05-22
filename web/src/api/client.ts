import type { StatusRecord, StatusRecordCreate, StatusRecordUpdate, StatusRecordListResponse, ApiError, StatusValue } from "../types/statusRecord";

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || "http://localhost:5000";

const API_STATUS_RECORDS_PATH = "/api/project/status";

interface RequestConfig extends RequestInit {
  data?: object;
}

const onRequestError = (error: Error): ApiError => ({
  error: {
    code: 0,
    message: error.message,
  },
});

const fetcher = async <T>(
  path: string,
  config: RequestConfig = {},
): Promise<T> => {
  const url = `${API_BASE_URL}${path}`;
  
  const headers: HeadersInit = {
    "Content-Type": "application/json",
    ...config.headers,
  };
  
  const options: RequestInit = {
    ...config,
    headers,
    body: config.data ? JSON.stringify(config.data) : undefined,
  };
  
  try {
    const response = await fetch(url, options);
    const body = await response.json();
    
    if (!response.ok) {
      return Promise.reject(body);
    }
    
    return body as T;
  } catch (error) {
    return Promise.reject(onRequestError(error instanceof Error ? error : new Error(String(error))));
  }
};

export const apiClient = {
  getRecords(params?: {
    page?: number;
    per_page?: number;
    status?: StatusValue;
  }): Promise<StatusRecordListResponse> {
    const searchParams = new URLSearchParams();
    
    if (params?.page) searchParams.set("page", String(params.page));
    if (params?.per_page) searchParams.set("per_page", String(params.per_page));
    if (params?.status) searchParams.set("status", params.status);
    
    return fetcher<StatusRecordListResponse>(
      `${API_STATUS_RECORDS_PATH}${searchParams.toString() ? `?${searchParams.toString()}` : ""}`,
      { method: "GET" },
    );
  },
  
  getRecord(id: string): Promise<StatusRecord> {
    return fetcher<StatusRecord>(`${API_STATUS_RECORDS_PATH}/${id}`, { method: "GET" });
  },
  
  createRecord(data: StatusRecordCreate): Promise<StatusRecord> {
    return fetcher<StatusRecord>(API_STATUS_RECORDS_PATH, {
      method: "POST",
      data,
    });
  },
  
  updateRecord(id: string, data: StatusRecordUpdate): Promise<StatusRecord> {
    return fetcher<StatusRecord>(`${API_STATUS_RECORDS_PATH}/${id}`, {
      method: "PATCH",
      data,
    });
  },
  
  deleteRecord(id: string): Promise<{ message: string }> {
    return fetcher<{ message: string }>(`${API_STATUS_RECORDS_PATH}/${id}`, {
      method: "DELETE",
    });
  },
};
