// Copyright 2014 Andy Kosela.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Concurrent ssh client.
// cs is a program for concurrently executing ssh(1)/scp(1) on a number
// of hosts.  It is intended to automate running remote commands or
// copying files between hosts on a network.  Public key authentication
// is used for establishing passwordless connection.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

const timeFmt = "02-Jan-2006 15:04:05"

func createFile(path string) *os.File {
	file, err := os.OpenFile(path,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
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

func run(command, hostname, id, login, out, path, port, timeout string,
	copy, download, recursive, quiet *bool, f *os.File) string {

	hostname = strings.Trim(hostname, "\n")
	tout := "ConnectTimeout=" + timeout
	var cmd *exec.Cmd
	if login != "" {
		cmd = exec.Command("/usr/bin/ssh", "-i", id, "-l",
			login, "-p", port, "-o", tout, hostname, command)
	} else if *copy && *recursive {
		cmd = exec.Command("/usr/bin/scp", "-r", "-i", id, "-P",
			port, "-o", tout, command, hostname+":"+path)
	} else if *copy {
		cmd = exec.Command("/usr/bin/scp", "-i", id, "-P", port,
			"-o", tout, command, hostname+":"+path)
	} else if *download && *recursive {
		cmd = exec.Command("/usr/bin/scp", "-r", "-i", id, "-P",
			port, "-o", tout, hostname+":"+command, path)
	} else if *download {
		cmd = exec.Command("/usr/bin/scp", "-i", id, "-P", port,
			"-o", tout, hostname+":"+command, path)
		fmt.Println(cmd)
	} else {
		cmd = exec.Command("/usr/bin/ssh", "-i", id, "-p", port,
			"-o", tout, hostname, command)
	}

	buf, err := cmd.CombinedOutput()
	if err != nil {
		return hostname + ": " + string(buf)
	}
	return hostname + ":\n" + string(buf)
}

func main() {
	flag.Usage = func() {
		fmt.Println(
`usage: cs [-cdfqrs] [-h hosts_file] [-i identity_file] [-l login_name]
	  [-o output_file] [-P port] [-p path] [-t timeout] {command | file}
	  [[user@]host] ...`)
		os.Exit(1)
	}

	copy := flag.Bool("c", false, "Copy")
	download := flag.Bool("d", false, "Download")
	file := flag.Bool("f", false, "Script file")
	hostsfile := flag.String("h", "", "Hosts file")
	id := flag.String("i", string(os.Getenv("HOME")+"/.ssh/id_rsa"),
		"Identity file")
	login := flag.String("l", "", "Login name")
	out := flag.String("o", "", "Output filename")
	port := flag.String("P", "22", "SSH port")
	path := flag.String("p", ".", "Path")
	quiet := flag.Bool("q", false, "Quiet")
	recursive := flag.Bool("r", false, "Recursive")
	sorted := flag.Bool("s", false, "Sort")
	timeout := flag.String("t", "5", "Timeout")
	flag.Parse()
	argv := flag.Args()

	if len(argv) < 1 {
		flag.Usage()
	}

	var f *os.File
	if *out != "" {
		now := time.Now()
		nowStr := now.Format(timeFmt)
		f = createFile(*out)
		f.WriteString("------ START: " + nowStr + " ------\n")
	}

	hosts := argv[1:]
	if *hostsfile != "" {
		f := openFile(*hostsfile)
		hosts = readFile(f)
	}

	if *file {
		f := openFile(argv[0])
		argv[0] = strings.Join(readFile(f), "")
	}

	var mk, mk2 []string
	if *sorted {
		mk = make([]string, len(hosts))
		if *quiet {
			mk2 = make([]string, len(hosts))
		}
	}

	output := make(chan string, 10)
	for _, hostname := range hosts {
		go func(hostname string) {
			output <- run(argv[0], hostname, *id, *login,
				*out, *path, *port, *timeout, copy,
				download, recursive, quiet, f)
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
			}
			if *sorted {
				if *quiet {
					if *out != "" {
						mk[i] = c
					}
					split := strings.Split(c, ":")
					if err == 1 {
						mk2[i] = split[0] +
							"\t[ERROR]\n"
						continue
					}
					mk2[i] = split[0] + "\t[OK]\n"
					continue
				}
				mk[i] = c
				continue
			}
			if *out != "" {
				f.WriteString(c)
			}
			if *quiet {
				split := strings.Split(c, ":")
				if err == 1 {
					fmt.Println(split[0] + "\t[ERROR]")
					continue
				}
				fmt.Println(split[0] + "\t[OK]")
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

	fmt.Printf("hosts = %d, errors = %d\n", len(hosts), e)
	if *out != "" {
		now := time.Now()
		nowStr := now.Format(timeFmt)
		f.WriteString("------ END: " + nowStr + "   ------\n")
		f.Close()
	}
}
