package status

import (
	"context"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func collectTopProcesses() []ProcessInfo {
	if runtime.GOOS != "darwin" {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	out, err := runCmd(ctx, "ps", "-Aceo", "pid=,ppid=,pcpu=,pmem=,rss=,comm=", "-r")
	if err != nil {
		return nil
	}

	lines := strings.Split(strings.TrimSpace(out), "\n")
	var procs []ProcessInfo
	for i, line := range lines {
		if i >= 20 {
			break
		}
		fields := strings.Fields(line)
		if len(fields) < 6 {
			continue
		}
		pid, _ := strconv.Atoi(fields[0])
		ppid, _ := strconv.Atoi(fields[1])
		cpuVal, _ := strconv.ParseFloat(fields[2], 64)
		rssKB, _ := strconv.ParseFloat(fields[4], 64)
		command := strings.Join(fields[5:], " ")
		name := command
		if idx := strings.LastIndex(name, "/"); idx >= 0 {
			name = name[idx+1:]
		}
		procs = append(procs, ProcessInfo{
			PID:     pid,
			PPID:    ppid,
			Name:    name,
			Command: command,
			CPU:     cpuVal,
			Memory:  rssKB / 1024,
		})
	}
	return procs
}
