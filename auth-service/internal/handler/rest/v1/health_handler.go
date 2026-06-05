package v1

import (
	"auth-service/internal/model"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"gorm.io/gorm"
)

type HealthHandler struct {
	StartAt time.Time
	Db      *gorm.DB
}

func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{StartAt: time.Now(),
		Db: db,
	}
}

func (h *HealthHandler) Health1(c *gin.Context) {
	//Memory
	vmStat, err := mem.VirtualMemory()
	totalMemory := uint64(0)
	freeMemory := uint64(0)
	if err == nil {
		totalMemory = vmStat.Total / 1024 / 1024
		freeMemory = vmStat.Free / 1024 / 1024
	}
	// cpu - get CPU number of cores and speed
	percentage, err := cpu.Percent(0, true)
	var arrayCPUs []string
	for idx, percent := range percentage {
		arrayCPUs = append(arrayCPUs, "CPU["+strconv.Itoa(idx)+"]: "+strconv.FormatFloat(percent, 'f', 2, 64)+"%")
	}
	// host or machine kernel, uptime, platform Info
	hostStat, err := host.Info()
	hostOS := ""
	hostID := ""
	if err == nil {
		hostOS = hostStat.OS
		hostID = hostStat.HostID
	}
	response := model.HealthCheckResponse{
		Name:        "go-template",
		Uptime:      time.Now().Sub(h.StartAt).String(),
		TotalMemory: strconv.FormatUint(totalMemory, 10) + " MB",
		FreeMemory:  strconv.FormatUint(freeMemory, 10) + " MB",
		UsedPercent: strconv.FormatFloat(vmStat.UsedPercent, 'f', 2, 64) + "%",
		Cpus:        arrayCPUs,
		HostOS:      hostOS,
		HostId:      hostID,
	}
	c.JSON(http.StatusOK, response)
}

func (h *HealthHandler) Health(c *gin.Context) {
	health := gin.H{
		"status":    "healthy",
		"service":   "image-service",
		"timestamp": time.Now().Format(time.RFC3339),
	}

	// Check database connection
	if h.Db != nil {
		sqlDB, err := h.Db.DB()
		if err != nil {
			health["status"] = "unhealthy"
			health["database"] = gin.H{
				"status": "error",
				"error":  err.Error(),
			}
			c.JSON(http.StatusServiceUnavailable, health)
			return
		}

		if err := sqlDB.Ping(); err != nil {
			health["status"] = "unhealthy"
			health["database"] = gin.H{
				"status": "disconnected",
				"error":  err.Error(),
			}
			c.JSON(http.StatusServiceUnavailable, health)
			return
		}

		// Get database stats
		stats := sqlDB.Stats()
		health["database"] = gin.H{
			"status":           "connected",
			"open_connections": stats.OpenConnections,
			"in_use":           stats.InUse,
			"idle":             stats.Idle,
		}
	}

	c.JSON(http.StatusOK, health)
}
