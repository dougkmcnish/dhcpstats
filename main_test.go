package main

import (
	"net"
	"strings"
	"testing"
)

const testLf = `# The format of this file is documented in the dhcpd.leases(5) manual page.
# This lease file was written by isc-dhcp-4.2.5


failover peer "production-network" state {
  my state normal at 4 2017/11/02 15:51:58;
  partner state normal at 2 2015/09/15 18:47:53;
}
lease 192.168.104.79 {
  starts 3 2017/11/01 21:28:33;
  ends 3 2017/11/01 22:34:25;
  tstp 3 2017/11/01 22:34:25;
  tsfp 4 2017/11/02 00:01:29;
  atsfp 4 2017/11/02 00:01:29;
  cltt 3 2017/11/01 21:28:33;
  binding state active;
  next binding state free;
  hardware ethernet de:ad:be:ef:fe:ed;
}
lease 192.168.104.84 {
  starts 3 2017/11/01 23:01:28;
  ends 3 2017/11/01 23:01:28;
  tstp 0 2016/10/23 15:36:59;
  tsfp 0 2016/10/23 15:36:59;
  atsfp 0 2016/10/23 15:36:59;
  cltt 3 2017/11/01 23:01:28;
  binding state free;
}
lease 192.168.104.189 {
  starts 4 2017/11/02 02:40:42;
  ends 4 2017/11/02 02:41:45;
  tstp 5 2017/10/20 21:28:19;
  tsfp 4 2017/11/02 02:40:42;
  atsfp 4 2017/11/02 02:40:42;
  cltt 0 2017/10/08 03:00:52;
  binding state active;
  next binding state free;
  hardware ethernet de:ad:be:ef:fe:ed;
}
lease 192.168.104.156 {
  starts 4 2017/11/02 04:05:59;
  ends 4 2017/11/02 06:05:59;
  tstp 4 2017/11/02 06:05:59;
  tsfp 4 2017/11/02 07:05:59;
  atsfp 4 2017/11/02 07:05:59;
  cltt 4 2017/11/02 04:05:59;
  binding state backup;
  next binding state free;
  hardware ethernet de:ad:be:ef:fe:ed;
}
lease 192.168.104.148 {
  starts 4 2017/11/02 08:23:15;
  ends 4 2017/11/02 06:26:05;
  tstp 4 2017/11/02 07:23:28;
  tsfp 4 2017/11/02 08:23:15;
  atsfp 4 2017/11/02 08:23:15;
  cltt 4 2017/11/02 06:24:05;
  binding state free;
  hardware ethernet de:ad:be:ef:fe:ed;
}
lease 192.168.104.227 {
  starts 4 2017/11/02 13:02:49;
  ends 4 2017/11/02 06:28:07;
  tstp 4 2017/11/02 12:03:02;
  tsfp 4 2017/11/02 13:02:49;
  atsfp 4 2017/11/02 13:02:49;
  cltt 4 2017/11/02 06:26:07;
  binding state free;
  hardware ethernet de:ad:be:ef:fe:ed;
}`

func TestLact(t *testing.T) {
	actString := "binding state active;"
	nextString := "next binding state free;"

	if !lact(actString) {
		t.Errorf("%v should be true, got %v", actString, lact(actString))
	}

	if lact(nextString) {
		t.Errorf("%v should be false, got %v", nextString, lact(nextString))
	}

}

func TestLfree(t *testing.T) {
	freeString := "binding state free;"
	nextString := "next binding state free;"

	if !lfree(freeString) {
		t.Errorf("%v should be true, got %v", freeString, lfree(freeString))
	}

	if lfree(nextString) {
		t.Errorf("%v should be false, got %v", nextString, lfree(nextString))
	}
}

func TestNumHosts(t *testing.T) {

	_, ipnet, _ := net.ParseCIDR("192.168.1.10/22")

	if numHosts(ipnet) != 1022 {
		t.Errorf("expected 1022, got %d", numHosts(ipnet))
	}

	_, ipnet, _ = net.ParseCIDR("192.168.1.10/24")

	if numHosts(ipnet) != 254 {
		t.Errorf("expected 254, got %d", numHosts(ipnet))
	}

	_, ipnet, _ = net.ParseCIDR("192.168.1.10/19")

	if numHosts(ipnet) != 8190 {
		t.Errorf("expected 8190, got %d", numHosts(ipnet))
	}
}

func TestActCount(t *testing.T) {

	ip, ipnet, _ := net.ParseCIDR("192.168.104.0/24")

	f := strings.NewReader(testLf)

	actLeases := actCount(f, ip, ipnet)

	if actLeases != 2 {
		t.Errorf("Expected 2, got %d", actLeases)
	}
}
