package cmd

import (
	"fmt"
	"net"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/go-ping/ping"
)

//执行ping命令
func Ping(ip string, count int, timeout time.Duration) string {
	strBuf := strings.Builder{}
	pinger, err := ping.NewPinger(ip)
	if err != nil {
		strBuf.WriteString(err.Error())
		return strBuf.String()
	}
	pinger.Timeout = timeout
	pinger.Count = count
	if runtime.GOOS == "windows" {
		pinger.SetPrivileged(true)
	}
	pinger.OnRecv = func(pkt *ping.Packet) {
		fmt.Fprintf(&strBuf, "%d bytes from %s: icmp_seq=%d time=%v\n",
			pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt)
	}

	pinger.OnDuplicateRecv = func(pkt *ping.Packet) {
		fmt.Fprintf(&strBuf, "%d bytes from %s: icmp_seq=%d time=%v ttl=%v (DUP!)\n",
			pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt, pkt.Ttl)
	}

	pinger.OnFinish = func(stats *ping.Statistics) {
		fmt.Fprintf(&strBuf, "\n--- %s ping statistics ---\n", stats.Addr)
		fmt.Fprintf(&strBuf, "%d packets transmitted, %d packets received, %v%% packet loss\n",
			stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
		fmt.Fprintf(&strBuf, "round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
			stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)
	}

	fmt.Fprintf(&strBuf, "PING %s (%s):\n", pinger.Addr(), pinger.IPAddr())
	err = pinger.Run()
	if err != nil {
		strBuf.WriteString(err.Error())
		return strBuf.String()
	}

	return strBuf.String()
}

func IsPortAvailable(ip string, port int) bool {
	conn, err := net.Dial("tcp", ip+":"+strconv.Itoa(port))
	if err != nil {
		return false
	}

	conn.SetDeadline(time.Now().Add(time.Second * 10))
	defer conn.Close()
	return true
}
