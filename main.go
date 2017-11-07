package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

var lf, listen *string

func init() {
	lf = flag.String("lf", "/var/lib/dhcpd/dhcpd.leases", "Full path to dhcpd.leases")
	listen = flag.String("listen", ":8080", "Listen address")
	flag.Parse()
}

func main() {
	http.HandleFunc("/dhcp/stats", stats)
	log.Fatal(http.ListenAndServe(*listen, nil))
}

func stats(w http.ResponseWriter, r *http.Request) {

	ip, ipnet, err := net.ParseCIDR(r.URL.Query().Get("subnet"))

	if err != nil {
		fmt.Fprintln(w, "Invalid Subnet")
		return
	}

	rbody := "<html><head><title>Subnets</title></head><body>%v Total: %d Active: %d  Free: %d Utilization Percent: %.0f</body></html>"

	if r.Header.Get("Accept") == "text/csv" {
		rbody = "%v,%d,%d,%d,%.0f"
	}

	leases := lcount(ip, ipnet)
	addresses := size(ipnet)
	fmt.Fprintf(w,
		rbody,
		ip,
		addresses,
		leases,
		(addresses - leases),
		pfree(leases, addresses))

}

func pfree(i, j int) float64 {
	return float64(i) / float64(j) * 100
}

func size(n *net.IPNet) int {
	mask, bits := n.Mask.Size()
	netmask := 0xffffffff << uint32(bits - mask) 
	host := 0xffffffff &^ netmask 
	return int(host) 
}

func lcount(i net.IP, n *net.IPNet) int {
	f, err := os.Open(*lf)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	r := bufio.NewScanner(f)
	addrs := make(map[string]bool)

	var count int
	for r.Scan() {
		if strings.Contains(r.Text(), "lease") {

			l := strings.Fields(r.Text())[1]
			ip := net.ParseIP(l)

			if n.Contains(ip) {
				for r.Scan() {

					if err := r.Err(); err != nil {
						log.Printf("File scanner error: %v", err)
					}

					line := r.Text()
					if lact(line) {
						addrs[l] = true
						break
					} else if lfree(line) {
						addrs[l] = false
						break
					}

				}
			}
		}

	}

	if err := r.Err(); err != nil {
		log.Printf("File scanner error: %v", err)
	}

	for k := range addrs {
		if addrs[k] {
			count++
		}
	}

	return count
}

func lact(s string) bool {
	return !strings.Contains(s, "next") &&
		strings.Contains(s, "binding state active")
}

func lfree(s string) bool {
	return !strings.Contains(s, "next") &&
		(strings.Contains(s, "binding state free") ||
			strings.Contains(s, "binding state backup"))
}
