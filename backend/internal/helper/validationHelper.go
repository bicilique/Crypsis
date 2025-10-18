package helper

// ValidateSortParams validates sortBy and order parameters to prevent SQL injection
// Returns validated sortBy and order, using defaults if invalid values are provided
func ValidateSortParams(sortBy, order string, allowedFields map[string]bool) (string, string) {
	// Validate and set default for sortBy
	if sortBy == "" || !allowedFields[sortBy] {
		sortBy = "created_at"
	}

	// Validate and set default for order
	if order != "asc" && order != "desc" {
		order = "desc"
	}

	return sortBy, order
}

// Common allowed sort fields for different entities
var (
	AllowedFileSortFields = map[string]bool{
		"created_at": true,
		"name":       true,
		"size":       true,
	}

	AllowedAppSortFields = map[string]bool{
		"created_at": true,
		"name":       true,
		"client_id":  true,
	}

	AllowedAdminSortFields = map[string]bool{
		"id":         true,
		"username":   true,
		"client_id":  true,
		"created_at": true,
		"updated_at": true,
	}

	AllowedLogSortFields = map[string]bool{
		"created_at": true,
		"file_id":    true,
		"action":     true,
		"timestamp":  true,
		"ip":         true,
		"user_agent": true,
	}
)
