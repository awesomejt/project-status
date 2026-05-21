export type StatusValue =
  | "active"
  | "paused"
  | "blocked"
  | "working"
  | "error"
  | "stopped"
  | "completed";

export interface StatusRecord {
  id: string;
  project_name: string;
  short_name: string;
  status: StatusValue;
  phase?: string;
  summary?: string;
  reason?: string;
  details?: string;
  tags?: string[];
  source?: string;
  created_at: string;
  updated_at: string;
}

export interface StatusRecordCreate {
  project_name: string;
  short_name: string;
  status: StatusValue;
  phase?: string;
  summary?: string;
  reason?: string;
  details?: string;
  tags?: string[];
  source?: string;
}

export interface StatusRecordUpdate {
  project_name?: string;
  short_name?: string;
  status?: StatusValue;
  phase?: string;
  summary?: string;
  reason?: string;
  details?: string;
  tags?: string[];
}

export interface StatusRecordListResponse {
  records: Omit<StatusRecord, "reason" | "details" | "tags" | "source">[];
  page: number;
  per_page: number;
  total: number;
  pages: number;
}

export interface ApiError {
  error: {
    code: number;
    message: string;
    details?: unknown;
  };
}
