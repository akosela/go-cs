
[![Build Status](https://travis-ci.org/akosela/go-cs.svg?branch=master)](https://travis-ci.org/akosela/go-cs)
[![Go Report Card](https://goreportcard.com/badge/github.com/akosela/go-cs)](https://goreportcard.com/report/github.com/akosela/go-cs)
[![GoDoc](https://godoc.org/github.com/akosela/go-cs/cs?status.svg)](https://godoc.org/github.com/akosela/go-cs/cs)

## go-cs 

cs is a program for concurrently executing local or remote commands on
multiple hosts.  It is using OpenSSH for running remote commands.  Its
main purpose is to help automate and manage large network of hosts.  So
in essence this tool is very similar to Ansible, Chef or Puppet, but
much more simpler and faster in execution.  Speed and simplicity are its
main goals.

## Install

#### Go get

```
$ go get github.com/akosela/go-cs/cs
# ln -s $GOPATH/bin/cs /usr/bin/cs
```

#### Linux (rpm)

```
$ curl -OL https://github.com/akosela/go-cs/releases/download/v0.8/go-cs-0.8-1.x86_64.rpm
# rpm -ivh go-cs-0.8-1.x86_64.rpm
```

#### Linux

```
$ curl -OL https://github.com/akosela/go-cs/releases/download/v0.8/go-cs-0.8.linux.tar.gz
$ tar xvf go-cs-0.8.linux.tar.gz
$ cd go-cs
$ gzip cs.1
# cp cs /usr/bin ; cp cs.1.gz /usr/share/man/man1
```

#### FreeBSD

Package

```
# pkg install go-cs
```

Port

```
# cd /usr/ports/net/go-cs
# make install clean
```

#### OpenBSD

Package

```
# pkg_add go-cs
```

Port

```
# cd /usr/ports/sysutils/go-cs
# make install clean
```

#### MacOS

```
$ curl -OL https://github.com/akosela/go-cs/releases/download/v0.8/go-cs-0.8.darwin.tar.gz
$ tar xvf go-cs-0.8.darwin.tar.gz
$ cd go-cs
$ gzip cs.1
# cp cs /opt/local/bin ; cp cs.1.gz /opt/local/share/man/man1
```

## Man page

```
CS(1)                      BSD General Commands Manual                   CS(1)

NAME
     cs - concurrent ssh client

SYNOPSIS
     cs [-eqrstuVv1] [-c file] [-cmd] [-cname] [-d file] [-dd] [-du path]
        [-f script.sh] [-h hosts_file] [-i identity_file] [-io] [-ip]
        [-l login_name] [-mx] [-nc] [-nmap] [-ns] [-o output_file] [-P port]
        [-p path] [-ping] [-soa] [-to timeout] [-top] [-tr] [-tri] [-uname]
        [-vm] [command] [[user@]host] ...

DESCRIPTION
     cs is a program for concurrently executing local or remote commands on
     multiple hosts.  It is using OpenSSH for running remote commands.  Its
     main purpose is to help automate and manage large network of hosts.

     The options are as follows:

     -c file
             Copy file to the remote machine.

     -cmd    Run an arbitrary local command concurrently on multiple hosts.

     -cname  Run a local DNS query of type CNAME.

     -d file
             Download file from the remote machine.  It will be saved in a
             directory named after the remote host only when you download from
             multiple servers.

     -dd     Prints basic hardware information for a specific host (sudo(8)
             and dmidecode(8) required).

     -du path
     	     Prints the list of largest files in a specific directory (sudo(8) 
	     required, units in M).

     -e      Prints hosts with errors only.

     -f script.sh
             Run a local shell script on the remote host.

     -h hosts_file
             Reads hostnames from the given hosts_file.  Lines in the
             hosts_file can include commented lines beginning with a `#' and
             only one host per line is allowed.

     -i identity_file
             Selects a file from which the identity (private key) for public
             key authentication is read.  The default is ~/.ssh/id_rsa.

     -io     Run remote iostat(1) on a specific host.

     -ip     Lists network interfaces.

     -l login_name
             Specifies the user to log in as on the remote machine.  This also
             may be specified on a per-host basis on the command line.

     -mx     Run a local DNS query of type MX.

     -nc     Tests specified port with nc(1).  Default is 22/tcp.

     -nmap   Scans host with nmap(1).

     -ns     Run a local DNS query of type NS.

     -o output_file
             Saves standard output and standard error to a file.

     -P port
             Port to connect to on the remote host.

     -p path
             Specifies remote or local path for files in a remote copy or
             download mode.

     -ping   Sends ICMP ECHO_REQUEST to network hosts.

     -q      Quiet mode.  Supresses verbose standard output from remote
             machines.  This mode reports success or failure only.

     -r      Recursively copy entire directories.  It follows symbolic links
             encountered in the tree traversal.

     -s      Sort output.

     -soa    Run a local DNS query of type SOA.

     -t      Force pseudo-tty allocation.	

     -to timeout
             Specifies the timeout (in seconds) used when connecting to the
             SSH server.  The default value is 4 seconds.

     -top    Run remote top(1) in batch mode on a specific host.

     -tr     Run local traceroute(8).

     -tri    Run local traceroute(8) using ICMP (local sudo(8) required).

     -uname  Prints remote system information including OS version.

     -u      Run remote uptime(1) on a specific host.

     -V      Displays the version number and exit.

     -v      Verbose mode.  Causes cs to print debugging messages from ssh(1)
             about its progress.  This is helpful in debugging connection,
             authentication, and configuration problems.  Multiple -v options
             increase the verbosity.  The maximum is 3.

     -vm     Run remote vmstat(8) on a specific host.

     -1      One line mode, useful for sorting output later.

AUTHENTICATION
     The default method for authentication is a public key authentication
     which serves its purpose when dealing with multiple hosts.  You can read
     more about public key authentication in ssh(1).

EXIT STATUS
     The cs utility exits 0 on success, and >0 if an error occurs.

EXAMPLES
     Run a series of commands on hosts foo and bar:

           $ cs 'uptime; uname -a' foo bar

     Run a command on multiple hosts specified in a hosts_file:

           $ cs -h hosts_file uptime

     Run a local shell script on multiple hosts:

           $ cs -f script.sh foo{1..100}

     Run a local shell script with sudo(8) on multiple hosts:

           $ cs -t -f script.sh foo{1..100}

     Copy file to multiple hosts using a specific remote path:

           $ cs -c file -h hosts_file -p /foo/bar

     Download file from host foo:~ to the current working directory:

           $ cs -d file foo

     Download recursively files from /foo/bar from multiple hosts to a speci-
     fic local path /tmp with subdirectories named after remote hosts:

           $ cs -r -d '/foo/bar/*' -h hosts_file -p /tmp

     Run a command on multiple hosts and sort the output:

           $ cs -1 -h hosts_file 'free -m | grep Swap' | sort -rnk4 | head

     Run local ping(1) on multiple hosts:

           $ cs -ping foo{1..100}

     Show name and size of all files matching the glob pattern:

           $ cs -du '/var/log/mes*' foo{1..100}

SEE ALSO
     scp(1), ssh(1), ssh_config(5), sudo(8)

AUTHORS
     Andy Kosela <akosela@andykosela.com>

BSD                              July 10, 2017                             BSD  
```
