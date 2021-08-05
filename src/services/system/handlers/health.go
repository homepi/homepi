package handlers

import (
	"encoding/json"
	"math"
	"net/http"
	"time"

	"github.com/homepi/homepi/src/core"
	"github.com/mrjosh/respond.go"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
)

func HandleGetHealth(ctx *core.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			ctx.WrapAuthAdminHandler(getHealth(ctx)).ServeHTTP(w, r)
			return
		default:
			respond.NewWithWriter(w).MethodNotAllowed()
		}
	})
}

func getHealth(ctx *core.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cpuInfo, _ := cpu.Info()
		cpuUsedPercent, _ := cpu.Percent(time.Second, false)

		memory, _ := mem.VirtualMemory()
		var memoryData map[string]interface{}
		_ = json.Unmarshal([]byte(memory.String()), &memoryData)

		loadAvg, _ := load.Avg()
		diskUsage, _ := disk.Usage("/")

		data := map[string]interface{}{
			"load": map[string]interface{}{
				"avg": loadAvg,
			},
			"cpu": map[string]interface{}{
				"name":            "NaN",
				"cores":           0,
				"cache_size":      0,
				"used_percentage": math.Round(cpuUsedPercent[0]),
			},
			"memory": memoryData,
			"disk": map[string]interface{}{
				"usage": diskUsage,
			},
		}

		if len(cpuInfo) == 1 {
			data["cpu"].(map[string]interface{})["name"] = cpuInfo[0].ModelName
			data["cpu"].(map[string]interface{})["cores"] = cpuInfo[0].Cores
			data["cpu"].(map[string]interface{})["cache_size"] = cpuInfo[0].CacheSize
		}

		respond.NewWithWriter(w).Succeed(data)

	})
}
