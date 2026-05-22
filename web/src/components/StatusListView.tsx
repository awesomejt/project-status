import { useState, useEffect, useCallback } from "react";
import { useNavigate } from "react-router-dom";
import { apiClient } from "../api/client";
import type { StatusRecord, StatusValue, ApiError } from "../types/statusRecord";

const StatusBadge = ({ status }: { status: StatusValue }) => {
  const statusColors: Record<StatusValue, string> = {
    active: "#22c55e",
    paused: "#eab308",
    blocked: "#ef4444",
    working: "#3b82f6",
    error: "#dc2626",
    stopped: "#6b7280",
    completed: "#16a34a",
  };

  return (
    <span
      style={{
        padding: "2px 8px",
        borderRadius: "12px",
        backgroundColor: statusColors[status],
        color: "white",
        fontSize: "12px",
        fontWeight: 600,
        textTransform: "uppercase",
      }}
    >
      {status}
    </span>
  );
};

const StatusListView = () => {
  const navigate = useNavigate();
  const [records, setRecords] = useState<StatusRecord[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<ApiError | null>(null);
  const [filterStatus, setFilterStatus] = useState<StatusValue | "all">("all");
  const [pagination, setPagination] = useState({
    page: 1,
    per_page: 20,
    total: 0,
    pages: 0,
  });

  const fetchRecords = useCallback(async (page = 1) => {
    setLoading(true);
    setError(null);

    try {
      const response = await apiClient.getRecords({
        page,
        per_page: 20,
        status: filterStatus === "all" ? undefined : filterStatus,
      });

      setRecords(response.records);
      setPagination({
        page: response.page,
        per_page: response.per_page,
        total: response.total,
        pages: response.pages,
      });
    } catch (err) {
      const apiError = err as ApiError;
      setError(apiError);
    } finally {
      setLoading(false);
    }
  }, [filterStatus]);

  useEffect(() => {
    fetchRecords(1);
  }, [fetchRecords]);

  const handlePreviousPage = () => {
    if (pagination.page > 1) {
      fetchRecords(pagination.page - 1);
    }
  };

  const handleNextPage = () => {
    if (pagination.page < pagination.pages) {
      fetchRecords(pagination.page + 1);
    }
  };

  const formatDateTime = (isoString: string) => {
    return new Date(isoString).toLocaleString();
  };

  return (
    <div style={{ maxWidth: "1200px", margin: "0 auto", padding: "20px" }}>
      <header style={{ marginBottom: "24px", display: "flex", justifyContent: "space-between", alignItems: "center" }}>
        <div>
          <h1 style={{ fontSize: "28px", marginBottom: "8px" }}>Project Status</h1>
          <p style={{ color: "#666", fontSize: "14px" }}>
            Track and manage project status across your organization
          </p>
        </div>
        <button
          onClick={() => navigate("/create")}
          aria-label="Create a new status record"
          style={{
            padding: "10px 20px",
            backgroundColor: "#3b82f6",
            color: "white",
            border: "none",
            borderRadius: "6px",
            fontSize: "14px",
            cursor: "pointer",
            fontWeight: 500,
          }}
        >
          + Create Status Record
        </button>
      </header>

      <section style={{ marginBottom: "20px" }}>
        <label htmlFor="status-filter" style={{ marginRight: "8px", fontWeight: 500 }}>Filter:</label>
        <select
          id="status-filter"
          value={filterStatus}
          onChange={(e) => {
            setFilterStatus(e.target.value as StatusValue | "all");
            fetchRecords(1);
          }}
          aria-label="Filter status records by status"
          style={{
            padding: "8px 12px",
            borderRadius: "6px",
            border: "1px solid #ccc",
            fontSize: "14px",
          }}
        >
          <option value="all">All Status</option>
          <option value="active">Active</option>
          <option value="paused">Paused</option>
          <option value="blocked">Blocked</option>
          <option value="working">Working</option>
          <option value="error">Error</option>
          <option value="stopped">Stopped</option>
          <option value="completed">Completed</option>
        </select>
      </section>

      {loading ? (
        <section aria-live="polite">
          <div style={{ textAlign: "center", padding: "40px", color: "#666" }}>
            <div
              style={{
                width: "40px",
                height: "40px",
                border: "4px solid #f3f3f3",
                borderTop: "4px solid #3b82f6",
                borderRadius: "50%",
                animation: "spin 1s linear infinite",
                margin: "0 auto 16px",
              }}
            ></div>
            <p>Loading status records...</p>
          </div>
          <style>
            {
              `
              @keyframes spin {
                0% { transform: rotate(0deg); }
                100% { transform: rotate(360deg); }
              }
            `
            }
          </style>
        </section>
      ) : error ? (
        <section aria-live="assertive">
          <div
            style={{
              padding: "20px",
              backgroundColor: "#fef2f2",
              border: "1px solid #fecaca",
              borderRadius: "8px",
              color: "#dc2626",
            }}
          >
            <strong>Error:</strong> {error.error.message}
            {error.error.code && ` (Code: ${error.error.code})`}
          </div>
        </section>
      ) : records.length === 0 ? (
        <section aria-live="polite">
          <div
            style={{
              textAlign: "center",
              padding: "60px 20px",
              backgroundColor: "#f9fafb",
              borderRadius: "8px",
            }}
          >
            <svg
              style={{ width: "64px", height: "64px", marginBottom: "16px", opacity: 0.5 }}
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={1.5}
                d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-6 9l2 2 4-4"
              />
            </svg>
            <h3 style={{ fontSize: "18px", marginBottom: "8px" }}>No status records found</h3>
            <p style={{ color: "#666" }}>
              {filterStatus !== "all"
                ? `No records with status "${filterStatus}". Try a different filter.`
                : "Click \"Create Status Record\" to add your first status record."}
            </p>
          </div>
        </section>
      ) : (
        <section>
          <div
            style={{
              overflowX: "auto",
              border: "1px solid #e5e7eb",
              borderRadius: "8px",
            }}
          >
            <table
              aria-label="Status records table showing project name, status, phase, summary, and last updated date"
              style={{
                width: "100%",
                borderCollapse: "collapse",
                fontSize: "14px",
              }}
           >
              <thead>
                <tr
                  style={{
                    backgroundColor: "#f9fafb",
                    textAlign: "left",
                  }}
                >
                  <th style={headerStyle}>Project</th>
                  <th style={headerStyle}>Status</th>
                  <th style={headerStyle}>Phase</th>
                  <th style={headerStyle}>Summary</th>
                  <th style={headerStyle}>Updated</th>
                </tr>
              </thead>
              <tbody>
                {records.map((record) => (
                  <tr
                    key={record.id}
                    role="row"
                    tabIndex={0}
                    style={{
                      borderBottom: "1px solid #e5e7eb",
                      cursor: "pointer",
                      outline: "none",
                    }}
                    onClick={() => navigate(`/detail/${record.id}`)}
                    onKeyDown={(e) => {
                      if (e.key === 'Enter' || e.key === ' ') {
                        e.preventDefault();
                        navigate(`/detail/${record.id}`);
                      }
                    }}
                    onMouseEnter={(e) => (e.currentTarget.style.backgroundColor = "#f9fafb")}
                    onMouseLeave={(e) => (e.currentTarget.style.backgroundColor = "transparent")}
                    onFocus={(e) => (e.currentTarget.style.backgroundColor = "#f9fafb")}
                    onBlur={(e) => (e.currentTarget.style.backgroundColor = "transparent")}
                  >
                    <td style={cellStyle}>
                      <div style={{ textDecoration: "underline", textUnderlineOffset: "4px" }}>
                        <strong>{record.project_name}</strong>
                      </div>
                      <div style={{ fontSize: "12px", color: "#6b7280" }}>
                        {record.short_name}
                      </div>
                    </td>
                    <td style={cellStyle}>
                      <StatusBadge status={record.status} />
                    </td>
                    <td style={cellStyle}>{record.phase || "—"}</td>
                    <td style={{ ...cellStyle, maxWidth: "250px" }}>
                      <span
                        title={record.summary || ""}
                        style={{
                          display: "block",
                          overflow: "hidden",
                          textOverflow: "ellipsis",
                          whiteSpace: "nowrap",
                        }}
                      >
                        {record.summary || "—"}
                      </span>
                    </td>
                    <td style={{ ...cellStyle, color: "#6b7280" }}>
                      {formatDateTime(record.updated_at)}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          {pagination.pages > 1 && (
            <div
              style={{
                display: "flex",
                justifyContent: "space-between",
                alignItems: "center",
                marginTop: "16px",
                padding: "0 8px",
              }}
            >
              <div style={{ fontSize: "14px", color: "#6b7280" }}>
                Page {pagination.page} of {pagination.pages} ({pagination.total} records)
              </div>
              <div style={{ display: "flex", gap: "8px" }}>
                <button
                  onClick={handlePreviousPage}
                  disabled={pagination.page <= 1}
                  style={{
                    ...buttonStyle,
                    opacity: pagination.page <= 1 ? 0.5 : 1,
                    cursor: pagination.page <= 1 ? "not-allowed" : "pointer",
                  }}
                >
                  Previous
                </button>
                <button
                  onClick={handleNextPage}
                  disabled={pagination.page >= pagination.pages}
                  style={{
                    ...buttonStyle,
                    opacity: pagination.page >= pagination.pages ? 0.5 : 1,
                    cursor: pagination.page >= pagination.pages ? "not-allowed" : "pointer",
                  }}
                >
                  Next
                </button>
              </div>
            </div>
          )}
        </section>
      )}
    </div>
  );
};

const headerStyle = {
  padding: "12px 16px",
  fontWeight: 600,
  color: "#374151",
  textTransform: "uppercase" as const,
  fontSize: "11px",
  letterSpacing: "0.5px",
};

const cellStyle = {
  padding: "12px 16px",
  verticalAlign: "middle",
};

const buttonStyle = {
  padding: "8px 16px",
  border: "1px solid #d1d5db",
  borderRadius: "6px",
  backgroundColor: "white",
  color: "#374151",
  fontSize: "14px",
  transition: "all 0.2s",
};

export default StatusListView;
