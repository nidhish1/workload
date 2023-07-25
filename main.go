package main

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	
)

func getCPUTemperature() (float64, error) {
	cmd := exec.Command("cat", "/sys/class/thermal/thermal_zone0/temp")
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	tempStr := strings.TrimSpace(string(output))
	temp, err := strconv.ParseFloat(tempStr, 64)
	if err != nil {
		return 0, err
	}

	return temp / 1000.0, nil
}

func getGPUTemperature() (float64, error) {
	cmd := exec.Command("nvidia-smi", "--query-gpu=temperature.gpu", "--format=csv,noheader,nounits")
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	tempStr := strings.TrimSpace(string(output))
	temp, err := strconv.ParseFloat(tempStr, 64)
	if err != nil {
		return 0, err
	}

	return temp, nil
}

func getCPUWorkload() (float64, error) {
	cmd := exec.Command("mpstat", "1", "1")
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	// Parse the CPU usage percentage from the output
	lines := strings.Split(string(output), "\n")
	if len(lines) < 4 {
		return 0, fmt.Errorf("unexpected mpstat output format")
	}

	fields := strings.Fields(lines[len(lines)-2])
	if len(fields) < 12 {
		return 0, fmt.Errorf("unexpected mpstat output format")
	}

	workloadStr := fields[11]
	workload, err := strconv.ParseFloat(workloadStr, 64)
	if err != nil {
		return 0, err
	}

	return 100.0 - workload, nil
}

func getGPUWorkload() (float64, error) {
	cmd := exec.Command("nvidia-smi", "--query-gpu=utilization.gpu", "--format=csv,noheader,nounits")
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	gpuUtilStr := strings.TrimSpace(string(output))
	gpuUtil, err := strconv.ParseFloat(gpuUtilStr, 64)
	if err != nil {
		return 0, err
	}

	return gpuUtil, nil
}

func main() {
	for {
		cpuTemp, err := getCPUTemperature()
		if err != nil {
			fmt.Println("Failed to read CPU temperature:", err)
		} else {
			fmt.Printf("CPU Temperature: %.2f °C\n", cpuTemp)
		}

		gpuTemp, err := getGPUTemperature()
		if err != nil {
			fmt.Println("Failed to read GPU temperature:", err)
		} else {
			fmt.Printf("GPU Temperature: %.2f °C\n", gpuTemp)
		}

		cpuWorkload, err := getCPUWorkload()
		if err != nil {
			fmt.Println("Failed to read CPU workload:", err)
		} else {
			fmt.Printf("CPU Workload: %.2f%%\n", cpuWorkload)
		}

		gpuWorkload, err := getGPUWorkload()
		if err != nil {
			fmt.Println("Failed to read GPU workload:", err)
		} else {
			fmt.Printf("GPU Workload: %.2f%%\n", gpuWorkload)
		}

		fmt.Println("--------------")
		time.Sleep(5 * time.Second)
	}
}
