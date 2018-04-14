package main

import (
	"fmt"
	"net"
	"time"
	"strings"
	"net/http"
	"github.com/miekg/dns"
)

func main() {
	addrs, _ := net.InterfaceAddrs()
	for _, addr := range addrs {
		ipnet := addr.(*net.IPNet)
		if ipnet.IP.To4() == nil { continue }
		if ipnet.IP.IsGlobalUnicast() == false { continue }
		fmt.Println("specify this as DNS server:", ipnet.IP.String())
	}
	
	dns.HandleFunc(".", recursor)
	dns.HandleFunc("nintendo.com.", self)
	dns.HandleFunc("nintendo.net.", self)
	dns.HandleFunc("nintendowifi.net.", self)
	
	go func () {
		server := &dns.Server{Addr: ":53", Net: "udp"}
		if err := server.ListenAndServe(); err != nil { panic(err) }
	}()
	
	http.HandleFunc("/", jumper)
	
	go func () {
		server := &http.Server{Addr: ":80"}
		if err := server.ListenAndServe(); err != nil { panic(err) }
	}()
	
	select {}
}

func self(writer dns.ResponseWriter, req *dns.Msg) {
	dump(writer, req)
	
	res := new(dns.Msg)
	res.SetReply(req)
	res.Compress = true
	res.RecursionAvailable = true
	
	qname := req.Question[0].Name
	
	addrs, _ := net.InterfaceAddrs()
	for _, addr := range addrs {
		ipnet := addr.(*net.IPNet)
		
		if ipnet.IP.To4() == nil { continue }
		if ipnet.IP.IsGlobalUnicast() == false { continue }
		
		rr := new(dns.A)
		rr.Hdr = dns.RR_Header{
			Name: qname,
			Rrtype: dns.TypeA,
			Class: dns.ClassINET,
			Ttl: 0,
		}
		rr.A = ipnet.IP
		res.Answer = append(res.Answer, rr)
	}
	
	writer.WriteMsg(res)
}

func recursor(writer dns.ResponseWriter, req *dns.Msg) {
	dump(writer, req)
	
	res := new(dns.Msg)
	res.SetReply(req)
	res.Compress = true
	res.RecursionAvailable = true
	
	qname := req.Question[0].Name
	
	// I don't care IPv6-only servers
	
	addrs, err := net.LookupIP(qname)
	if err != nil {
		res.SetRcode(req, dns.RcodeServerFailure)
		writer.WriteMsg(res)
		return
	}
	
	for _, addr := range addrs {
		if addr.To4() == nil { continue }
		
		rr := new(dns.A)
		rr.Hdr = dns.RR_Header{
			Name: qname,
			Rrtype: dns.TypeA,
			Class: dns.ClassINET,
			Ttl: 0,
		}
		rr.A = addr
		res.Answer = append(res.Answer, rr)
		
		break
	}
	
	writer.WriteMsg(res)
}

func dump(writer dns.ResponseWriter, msg *dns.Msg) {
	var cols []string
	
	now := time.Now()
	cols = append(cols, now.Format("2006-01-02 15:04:05"))
	
	cols = append(cols, msg.Question[0].Name)
	cols = append(cols, dns.TypeToString[msg.Question[0].Qtype])
	//cols = append(cols, fmt.Sprintf("%d", msg.MsgHdr.Id))
	
	//cols = append(cols, writer.RemoteAddr().Network())
	cols = append(cols, writer.RemoteAddr().String())
	
	fmt.Println(strings.Join(cols, "\t"))
}

func jumper(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/html")
	res.Write([]byte(`<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<meta name="viewport" content="initial-scale=1">
</head>
<body>
<input type="url" value="https://" style="width: 100%;"><br>
<input type="submit" onclick="window.location = document.querySelector('input').value;">
</body>
</html>
`))
}
