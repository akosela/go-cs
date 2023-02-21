// Copyright 2014 Andy Kosela.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Concurrent ssh client.
// cs is a program for concurrently executing local or remote commands
// on multiple hosts.  It is using OpenSSH for running remote commands.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strings"
	"time"
)

const (
	echo	= "/bin/echo"
	du      = "/usr/bin/du"
	host    = "/usr/bin/host"
	hostn	= "/bin/hostname"
	ifconf 	= "/sbin/ifconfig"
	lspci	= "/sbin/lspci"
	mpath	= "/sbin/multipath"
	nc      = "/usr/bin/nc"
	scp     = "/usr/bin/scp"
	ssh     = "/usr/bin/ssh"
	top     = "/usr/bin/top"
	uptime  = "/usr/bin/uptime"
	timeFmt = "02-Jan-2006 15:04:05"
)

func createFile(path string) *os.File {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0666)
	if err != nil {
		fmt.Println("cs:", err)
		os.Exit(1)
	}
	return file
}

func openFile(path string) *os.File {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("cs:", err)
		os.Exit(1)
	}
	return file
}

func readFile(file *os.File) []string {
	var s []string
	r := bufio.NewReader(file)
	for {
		line, _ := r.ReadString('\n')
		if strings.HasPrefix(line, "#") {
			continue
		} else if line == "" {
			break
		}
		s = append(s, line)
	}
	defer file.Close()
	return s
}

func checkArgs(nocmd int, copy, download, file, hostsfile string, tty *bool,
	argv []string) (string, []string) {
	var command string
	var hosts []string

	if nocmd == 1 {
		if hostsfile == "" {
			if len(argv) < 1 {
				flag.Usage()
			}
			command = ""
			hosts = argv[0:]
		} else {
			f := openFile(hostsfile)
			hosts = readFile(f)
			command = ""
		}

		return command, hosts
	}

	if file != "" || copy != "" || download != "" {
		if hostsfile == "" {
			if len(argv) < 1 {
				flag.Usage()
			}
			hosts = argv[0:]
		} else {
			f := openFile(hostsfile)
			hosts = readFile(f)
		}
		if file != "" {
			f := openFile(file)
			if *tty {
				c := strings.Join(readFile(f), "")
				command = "sudo sh <<'EOF'" + c + "EOF"
			} else {
				command = strings.Join(readFile(f), "")
			}
		}
	} else {
		if hostsfile == "" {
			if len(argv) < 2 {
				flag.Usage()
			}
			command = argv[0]
			hosts = argv[1:]
		} else {
			f := openFile(hostsfile)
			hosts = readFile(f)
			command = argv[0]
		}
	}

	return command, hosts
}

func exist(hostname, path string) string {
	if fi, err := os.Stat(path + "/" + hostname); err != nil {
		os.Mkdir(path+"/"+hostname, 0755)
		return path + "/" + hostname
	} else if fi.IsDir() {
		return path + "/" + hostname
	}
	os.Mkdir(path+"/"+hostname+".host", 0755)
	return path + "/" + hostname + ".host"
}

func run(command, hostname, id, login, path, port, timeout, copy, disku,
	download string, dd, cname, hba, io, ip, lcmd, lip, netcat, nmap, ns,
	mount, mx, one, png, recursive, soa, up, verbose1, verbose2, verbose3,
	top1, tr, tri, tty, uname, vm *bool, ddir int, f *os.File) string {

	hostname = strings.Trim(hostname, "\n")
	batchmode := "-oBatchMode=yes"
	strict := "-oStrictHostKeyChecking=no"
	tout := "-oConnectTimeout=" + timeout
	flag := "-"
	flag2 := ""

	if id != "" {
		id = "-i" + id
	} else {
		id = "-2"
	}

	if login != "" {
		login = "-l" + login
	} else {
		login = "-2"
	}

	if *verbose1 {
		flag = "-v"
	} else if *verbose2 {
		flag = "-vv"
	} else if *verbose3 {
		flag = "-vvv"
	}

	var cmd *exec.Cmd
	if copy != "" && *recursive {
		cmd = exec.Command(scp, login, id, flag+"r", "-P", port,
			batchmode, strict, tout, copy, hostname+":"+ path)
	} else if copy != "" {
		cmd = exec.Command(scp, login, id, "-P", port, batchmode,
			strict, tout, copy, hostname+":"+path)
	} else if download != "" && *recursive {
		if ddir == 1 {
			path = exist(hostname, path)
		}

		cmd = exec.Command(scp, login, id, flag+"r", "-P", port,
			batchmode, strict, tout, hostname+":"+download, path)
	} else if download != "" {
		if ddir == 1 {
			path = exist(hostname, path)
		}

		cmd = exec.Command(scp, login, id, "-P", port, batchmode,
			strict, tout, hostname+":"+download, path)
	} else if *cname {
		cmd = exec.Command(host, "-tcname", hostname)
	} else if *dd {
		c := "sudo dmidecode -s system-product-name " +
			"; sudo dmidecode -s system-serial-number"

		cmd = exec.Command(ssh, login, id, "-q", flag+"tt", "-p",
			port, batchmode, strict, tout, hostname, c)
	} else if disku != "" {
		c := "if [ `uname -s` == OpenBSD ]; then sudo " + du +
			" -akx " + disku + " |sort -rn 2>/dev/null |head -20" +
			"; else sudo " + du + " -amx " + disku +
			" |sort -rn 2>/dev/null |head -20; fi"

		cmd = exec.Command(ssh, login, id, "-q", flag+"tt", "-p",
			port, batchmode, strict, tout, hostname, c)
	} else if *hba {
		c := "if [ `uname -s` == Linux ]; then " + lspci + " |" +
			"grep -i fibre; grep -v zZzZ -H " +
			"/sys/class/fc_host/host*/port_name; else echo " +
			"'cs: your OS is not supported'; fi"

		cmd = exec.Command(ssh, login, id, "-p", port, batchmode,
			strict, tout, hostname, c)
	} else if *io {
		c := "if [ `uname -s` == Linux ]; then iostat -mx |grep Dev " +
			"|awk '{OFS=\"\t\" ; print $1, $4, $5, $6, $7, $10, " +
			"$12, $9}'; iostat -mx |grep -e dm- -e sd -e hd | " +
			"sort -rnk10 2>/dev/null |head -20 | " +
			"awk '{OFS=\"\t\" ; print $1, $4, $5, $6, $7, $10, " +
			"$12, $9}'; elif [ `uname -s` == FreeBSD ]; then " +
			"iostat -x | grep qlen ; iostat -x |grep -v device | " +
			"sort -rnk9 2>/dev/null |head -20; " +
			"else echo 'cs: your OS is not supported'; fi"

		cmd = exec.Command(ssh, login, id, "-p", port, batchmode,
			strict, tout, hostname, c)
	} else if *ip {
		cmd = exec.Command(ssh, login, id, "-p", port, batchmode,
			strict, tout, hostname, ifconf)
	} else if *lip {
		c := "if [ `uname -s` == Linux ]; then for e in `ls " +
			"/sys/class/fc_host/host*/issue_lip`; do echo $e; " +
			"echo '1' >$e; done; for e in `ls " +
			"/sys/class/scsi_host/host*/scan`; do echo $e; echo " +
			"'- - -' >$e; done; echo '# multipath -F'; " + mpath +
			" -F; echo '# multipath'; " + mpath +
			"; else echo 'cs: your OS is not supported'; fi"
		command = "sudo sh <<'EOF'" + "\n" + c + "\n" + "EOF"

		cmd = exec.Command(ssh, login, id, "-q", flag+"tt", "-p",
			port, batchmode, strict, tout, hostname, command)
	} else if *mount {
		cmd = exec.Command(ssh, login, id, "-p", port, batchmode,
			strict, tout, hostname, "mount")
	} else if *mx {
		cmd = exec.Command(host, "-tmx", hostname)
	} else if *ns {
		cmd = exec.Command(host, hostname)
	} else if *nmap {
		cmd = exec.Command("nmap", hostname)
	} else if *netcat {
		cmd = exec.Command(nc, "-w1", hostname, port)
	} else if *png {
		var c string
		if runtime.GOOS == "linux" {
			c = "ping -nc1 -s16 -W3 " + hostname + " |grep from"
		} else if runtime.GOOS == "freebsd" || runtime.GOOS == "darwin" {
			c = "ping -nc1 -s16 -W3000 " + hostname + " |grep from"
		} else {
			c = "ping -nc1 -s16 -w3 " + hostname + " |grep from"
		}
		cmd = exec.Command("/bin/sh", "-c", c)
	} else if *lcmd {
		scmd := strings.Split(command, " ")

		switch len(scmd) {
		case 1:
			cmd = exec.Command(scmd[0], hostname)
		case 2:
			cmd = exec.Command(scmd[0], scmd[1], hostname)
		case 3:
			cmd = exec.Command(scmd[0], scmd[1], scmd[2], hostname)
		case 4:
			cmd = exec.Command(scmd[0], scmd[1], scmd[2], scmd[3],
				hostname)
		case 5:
			cmd = exec.Command(scmd[0], scmd[1], scmd[2], scmd[3],
				scmd[4], hostname)
		}
	} else if *soa {
		cmd = exec.Command(host, "-tsoa", hostname)
	} else if *top1 {
		c := "if [ `uname -s` == Linux ]; then " + top +
			" -cbn1 |grep -ve '\\[' -e " + top + " 2>/dev/null" +
			" |head -20; else " + top + " -b 13; fi"

		cmd = exec.Command(ssh, login, id, "-p", port, batchmode,
			strict, tout, hostname, c)
	} else if *tr {
		cmd = exec.Command("traceroute", hostname)
	} else if *tri {
		cmd = exec.Command("sudo", "traceroute", "-I", hostname)
	} else if *uname {
		c := "if [ `uname -s` == Linux ]; then uname -a; " +
			"cat /etc/redhat-release 2>/dev/null; else uname -a; fi"

		cmd = exec.Command(ssh, login, id, "-p", port, batchmode,
			strict, tout, hostname, c)
	} else if *up {
		cmd = exec.Command(ssh, login, id, "-p", port, batchmode,
			strict, tout, hostname, uptime)
	} else if *vm {
		c := "if [ `uname -s` == Linux ]; then vmstat -SM; " +
			"elif [ `uname -s` == FreeBSD ]; then vmstat -h; " +
			"else vmstat; fi"

		cmd = exec.Command(ssh, login, id, "-p", port, batchmode,
			strict, tout, hostname, c)
	} else {
		if *tty {
			flag2 = "tt"
		}

		cmd = exec.Command(ssh, login, id, "-q", flag+flag2+"p",
			port, batchmode, strict, tout, hostname, command)
	}

	buf, err := cmd.CombinedOutput()
	if err != nil {
		return hostname + ": " + string(buf)
	}

	if *one {
		return hostname + ": " + string(buf)
	}
	return hostname + ":\n" + string(buf)
}

func main() {
	flag.Usage = func() {
		fmt.Println(
`usage: cs [-emqrstuVv1] [-c file] [-cmd] [-cname] [-d file] [-dd] [-du path]
	  [-f script.sh] [-h hosts_file] [-hba] [-i identity_file] [-io] [-ip]
	  [-l login_name] [-lip] [-mx] [-nc] [-nmap] [-ns] [-o output_file]
	  [-P port] [-p path] [-ping] [-soa] [-to timeout] [-top] [-tr]
	  [-tri] [-uname] [-vm] [command] [[user@]host] ...`)
		os.Exit(1)
	}

	lcmd := flag.Bool("cmd", false, "Local command")
	cname := flag.Bool("cname", false, "CNAME")
	copy := flag.String("c", "", "Copy")
	dd := flag.Bool("dd", false, "Dmidecode")
	disku := flag.String("du", "", "Du")
	download := flag.String("d", "", "Download")
	error := flag.Bool("e", false, "Error")
	file := flag.String("f", "", "Script file")
	hba := flag.Bool("hba", false, "HBA")
	hostsfile := flag.String("h", "", "Hosts file")
	id := flag.String("i", "", "Identity file")
	io := flag.Bool("io", false, "Iostat")
	ip := flag.Bool("ip", false, "IP")
	lip := flag.Bool("lip", false, "LIP")
	login := flag.String("l", "", "Login name")
	mount := flag.Bool("m", false, "Mounts")
	mx := flag.Bool("mx", false, "MX")
	netcat := flag.Bool("nc", false, "Netcat")
	nmap := flag.Bool("nmap", false, "Nmap")
	ns := flag.Bool("ns", false, "NS")
	one := flag.Bool("1", false, "One line")
	out := flag.String("o", "", "Output filename")
	path := flag.String("p", ".", "Path")
	png := flag.Bool("ping", false, "Ping")
	port := flag.String("P", "22", "SSH port")
	quiet := flag.Bool("q", false, "Quiet")
	recursive := flag.Bool("r", false, "Recursive")
	soa := flag.Bool("soa", false, "SOA")
	sorted := flag.Bool("s", false, "Sort")
	timeout := flag.String("to", "4", "Timeout")
	top1 := flag.Bool("top", false, "Top")
	tr := flag.Bool("tr", false, "Traceroute")
	tri := flag.Bool("tri", false, "Traceroute -I")
	tty := flag.Bool("t", false, "Force pseudo-tty allocation")
	uname := flag.Bool("uname", false, "Uname")
	up := flag.Bool("u", false, "Uptime")
	version := flag.Bool("V", false, "Version")
	verbose1 := flag.Bool("v", false, "Verbose mode 1")
	verbose2 := flag.Bool("vv", false, "Verbose mode 2")
	verbose3 := flag.Bool("vvv", false, "Verbose mode 3")
	vm := flag.Bool("vm", false, "Vmstat")
	flag.Parse()
	argv := flag.Args()

	if *version {
		fmt.Println("cs 1.1")
		os.Exit(1)
	}

	nocmd := 0
	if *cname || *dd || *disku != "" || *hba || *io || *ip || *lip ||
		*mount || *netcat || *nmap || *ns || *mx || *png || *soa ||
		*top1 || *tr || *tri || *uname || *up || *vm {
		nocmd = 1
	}

	command, hosts := checkArgs(nocmd, *copy, *download, *file, *hostsfile,
		tty, argv)

	ddir := 0
	if len(hosts) > 1 {
		ddir = 1
	}

	var f *os.File
	if *out != "" {
		now := time.Now()
		nowStr := now.Format(timeFmt)
		f = createFile(*out)
		f.WriteString("------ START: " + nowStr + " ------\n")
	}

	var mk, mk2 []string
	if *sorted {
		mk = make([]string, len(hosts))
		if *quiet {
			mk2 = make([]string, len(hosts))
		}
	}

	output := make(chan string)
	for _, hostname := range hosts {
		go func(hostname string) {
			output <- run(command, hostname, *id, *login, *path,
				*port, *timeout, *copy, *disku, *download, dd,
				cname, hba, io, ip, lcmd, lip, netcat, nmap, ns,
				mount, mx, one, png, recursive, soa, up,
				verbose1, verbose2, verbose3, top1, tr, tri,
				tty, uname, vm, ddir, f)
		}(hostname)
	}

	e, err := 0, 0
	for i := 0; i < len(hosts); i++ {
		select {
		case c := <-output:
			err = 0
			if !strings.Contains(c, ":\n") {
				e++
				err = 1
				if !strings.Contains(c, "\n") {
					c = c + "\n"
				}
			}

			if *sorted {
				if *error && err == 0 {
					continue
				}
				if *quiet {
					if *out != "" {
						mk[i] = c
					}
					split := strings.Split(c, ":")
					if err == 1 {
						mk2[i] = split[0] +"\t[ERROR]\n"
						continue
					}
					mk2[i] = split[0] + "\t[OK]\n"
					continue
				}
				mk[i] = c
				continue
			}
			if *out != "" {
				if *error && err == 0 {
					continue
				}
				f.WriteString(c)
			}
			if *quiet {
				split := strings.Split(c, ":")
				if *error && err == 0 {
					continue
				}

				if err == 1 {
					fmt.Println(split[0] + "\t[ERROR]")
					continue
				}
				fmt.Println(split[0] + "\t[OK]")
				continue
			}

			if *error && err == 0 {
				continue
			}
			fmt.Print(c)
		}
	}

	if *sorted {
		sort.Strings(mk)
		if *out != "" {
			for _, v := range mk {
				f.WriteString(v)
			}
		}
		if *quiet {
			sort.Strings(mk2)
			for _, v := range mk2 {
				fmt.Print(v)
			}
		} else {
			for _, v := range mk {
				fmt.Print(v)
			}
		}
	}

	if !*one {
		fmt.Printf("hosts = %d, errors = %d\n", len(hosts), e)
	}
	if *out != "" {
		now := time.Now()
		nowStr := now.Format(timeFmt)
		f.WriteString("------   END: " + nowStr + " ------\n")
		f.Close()
	}
}
