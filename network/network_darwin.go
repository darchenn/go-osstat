// +build darwin

package network

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"
)

// Get network statistics
func Get() ([]NetworkStats, error) {
	cmd := exec.Command("netstat", "-bni")
	out, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	networks, err := collectNetworkStats(out)
	if err != nil {
		return nil, err
	}
	if err := cmd.Wait(); err != nil {
		return nil, err
	}
	return networks, nil
}

// NetworkStats represents network statistics for darwin
type NetworkStats struct {
	Name             string
	RxBytes, TxBytes uint64
}

func collectNetworkStats(out io.Reader) ([]NetworkStats, error) {
	scanner := bufio.NewScanner(out)
	if !scanner.Scan() { // skip the first line
		return nil, fmt.Errorf("failed to scan output of netstat")
	}

	var networks []NetworkStats
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		name := strings.TrimSuffix(fields[0], "*")
		if strings.HasPrefix(name, "lo") || !strings.HasPrefix(fields[2], "<Link#") {
			continue
		}
		rxBytesIdx, txBytesIdx := 6, 9
		if len(fields) == 10 {
			rxBytesIdx, txBytesIdx = 5, 8
		}
		rxBytes, err := strconv.ParseUint(fields[rxBytesIdx], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse Ibytes of %s", name)
		}
		txBytes, err := strconv.ParseUint(fields[txBytesIdx], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse Obytes of %s", name)
		}
		networks = append(networks, NetworkStats{Name: name, RxBytes: rxBytes, TxBytes: txBytes})
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan error for netstat: %s", err)
	}

	return networks, nil
}
