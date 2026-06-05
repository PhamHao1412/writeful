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

type UploadResult struct {
	PublicId     string `json:"public_id"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	Format       string `json:"format"`
	ResourceType string `json:"resource_type"`
	URL          string `json:"url"`
	SecureURL    string `json:"secure_url"`
	Bytes        int    `json:"bytes"`
}
