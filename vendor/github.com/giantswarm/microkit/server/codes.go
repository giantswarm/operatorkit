package server

var (
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
	// CodeResourceNotFound indicates a resource could not be found.
	CodeResourceNotFound = "RESOURCE_NOT_FOUND"
	// CodeResourceUpdated indicates a resource has been updated.
	CodeResourceUpdated = "RESOURCE_UPDATED"
	// CodeUnknownAttribute indicates the provided data structure contains
	// unexpected fields.
	CodeUnknownAttribute = "UNKNOWN_ATTRIBUTE"
	// CodeUnknownError indicates something went wrong in unpredictable ways.
	CodeUnknownError = "UNKNOWN_ERROR"
)
