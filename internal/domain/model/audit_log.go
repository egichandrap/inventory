package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// AuditAction represents the type of audit action
type AuditAction string

const (
	ActionCreate AuditAction = "CREATE"
	ActionUpdate AuditAction = "UPDATE"
	ActionDelete AuditAction = "DELETE"
	ActionLogin AuditAction = "LOGIN"
	ActionLogout AuditAction = "LOGOUT"
	ActionCheckout AuditAction = "CHECKOUT"
	ActionCancel AuditAction = "CANCEL"
	ActionRefund AuditAction = "REFUND"
	ActionAdjustStock AuditAction = "ADJUST_STOCK"
	ActionChangePassword AuditAction = "CHANGE_PASSWORD"
)

// AuditLog represents an audit log entry
type AuditLog struct {
	id         string
	timestamp  time.Time
	userID     string
	userName   string
	action     AuditAction
	entityType string
	entityID   string
	details    map[string]interface{}
	ipAddress  string
	userAgent  string
	success    bool
	errorMsg   string
}

// NewAuditLog creates a new audit log entry
func NewAuditLog(
	userID, userName string,
	action AuditAction,
	entityType, entityID string,
	details map[string]interface{},
	ipAddress, userAgent string,
	success bool,
) *AuditLog {
	return &AuditLog{
		id:         uuid.New().String(),
		timestamp:  time.Now(),
		userID:     userID,
		userName:   userName,
		action:     action,
		entityType: entityType,
		entityID:   entityID,
		details:    details,
		ipAddress:  ipAddress,
		userAgent:  userAgent,
		success:    success,
	}
}

// Accessors
func (a *AuditLog) ID() string            { return a.id }
func (a *AuditLog) Timestamp() time.Time  { return a.timestamp }
func (a *AuditLog) UserID() string        { return a.userID }
func (a *AuditLog) UserName() string      { return a.userName }
func (a *AuditLog) Action() AuditAction   { return a.action }
func (a *AuditLog) EntityType() string    { return a.entityType }
func (a *AuditLog) EntityID() string      { return a.entityID }
func (a *AuditLog) Details() map[string]interface{} { return a.details }
func (a *AuditLog) IPAddress() string     { return a.ipAddress }
func (a *AuditLog) UserAgent() string     { return a.userAgent }
func (a *AuditLog) Success() bool         { return a.success }
func (a *AuditLog) ErrorMsg() string      { return a.errorMsg }

// SetError sets the error message
func (a *AuditLog) SetError(errMsg string) {
	a.errorMsg = errMsg
	a.success = false
}

// MarshalDetailsJSON marshals details to JSON
func (a *AuditLog) MarshalDetailsJSON() (string, error) {
	if a.details == nil {
		return "{}", nil
	}
	data, err := json.Marshal(a.details)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// UnmarshalDetailsJSON unmarshals JSON to details
func (a *AuditLog) UnmarshalDetailsJSON(jsonStr string) error {
	if jsonStr == "" || jsonStr == "{}" {
		a.details = make(map[string]interface{})
		return nil
	}
	return json.Unmarshal([]byte(jsonStr), &a.details)
}
