package commands

import (
	"TuruBot-Go/internal/app/types"
	"fmt"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/process"
	"os"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/load"
	"github.com/shirou/gopsutil/v4/net"
)

func (cmd *Command) StatusHandler(ctx *types.BotContext) error {
	uptime, err := host.Uptime()
	if err != nil {
		return ctx.Reply(fmt.Sprintf("Failed to get uptime: %v", err))
	}
	uptimeDur := time.Duration(uptime) * time.Second
	uptimeString := formatDurationHuman(uptimeDur)

	avg, _ := load.Avg() // ignore error dan tampilkan kosong jika error
	memStat, _ := getAppMemoryUsageMB()
	cpuPercent, _ := cpu.Percent(0, false)
	cputimes, _ := cpu.Times(false)
	diskUsage, _ := disk.Usage("/")
	netIO, _ := net.IOCounters(false)

	var sb strings.Builder
	sb.WriteString("*System Host Status*\n\n")
	sb.WriteString(fmt.Sprintf("Uptime: %v\n", uptimeString))

	if avg != nil {
		sb.WriteString(fmt.Sprintf("Load Average: %.2f (1m), %.2f (5m), %.2f (15m)\n", avg.Load1, avg.Load5, avg.Load15))
	} else {
		sb.WriteString("Load Average: N/A\n")
	}

	if memStat != "" {
		sb.WriteString(fmt.Sprintf("Memory Usage: %s\n", memStat))
	} else {
		sb.WriteString("Memory Usage: N/A\n")
	}

	if len(cpuPercent) > 0 {
		sb.WriteString(fmt.Sprintf("CPU Usage: %.2f%%\n", cpuPercent[0]))
	} else {
		sb.WriteString("CPU Usage: N/A\n")
	}

	if len(cputimes) > 0 {
		cpuTotal := cputimes[0]
		sb.WriteString(fmt.Sprintf("CPU Times (s): user=%.2f system=%.2f idle=%.2f\n", cpuTotal.User, cpuTotal.System, cpuTotal.Idle))
	}

	if diskUsage != nil {
		sb.WriteString(fmt.Sprintf("Disk Usage (/): %.2f%% (%v / %v)\n", diskUsage.UsedPercent, formatBytes(diskUsage.Used), formatBytes(diskUsage.Total)))
	} else {
		sb.WriteString("Disk Usage (/): N/A\n")
	}

	if len(netIO) > 0 {
		io := netIO[0]
		sb.WriteString(fmt.Sprintf("Network I/O: sent=%v received=%v\n", formatBytes(io.BytesSent), formatBytes(io.BytesRecv)))
	}

	return ctx.Reply(sb.String())
}

func getAppMemoryUsageMB() (string, error) {
	pid := int32(os.Getpid())
	p, err := process.NewProcess(pid)
	if err != nil {
		return "", err
	}
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return "", err
	}

	memInfo, err := p.MemoryInfo()
	if err != nil {
		return "", err
	}

	memPercent, err := p.MemoryPercent()
	if err != nil {
		return "", err
	}

	vmStatTotal := formatBytes(vmStat.Total)
	rss := formatBytes(memInfo.RSS)

	return fmt.Sprintf(
		"%.2f%% (%v/%v)",
		memPercent,
		rss,
		vmStatTotal,
	), nil
}

func formatBytes(b uint64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := uint64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}

func formatDurationHuman(d time.Duration) string {
	seconds := int(d.Seconds())
	days := seconds / 86400
	seconds %= 86400
	hours := seconds / 3600
	seconds %= 3600
	minutes := seconds / 60
	seconds %= 60

	parts := []string{}
	if days > 0 {
		parts = append(parts, fmt.Sprintf("%d hari", days))
	}
	if hours > 0 {
		parts = append(parts, fmt.Sprintf("%d jam", hours))
	}
	if minutes > 0 {
		parts = append(parts, fmt.Sprintf("%d menit", minutes))
	}
	if seconds > 0 {
		parts = append(parts, fmt.Sprintf("%d detik", seconds))
	}
	if len(parts) == 0 {
		return "0 detik"
	}
	return strings.Join(parts, " ")
}
