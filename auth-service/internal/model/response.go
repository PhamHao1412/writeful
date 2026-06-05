package model

type HealthCheckResponse struct {
	Name        string   `json:"name"`
	Uptime      string   `json:"uptime"`
	TotalMemory string   `json:"total_memory"`
	FreeMemory  string   `json:"free_memory"`
	UsedPercent string   `json:"used_memory"`
	Cpus        []string `json:"cpus"`
	HostOS      string   `json:"host_os"`
	HostId      string   `json:"host_id"`
}

type GetUserInfoResponse struct {
	ID       string   `json:"id"`
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Status   string   `json:"status"`
	AvtarURL string   `json:"avatar_url"`
	Roles    []string `json:"roles"`
}

type PaginationResponse struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalPages int   `json:"total_pages"`
	TotalItems int64 `json:"total_items"`
}

type GetUsersResponse struct {
	Users      []GetUserInfoResponse `json:"users"`
	Pagination PaginationResponse    `json:"pagination"`
}
