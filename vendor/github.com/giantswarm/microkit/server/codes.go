package server

var (
	// CodeFailure indicates the requested action failed.
	CodeFailure = "FAILURE"
	// CodeInvalidCredentials indicates the provided credentials are not valid.
	CodeInvalidCredentials = "INVALID_CREDENTIALS"
	// CodePermissionDenied indicates the provided credentials are valid, but the
	// requested resource requires other permissions.
	CodePermissionDenied = "PERMISSION_DENIED"
	// CodeResourceAlreadyExists indicates a resource does already exist.
	CodeResourceAlreadyExists = "RESOURCE_ALREADY_EXISTS"
	// CodeResourceCreated indicates a resource has been created.
	CodeResourceCreated = "RESOURCE_CREATED"
	// CodeResourceDeleted indicates a resource has been deleted.
	CodeResourceDeleted = "RESOURCE_DELETED"
	// CodeResourceDeletionStarted indicates a resource will be deleted.
	CodeResourceDeletionStarted = "RESOURCE_DELETION_STARTED"
	// CodeResourceNotFound indicates a resource could not be found.
	CodeResourceNotFound = "RESOURCE_NOT_FOUND"
	// CodeResourceUpdated indicates a resource has been updated.
	CodeResourceUpdated = "RESOURCE_UPDATED"
	// CodeSuccess indicates the requested action successed.
	CodeSuccess = "SUCCESS"
	// CodeImmutableAttribute indicates the provided data structure contains
	// fields that are immutable.
	CodeImmutableAttribute = "IMMUTABLE_ATTRIBUTE"
	// CodeUnknownAttribute indicates the provided data structure contains
	// unexpected fields.
	CodeUnknownAttribute = "UNKNOWN_ATTRIBUTE"
	// CodeInternalError represents an error we don't want to give more details
	// about (usually HTTP status 500).
	CodeInternalError = "INTERNAL_ERROR"
)
