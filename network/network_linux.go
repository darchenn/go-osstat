// +build linux

package network

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// Get network statistics
func Get() ([]NetworkStats, error) {
	file, err := os.Open("/proc/net/dev")
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return collectNetworkStats(file)
}

// NetworkStats represents network statistics for linux
type NetworkStats struct {
	Name             string
	RxBytes, TxBytes uint64
}

func collectNetworkStats(out io.Reader) ([]NetworkStats, error) {
	scanner := bufio.NewScanner(out)
	var networks []NetworkStats
	rxBytesIdx, txBytesIdx := 1, 9
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 17 || len(fields) > 0 && !strings.HasSuffix(fields[0], ":") {
			continue
		}
		name := fields[0][:len(fields[0])-1]
		if name == "lo" {
			continue
		}
		rxBytes, err := strconv.ParseUint(fields[rxBytesIdx], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse rxBytes of %s", name)
		}
		txBytes, err := strconv.ParseUint(fields[txBytesIdx], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse txBytes of %s", name)
		}
		networks = append(networks, NetworkStats{Name: name, RxBytes: rxBytes, TxBytes: txBytes})
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan error for /proc/net/dev: %s", err)
	}

	return networks, nil
}
