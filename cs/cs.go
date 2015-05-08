// Copyright 2014 Andy Kosela.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Concurrent ssh client.
// cs is a program for concurrently executing ssh(1) or scp(1) on a
// number of hosts.  It is intended to automate running remote commands
// or copying files between hosts on a network.  Public key
// authentication is used for establishing passwordless connection.
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

const (
	scp     = "/usr/bin/scp"
	ssh     = "/usr/bin/ssh"
	timeFmt = "02-Jan-2006 15:04:05"
)

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

func checkArgs(copy, download, downloadFlat, file, hostsfile string,
	argv []string) (string, []string) {
	var command string
	var hosts []string

	if file != "" || copy != "" || download != "" || downloadFlat != "" {
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
			command = strings.Join(readFile(f), "")
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

func run(command, hostname, id, login, path, port, timeout, copy,
	download, downloadFlat string, one, recursive, verbose1, verbose2,
	verbose3 *bool, f *os.File) string {

	hostname = strings.Trim(hostname, "\n")
	strict := "StrictHostKeyChecking=no"
	tout := "ConnectTimeout=" + timeout
	flag := "-"
	if *verbose1 {
		flag = "-v"
	} else if *verbose2 {
		flag = "-vv"
	} else if *verbose3 {
		flag = "-vvv"
	}

	var cmd *exec.Cmd
	if copy != "" && *recursive {
		if login != "" {
			cmd = exec.Command(scp, flag+"r", "-i", id, "-P", port,
				"-o", strict, "-o", tout, copy, login+"@"+
				hostname+":"+path)
		} else {
			cmd = exec.Command(scp, flag+"r", "-i", id, "-P", port,
				"-o", strict, "-o", tout, copy, hostname+":"+
				path)
		}
	} else if copy != "" {
		if login != "" {
			cmd = exec.Command(scp, flag+"i", id, "-P", port, "-o",
				strict, "-o", tout, copy, login+"@"+hostname+
				":"+path)
		} else {
			cmd = exec.Command(scp, flag+"i", id, "-P", port, "-o",
				strict, "-o", tout, copy, hostname+":"+path)
		}
	} else if download != "" && *recursive || downloadFlat != "" &&
		*recursive {
		if downloadFlat == "" {
			path = exist(hostname, path)
		} else {
			download = downloadFlat
		}

		if login != "" {
			cmd = exec.Command(scp, flag+"r", "-i", id, "-P", port,
				"-o", strict, "-o", tout, login+"@"+hostname+
				":"+download, path)
		} else {
			cmd = exec.Command(scp, flag+"r", "-i", id, "-P", port,
				"-o", strict, "-o", tout, hostname+":"+download,
				path)
		}
	} else if download != "" || downloadFlat != "" {
		if downloadFlat == "" {
			path = exist(hostname, path)
		} else {
			download = downloadFlat
		}

		if login != "" {
			cmd = exec.Command(scp, flag+"i", id, "-P", port, "-o",
				strict, "-o", tout, login+"@"+hostname+":"+
				download, path)
		} else {
			cmd = exec.Command(scp, flag+"i", id, "-P", port, "-o",
				strict, "-o", tout, hostname+":"+download, path)
		}
	} else {
		if login != "" {
			cmd = exec.Command(ssh, flag+"i", id, "-l", login, "-p",
				port, "-o", strict, "-o", tout, hostname,
				command)
		} else {
			cmd = exec.Command(ssh, flag+"i", id, "-p", port, "-o",
				strict, "-o", tout, hostname, command)
		}
	}

	buf, err := cmd.CombinedOutput()
	if err != nil {
		return hostname + ": " + string(buf)
	}
	if *one {
		return hostname + ": " + string(buf)
	} else {
		return hostname + ":\n" + string(buf)
	}
}

func main() {
	flag.Usage = func() {
		fmt.Println(
`usage: cs [-qrsVv1] [-c file] [-d file] [-df file] [-f script.sh]
	  [-h hosts_file] [-i identity_file] [-l login_name] [-o output_file]
	  [-P port] [-p path] [-t timeout] [command] [[user@]host] ...`)
		os.Exit(1)
	}

	copy := flag.String("c", "", "Copy")
	download := flag.String("d", "", "Download")
	downloadFlat := flag.String("df", "", "Download flat")
	file := flag.String("f", "", "Script file")
	hostsfile := flag.String("h", "", "Hosts file")
	id := flag.String("i", string(os.Getenv("HOME")+"/.ssh/id_rsa"),
		"Identity file")
	login := flag.String("l", "", "Login name")
	one := flag.Bool("1", false, "One line")
	out := flag.String("o", "", "Output filename")
	port := flag.String("P", "22", "SSH port")
	path := flag.String("p", ".", "Path")
	quiet := flag.Bool("q", false, "Quiet")
	recursive := flag.Bool("r", false, "Recursive")
	sorted := flag.Bool("s", false, "Sort")
	timeout := flag.String("t", "5", "Timeout")
	version := flag.Bool("V", false, "Version")
	verbose1 := flag.Bool("v", false, "Verbose mode 1")
	verbose2 := flag.Bool("vv", false, "Verbose mode 2")
	verbose3 := flag.Bool("vvv", false, "Verbose mode 3")
	flag.Parse()
	argv := flag.Args()

	if *version {
		fmt.Println("cs v0.4")
		os.Exit(1)
	}

	command, hosts := checkArgs(*copy, *download, *downloadFlat, *file,
		*hostsfile, argv)

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

	output := make(chan string, 10)
	for _, hostname := range hosts {
		go func(hostname string) {
			output <- run(command, hostname, *id, *login, *path,
				*port, *timeout, *copy, *download,
				*downloadFlat, one, recursive, verbose1,
				verbose2, verbose3, f)
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

	if !*one {
		fmt.Printf("hosts = %d, errors = %d\n", len(hosts), e)
	}
	if *out != "" {
		now := time.Now()
		nowStr := now.Format(timeFmt)
		f.WriteString("------ END: " + nowStr + "   ------\n")
		f.Close()
	}
}
