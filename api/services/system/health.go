package system

import (
	"encoding/json"
	"github.com/MrJoshLab/go-respond"
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"math"
	"time"
)

func (s *Service) GetHealthCharts(ctx *gin.Context)  {

	cpuInfo, _ := cpu.Info()
	cpuUsedPercent, _ := cpu.Percent(time.Second, false)

	memory, _ := mem.VirtualMemory()
	var memoryData map[string] interface{}
	_ = json.Unmarshal([]byte(memory.String()), &memoryData)

	loadAvg, _ := load.Avg()
	diskUsage, _ := disk.Usage("/")

	data := map[string] interface{} {
		"load": map[string] interface{} {
			"avg": loadAvg,
		},
		"cpu": map[string] interface{} {
			"name":              "NaN",
			"cores":             0,
			"cache_size":        0,
			"used_percentage":   math.Round(cpuUsedPercent[0]),
		},
		"memory": memoryData,
		"disk": map[string] interface{} {
			"usage": diskUsage,
		},
	}

	if len(cpuInfo) == 1 {
		data["cpu"].(map[string] interface{})["name"] = cpuInfo[0].ModelName
		data["cpu"].(map[string] interface{})["cores"] = cpuInfo[0].Cores
		data["cpu"].(map[string] interface{})["cache_size"] = cpuInfo[0].CacheSize
	}

	ctx.JSON(respond.Default.Succeed(data))
}
