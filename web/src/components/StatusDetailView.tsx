import { useState, useEffect } from "react";
import { useNavigate, useParams } from "react-router-dom";
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

const StatusDetailView = () => {
  const navigate = useNavigate();
  const { id } = useParams<{ id: string }>();
  const [record, setRecord] = useState<StatusRecord | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<ApiError | null>(null);
  const [deleting, setDeleting] = useState(false);
  const [showDeleteConfirm, setShowDeleteConfirm] = useState(false);

  useEffect(() => {
    if (id) {
      fetchRecord();
    }
  }, [id]);

  const fetchRecord = async () => {
    if (!id) return;
    
    setLoading(true);
    setError(null);
    
    try {
      const data = await apiClient.getRecord(id);
      setRecord(data as StatusRecord);
    } catch (err) {
      setError(err as ApiError);
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async () => {
    if (!id) return;
    
    setDeleting(true);
    try {
      await apiClient.deleteRecord(id);
      navigate("/");
    } catch (err) {
      setError(err as ApiError);
      setDeleting(false);
      setShowDeleteConfirm(false);
    }
  };

  const handleEdit = () => {
    if (id) {
      navigate(`/edit/${id}`);
    }
  };

  const formatDateTime = (isoString: string) => {
    return new Date(isoString).toLocaleString();
  };

  const handleCancelDelete = () => {
    setShowDeleteConfirm(false);
  };

  if (loading) {
    return (
      <div style={{ maxWidth: "900px", margin: "0 auto", padding: "20px" }}>
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
          <p>Loading status record...</p>
          <style
            dangerouslySetInnerHTML={{
              __html: `@keyframes spin { 0% { transform: rotate(0deg); } 100% { transform: rotate(360deg); } }`,
            }}
          />
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div style={{ maxWidth: "900px", margin: "0 auto", padding: "20px" }}>
        <div style={{ marginBottom: "16px" }}>
          <button
            onClick={() => navigate("/")}
            style={{
              padding: "8px 16px",
              border: "1px solid #d1d5db",
              borderRadius: "6px",
              backgroundColor: "white",
              color: "#374151",
              fontSize: "14px",
              cursor: "pointer",
            }}
          >
            ← Back to List
          </button>
        </div>
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
      </div>
    );
  }

  if (!record) {
    return (
      <div style={{ maxWidth: "900px", margin: "0 auto", padding: "20px" }}>
        <div style={{ textAlign: "center", padding: "60px 20px", backgroundColor: "#f9fafb", borderRadius: "8px" }}>
          <h3 style={{ fontSize: "18px", marginBottom: "8px" }}>Status record not found</h3>
          <p style={{ color: "#666", marginBottom: "16px" }}>The requested status record does not exist or has been deleted.</p>
          <button
            onClick={() => navigate("/")}
            style={{
              padding: "10px 20px",
              backgroundColor: "#3b82f6",
              color: "white",
              border: "none",
              borderRadius: "6px",
              fontSize: "14px",
              cursor: "pointer",
            }}
          >
            Back to List
          </button>
        </div>
      </div>
    );
  }

  return (
    <div style={{ maxWidth: "900px", margin: "0 auto", padding: "20px" }}>
      <div style={{ marginBottom: "24px", display: "flex", justifyContent: "space-between", alignItems: "center" }}>
        <button
          onClick={() => navigate("/")}
          style={{
            padding: "8px 16px",
            border: "1px solid #d1d5db",
            borderRadius: "6px",
            backgroundColor: "white",
            color: "#374151",
            fontSize: "14px",
            cursor: "pointer",
          }}
        >
          ← Back to List
        </button>
        <div style={{ display: "flex", gap: "8px" }}>
          {!showDeleteConfirm && (
            <>
              <button
                onClick={handleEdit}
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
                Edit Record
              </button>
              <button
                onClick={() => setShowDeleteConfirm(true)}
                style={{
                  padding: "10px 20px",
                  backgroundColor: "#ef4444",
                  color: "white",
                  border: "none",
                  borderRadius: "6px",
                  fontSize: "14px",
                  cursor: "pointer",
                  fontWeight: 500,
                }}
              >
                Delete
              </button>
            </>
          )}
        </div>
      </div>

      {showDeleteConfirm && (
        <div
          style={{
            marginBottom: "20px",
            padding: "16px",
            backgroundColor: "#fff7ed",
            border: "1px solid #fed7aa",
            borderRadius: "8px",
          }}
        >
          <p style={{ fontWeight: 600, marginBottom: "8px" }}>
            Are you sure you want to delete "{record.project_name}"?
          </p>
          <p style={{ color: "#666", fontSize: "14px", marginBottom: "16px" }}>
            This action cannot be undone.
          </p>
          <div style={{ display: "flex", gap: "8px" }}>
            <button
              onClick={handleDelete}
              disabled={deleting}
              style={{
                padding: "8px 16px",
                backgroundColor: "#dc2626",
                color: "white",
                border: "none",
                borderRadius: "6px",
                fontSize: "14px",
                cursor: deleting ? "not-allowed" : "pointer",
                opacity: deleting ? 0.7 : 1,
              }}
            >
              {deleting ? "Deleting..." : "Yes, Delete"}
            </button>
            <button
              onClick={handleCancelDelete}
              disabled={deleting}
              style={{
                padding: "8px 16px",
                backgroundColor: "white",
                color: "#374151",
                border: "1px solid #d1d5db",
                borderRadius: "6px",
                fontSize: "14px",
                cursor: deleting ? "not-allowed" : "pointer",
              }}
            >
              Cancel
            </button>
          </div>
        </div>
      )}

      <section>
        <div style={{ marginBottom: "8px" }}>
          <StatusBadge status={record.status} />
        </div>
        <h1 style={{ fontSize: "32px", marginBottom: "8px" }}>{record.project_name}</h1>
        <p style={{ color: "#6b7280", fontSize: "16px" }}>{record.short_name}</p>
      </section>

      <section
        style={{
          marginTop: "32px",
          padding: "24px",
          backgroundColor: "#f9fafb",
          borderRadius: "8px",
        }}
      >
        <div
          style={{
            display: "grid",
            gridTemplateColumns: "repeat(auto-fit, minmax(250px, 1fr))",
            gap: "20px",
          }}
        >
          <div>
            <div style={{ fontSize: "12px", color: "#6b7280", textTransform: "uppercase", fontWeight: 600, marginBottom: "4px" }}>
              Phase
            </div>
            <div style={{ fontSize: "16px", color: "#111827" }}>{record.phase || "—"}</div>
          </div>
          <div>
            <div style={{ fontSize: "12px", color: "#6b7280", textTransform: "uppercase", fontWeight: 600, marginBottom: "4px" }}>
              Source
            </div>
            <div style={{ fontSize: "16px", color: "#111827" }}>{record.source || "—"}</div>
          </div>
          <div>
            <div style={{ fontSize: "12px", color: "#6b7280", textTransform: "uppercase", fontWeight: 600, marginBottom: "4px" }}>
              Created
            </div>
            <div style={{ fontSize: "16px", color: "#111827" }}>{formatDateTime(record.created_at)}</div>
          </div>
          <div>
            <div style={{ fontSize: "12px", color: "#6b7280", textTransform: "uppercase", fontWeight: 600, marginBottom: "4px" }}>
              Last Updated
            </div>
            <div style={{ fontSize: "16px", color: "#111827" }}>{formatDateTime(record.updated_at)}</div>
          </div>
        </div>
      </section>

      <section style={{ marginTop: "24px" }}>
        <h2 style={{ fontSize: "18px", fontWeight: 600, marginBottom: "12px" }}>Summary</h2>
        <div style={{ fontSize: "16px", color: "#374151", lineHeight: 1.6 }}>
          {record.summary || "No summary provided."}
        </div>
      </section>

      {record.reason && (
        <section style={{ marginTop: "24px" }}>
          <h2 style={{ fontSize: "18px", fontWeight: 600, marginBottom: "12px" }}>Reason</h2>
          <div style={{ fontSize: "16px", color: "#374151", lineHeight: 1.6 }}>
            {record.reason}
          </div>
        </section>
      )}

      {record.details && (
        <section style={{ marginTop: "24px" }}>
          <h2 style={{ fontSize: "18px", fontWeight: 600, marginBottom: "12px" }}>Details</h2>
          <div
            style={{
              fontSize: "16px",
              color: "#374151",
              lineHeight: 1.6,
              whiteSpace: "pre-wrap",
            }}
          >
            {record.details}
          </div>
        </section>
      )}

      {record.tags && record.tags.length > 0 && (
        <section style={{ marginTop: "24px" }}>
          <h2 style={{ fontSize: "18px", fontWeight: 600, marginBottom: "12px" }}>Tags</h2>
          <div style={{ display: "flex", flexWrap: "wrap", gap: "8px" }}>
            {record.tags.map((tag) => (
              <span
                key={tag}
                style={{
                  padding: "4px 12px",
                  backgroundColor: "#e5e7eb",
                  color: "#374151",
                  borderRadius: "16px",
                  fontSize: "13px",
                  fontWeight: 500,
                }}
              >
                {tag}
              </span>
            ))}
          </div>
        </section>
      )}
    </div>
  );
};

export default StatusDetailView;
