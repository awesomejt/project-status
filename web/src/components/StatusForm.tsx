import { useState, useEffect } from "react";
import { useNavigate, useParams, useSearchParams } from "react-router-dom";
import { apiClient } from "../api/client";
import type { StatusRecordCreate, StatusRecordUpdate, ApiError } from "../types/statusRecord";

const StatusForm = () => {
  const navigate = useNavigate();
  const { id } = useParams<{ id: string }>();
  const [searchParams] = useSearchParams();
  
  const isEditing = Boolean(id);
  const source = searchParams.get("source") || "web";
  
  const [formData, setFormData] = useState<StatusRecordCreate>({
    project_name: "",
    short_name: "",
    status: "active",
    phase: "",
    summary: "",
    reason: "",
    details: "",
    tags: [],
    source,
  });
  
  const [tagInput, setTagInput] = useState("");
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<ApiError | null>(null);
  const [validationErrors, setValidationErrors] = useState<Record<string, string>>({});
  
  const peek = searchParams.get("peek") === "true";
  
  const resetForm = () => {
    setFormData({
      project_name: "",
      short_name: "",
      status: "active",
      phase: "",
      summary: "",
      reason: "",
      details: "",
      tags: [],
      source,
    });
    setTagInput("");
    setValidationErrors({});
    setError(null);
  };
  
  useEffect(() => {
    if (isEditing && id) {
      const fetchRecord = async () => {
        setLoading(true);
        try {
          const record = await apiClient.getRecord(id);
          setFormData({
            project_name: record.project_name,
            short_name: record.short_name,
            status: record.status,
            phase: record.phase || "",
            summary: record.summary || "",
            reason: record.reason || "",
            details: record.details || "",
            tags: record.tags || [],
            source: source || "web",
          });
        } catch (err) {
          setError(err as ApiError);
        } finally {
          setLoading(false);
        }
      };
      fetchRecord();
    } else {
      setLoading(false);
      resetForm();
    }
  }, [id, source]);
  
  const validateForm = (): boolean => {
    const errors: Record<string, string> = {};
    
    if (!formData.project_name || formData.project_name.trim().length === 0) {
      errors.project_name = "Project name is required";
    } else if (formData.project_name.length > 255) {
      errors.project_name = "Project name must be 255 characters or less";
    }
    
    if (!formData.short_name || formData.short_name.trim().length === 0) {
      errors.short_name = "Short name is required";
    } else if (formData.short_name.length > 64) {
      errors.short_name = "Short name must be 64 characters or less";
    }
    
    if (formData.phase && formData.phase.length > 64) {
      errors.phase = "Phase must be 64 characters or less";
    }
    
    if (formData.summary && formData.summary.length > 512) {
      errors.summary = "Summary must be 512 characters or less";
    }
    
    if (formData.reason && formData.reason.length > 1024) {
      errors.reason = "Reason must be 1024 characters or less";
    }
    
    if (formData.details && formData.details.length > 4096) {
      errors.details = "Details must be 4096 characters or less";
    }
    
    setValidationErrors(errors);
    return Object.keys(errors).length === 0;
  };
  
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!validateForm()) {
      return;
    }
    
    setSubmitting(true);
    setError(null);
    
    try {
      if (isEditing && id) {
        const updateData: StatusRecordUpdate = {
          project_name: formData.project_name,
          short_name: formData.short_name,
          status: formData.status,
          phase: formData.phase || undefined,
          summary: formData.summary || undefined,
          reason: formData.reason || undefined,
          details: formData.details || undefined,
          tags: formData.tags && formData.tags.length > 0 ? formData.tags : undefined,
        };
        await apiClient.updateRecord(id, updateData);
      } else {
        await apiClient.createRecord(formData);
      }
      navigate("/");
    } catch (err) {
      setError(err as ApiError);
    } finally {
      setSubmitting(false);
    }
  };
  
  const handleTagAdd = () => {
    const tag = tagInput.trim();
    if (tag && !formData.tags?.includes(tag)) {
      setFormData((prev) => ({
        ...prev,
        tags: [...(prev.tags || []), tag],
      }));
      setTagInput("");
    }
  };
  
  const handleTagRemove = (tag: string) => {
    setFormData((prev) => ({
      ...prev,
      tags: prev.tags?.filter((t) => t !== tag) || [],
    }));
  };
  
  const handleKeyDownTag = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === "Enter") {
      e.preventDefault();
      handleTagAdd();
    } else if (e.key === "Backspace" && tagInput === "") {
      if (formData.tags && formData.tags.length > 0) {
        handleTagRemove(formData.tags[formData.tags.length - 1]);
      }
    }
  };
  
  const handleCancel = () => {
    navigate("/");
  };
  
  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement | HTMLTextAreaElement>) => {
    const { name, value } = e.target;
    setFormData((prev) => ({
      ...prev,
      [name]: value,
    }));
    if (validationErrors[name]) {
      setValidationErrors((prev) => {
        const newErrors = { ...prev };
        delete newErrors[name];
        return newErrors;
      });
    }
  };
  
  if (loading) {
    return (
      <div style={{ maxWidth: "800px", margin: "0 auto", padding: "20px" }}>
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
          <style dangerouslySetInnerHTML={{ __html: `@keyframes spin { 0% { transform: rotate(0deg); } 100% { transform: rotate(360deg); } }` }} />
        </div>
      </div>
    );
  }
  
  return (
    <div style={{ maxWidth: "800px", margin: "0 auto", padding: "20px" }}>
      <header style={{ marginBottom: "24px" }}>
        <h1 id="form-title" style={{ fontSize: "28px", marginBottom: "8px" }}>
          {isEditing ? "Edit Status Record" : peek ? "Review Status Record" : "Create Status Record"}
        </h1>
        <p id="form-description" style={{ color: "#666", fontSize: "14px" }}>
          {isEditing ? "Update an existing status record" : "Add a new project status record"}
          {!(isEditing || peek) && " Fields marked with an asterisk (*) are required."}
        </p>
      </header>
      
      {error && (
        <div
          style={{
            padding: "16px",
            backgroundColor: "#fef2f2",
            border: "1px solid #fecaca",
            borderRadius: "8px",
            color: "#dc2626",
            marginBottom: "20px",
          }}
        >
          <strong>Error:</strong> {error.error.message}
          {error.error.code && ` (Code: ${error.error.code})`}
        </div>
      )}
      
      <form onSubmit={handleSubmit} aria-labelledby="form-title" aria-describedby="form-description">
        <div style={{ marginBottom: "20px" }}>
          <label htmlFor="project_name" style={{ display: "block", marginBottom: "6px", fontWeight: 500 }}>
            Project Name <span style={{ color: "#dc2626" }}>*</span>
          </label>
          <input
            type="text"
            id="project_name"
            name="project_name"
            value={formData.project_name}
            onChange={handleChange}
            disabled={peek}
            style={{
              width: "100%",
              padding: "10px 12px",
              border: validationErrors.project_name ? "2px solid #dc2626" : "1px solid #d1d5db",
              borderRadius: "6px",
              fontSize: "14px",
              backgroundColor: peek ? "#f3f4f6" : "white",
            }}
            placeholder="Enter project name"
            maxLength={255}
          />
          {validationErrors.project_name && (
            <span style={{ color: "#dc2626", fontSize: "12px", marginTop: "4px", display: "block" }}>
              {validationErrors.project_name}
            </span>
          )}
        </div>
        
        <div style={{ marginBottom: "20px" }}>
          <label htmlFor="short_name" style={{ display: "block", marginBottom: "6px", fontWeight: 500 }}>
            Short Name <span style={{ color: "#dc2626" }}>*</span>
          </label>
          <input
            type="text"
            id="short_name"
            name="short_name"
            value={formData.short_name}
            onChange={handleChange}
            disabled={peek}
            style={{
              width: "100%",
              padding: "10px 12px",
              border: validationErrors.short_name ? "2px solid #dc2626" : "1px solid #d1d5db",
              borderRadius: "6px",
              fontSize: "14px",
              backgroundColor: peek ? "#f3f4f6" : "white",
            }}
            placeholder="Enter short name"
            maxLength={64}
          />
          {validationErrors.short_name && (
            <span style={{ color: "#dc2626", fontSize: "12px", marginTop: "4px", display: "block" }}>
              {validationErrors.short_name}
            </span>
          )}
        </div>
        
        <div style={{ marginBottom: "20px" }}>
          <label htmlFor="status" style={{ display: "block", marginBottom: "6px", fontWeight: 500 }}>
            Status
          </label>
          <select
            id="status"
            name="status"
            value={formData.status}
            onChange={handleChange}
            disabled={peek}
            style={{
              width: "100%",
              padding: "10px 12px",
              border: "1px solid #d1d5db",
              borderRadius: "6px",
              fontSize: "14px",
              backgroundColor: peek ? "#f3f4f6" : "white",
            }}
          >
            <option value="active">Active</option>
            <option value="paused">Paused</option>
            <option value="blocked">Blocked</option>
            <option value="working">Working</option>
            <option value="error">Error</option>
            <option value="stopped">Stopped</option>
            <option value="completed">Completed</option>
          </select>
        </div>
        
        <div style={{ marginBottom: "20px" }}>
          <label htmlFor="phase" style={{ display: "block", marginBottom: "6px", fontWeight: 500 }}>
            Phase
          </label>
          <input
            type="text"
            id="phase"
            name="phase"
            value={formData.phase}
            onChange={handleChange}
            disabled={peek}
            style={{
              width: "100%",
              padding: "10px 12px",
              border: validationErrors.phase ? "2px solid #dc2626" : "1px solid #d1d5db",
              borderRadius: "6px",
              fontSize: "14px",
              backgroundColor: peek ? "#f3f4f6" : "white",
            }}
            placeholder="e.g., planning, implementation, validation, release"
            maxLength={64}
          />
          {validationErrors.phase && (
            <span style={{ color: "#dc2626", fontSize: "12px", marginTop: "4px", display: "block" }}>
              {validationErrors.phase}
            </span>
          )}
        </div>
        
        <div style={{ marginBottom: "20px" }}>
          <label htmlFor="summary" style={{ display: "block", marginBottom: "6px", fontWeight: 500 }}>
            Summary
          </label>
          <textarea
            id="summary"
            name="summary"
            value={formData.summary}
            onChange={handleChange}
            disabled={peek}
            rows={3}
            style={{
              width: "100%",
              padding: "10px 12px",
              border: validationErrors.summary ? "2px solid #dc2626" : "1px solid #d1d5db",
              borderRadius: "6px",
              fontSize: "14px",
              resize: "vertical",
              backgroundColor: peek ? "#f3f4f6" : "white",
            }}
            placeholder="Short user-facing status summary"
            maxLength={512}
          />
          {validationErrors.summary && (
            <span style={{ color: "#dc2626", fontSize: "12px", marginTop: "4px", display: "block" }}>
              {validationErrors.summary}
            </span>
          )}
        </div>
        
        <div style={{ marginBottom: "20px" }}>
          <label htmlFor="reason" style={{ display: "block", marginBottom: "6px", fontWeight: 500 }}>
            Reason
          </label>
          <textarea
            id="reason"
            name="reason"
            value={formData.reason}
            onChange={handleChange}
            disabled={peek}
            rows={3}
            style={{
              width: "100%",
              padding: "10px 12px",
              border: validationErrors.reason ? "2px solid #dc2626" : "1px solid #d1d5db",
              borderRadius: "6px",
              fontSize: "14px",
              resize: "vertical",
              backgroundColor: peek ? "#f3f4f6" : "white",
            }}
            placeholder="Explanation for paused, blocked, error, or stopped states"
            maxLength={1024}
          />
          {validationErrors.reason && (
            <span style={{ color: "#dc2626", fontSize: "12px", marginTop: "4px", display: "block" }}>
              {validationErrors.reason}
            </span>
          )}
        </div>
        
        <div style={{ marginBottom: "20px" }}>
          <label htmlFor="details" style={{ display: "block", marginBottom: "6px", fontWeight: 500 }}>
            Details
          </label>
          <textarea
            id="details"
            name="details"
            value={formData.details}
            onChange={handleChange}
            disabled={peek}
            rows={6}
            style={{
              width: "100%",
              padding: "10px 12px",
              border: validationErrors.details ? "2px solid #dc2626" : "1px solid #d1d5db",
              borderRadius: "6px",
              fontSize: "14px",
              resize: "vertical",
              backgroundColor: peek ? "#f3f4f6" : "white",
            }}
            placeholder="Longer notes and details"
            maxLength={4096}
          />
          {validationErrors.details && (
            <span style={{ color: "#dc2626", fontSize: "12px", marginTop: "4px", display: "block" }}>
              {validationErrors.details}
            </span>
          )}
        </div>
        
        <div style={{ marginBottom: "24px" }}>
          <div>
            <label htmlFor="tags-input" style={{ display: "block", marginBottom: "6px", fontWeight: 500 }}>
              Tags
            </label>
            <span style={{ fontSize: "12px", color: "#6b7280" }}>
              Type a tag and press Enter or click Add. Press Backspace to remove the last tag.
            </span>
          </div>
          <div style={{ display: "flex", gap: "8px", marginBottom: "8px" }}>
            <input
              type="text"
              id="tags-input"
              aria-describedby="tags-instructions"
              value={tagInput}
              onChange={(e) => setTagInput(e.target.value)}
              onKeyDown={handleKeyDownTag}
              disabled={peek}
              style={{
                flex: 1,
                padding: "10px 12px",
                border: "1px solid #d1d5db",
                borderRadius: "6px",
                fontSize: "14px",
                backgroundColor: peek ? "#f3f4f6" : "white",
              }}
              placeholder="Type a tag and press Enter"
            />
            <button
              type="button"
              onClick={handleTagAdd}
              aria-label="Add tag"
              disabled={peek || !tagInput.trim()}
              style={{
                padding: "10px 20px",
                backgroundColor: peek ? "#9ca3af" : "#3b82f6",
                color: "white",
                border: "none",
                borderRadius: "6px",
                fontSize: "14px",
                cursor: peek ? "not-allowed" : "pointer",
              }}
            >
              Add
            </button>
          </div>
          <span id="tags-instructions" style={{ display: "none" }}>Type a tag and press Enter or click Add. Press Backspace to remove the last tag.</span>
          {formData.tags && formData.tags.length > 0 && (
            <div style={{ display: "flex", flexWrap: "wrap", gap: "6px" }}>
              {formData.tags.map((tag) => (
                <span
                  key={tag}
                  style={{
                    display: "inline-flex",
                    alignItems: "center",
                    padding: "4px 8px",
                    backgroundColor: "#e5e7eb",
                    borderRadius: "12px",
                    fontSize: "12px",
                    fontWeight: 500,
                  }}
                >
                  {tag}
                   {!peek && (
                     <button
                       type="button"
                       onClick={() => handleTagRemove(tag)}
                       aria-label={`Remove tag "${tag}"`}
                       style={{
                         marginLeft: "6px",
                         border: "none",
                         background: "none",
                         cursor: "pointer",
                         color: "#6b7280",
                         fontSize: "14px",
                         fontWeight: 600,
                         padding: "0",
                         lineHeight: 1,
                         width: "16px",
                         height: "16px",
                         display: "flex",
                         alignItems: "center",
                         justifyContent: "center",
                       }}
                     >
                       ×
                     </button>
                   )}
                </span>
              ))}
            </div>
          )}
        </div>
        
        <div style={{ display: "flex", gap: "12px", justifyContent: "flex-end" }}>
          <button
            type="button"
            onClick={handleCancel}
            style={{
              padding: "10px 20px",
              backgroundColor: "white",
              color: "#374151",
              border: "1px solid #d1d5db",
              borderRadius: "6px",
              fontSize: "14px",
              cursor: "pointer",
            }}
          >
            Cancel
          </button>
          {!peek && (
            <button
              type="submit"
              disabled={submitting}
              style={{
                padding: "10px 20px",
                backgroundColor: "#22c55e",
                color: "white",
                border: "none",
                borderRadius: "6px",
                fontSize: "14px",
                cursor: submitting ? "not-allowed" : "pointer",
                opacity: submitting ? 0.7 : 1,
              }}
            >
              {submitting ? (isEditing ? "Updating..." : "Creating...") : isEditing ? "Update Record" : "Create Record"}
            </button>
          )}
        </div>
      </form>
    </div>
  );
};

export default StatusForm;
