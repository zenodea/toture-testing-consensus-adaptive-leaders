package util

import (
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"time"
)

type Performance struct {
	Option map[string]string
}

func NewBasicPerformance() *Performance {
	basic_stats := make(map[string]string)
	basic_stats["throughput"] = ""
	basic_stats["median_latency"] = ""
	basic_stats["p99_latency"] = ""
	basic_stats["average_latency"] = ""
	return &Performance{
		Option: basic_stats,
	}
}

func NewPerformance() Performance {
	return Performance{Option: make(map[string]string)}
}

func NewPerformanceWithOptions(options map[string]string) Performance {
	return Performance{
		Option: options,
	}
}

// retrieve the current CPU usage percentage.

func GetCPUUsage() float64 {
	percentages, err := cpu.Percent(time.Second, false)
	if err != nil {
		panic(err.Error())
	}
	if len(percentages) == 0 {
		panic("no CPU usage data")
	}
	return percentages[0]
}

// retrieve the current memory usage percentage.

func GetMemoryUsage() float64 {
	vm, err := mem.VirtualMemory()
	if err != nil {
		panic(err.Error())
	}
	return vm.UsedPercent
}

// retrieve the network packets in/out rate.

func GetNetworkStats() (float64, float64) {
	initialStats, err := net.IOCounters(false)
	if err != nil {
		panic(err.Error())
	}

	time.Sleep(100 * time.Millisecond)

	finalStats, err := net.IOCounters(false)

	if err != nil {
		panic(err.Error())
	}

	if len(initialStats) == 0 || len(finalStats) == 0 {
		panic("no network stats data")
	}

	// Calculate packet rates (packets per second)
	packetsInRate := float64(finalStats[0].PacketsRecv-initialStats[0].PacketsRecv) * 10
	packetsOutRate := float64(finalStats[0].PacketsSent-initialStats[0].PacketsSent) * 10

	return packetsInRate, packetsOutRate
}
