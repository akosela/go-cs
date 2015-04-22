## go-cs

cs is a program for concurrently executing ssh(1)/scp(1) on a number of
hosts.  It is intended to automate running remote commands or copying
files between hosts on a network.  Public key authentication is used for
establishing passwordless connection.

## Install

#### Go get

```
$ go get github.com/akosela/go-cs/cs
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

#### Linux (rpm)

```
$ curl -OL https://github.com/akosela/go-cs/releases/download/v0.3/go-cs-0.3-1.x86_64.rpm
# rpm -ivh go-cs-0.3-1.x86_64.rpm
```

#### Linux

```
$ curl -OL https://github.com/akosela/go-cs/releases/download/v0.3/go-cs-0.3.linux.amd64.tar.gz
$ tar xvf go-cs-0.3.linux.amd64.tar.gz
$ cd go-cs
$ gzip cs.1
# cp cs /usr/bin ; cp cs.1.gz /usr/share/man/man1
```

#### OS X

```
$ curl -OL https://github.com/akosela/go-cs/releases/download/v0.3/go-cs-0.3.darwin.amd64.tar.gz
$ tar xvf go-cs-0.3.darwin.amd64.tar.gz
$ cd go-cs
$ gzip cs.1
# cp cs /opt/local/bin ; cp cs.1.gz /opt/local/share/man/man1
```

## Man page

```
CS(1)                   FreeBSD General Commands Manual                  CS(1)

NAME
     cs -- concurrent ssh client

SYNOPSIS
     cs [-qrsVv1] [-c file] [-d file] [-f script.sh] [-h hosts_file]
        [-i identity_file] [-l login_name] [-o output_file] [-P port]
        [-p path] [-t timeout] [command] [[user@]host] ...

DESCRIPTION
     cs is a program for concurrently executing ssh(1) or scp(1) on a number
     of hosts.  It is intended to automate running remote commands or copying
     files between hosts on a network.  Public key authentication is used for
     establishing passwordless connection.

     The options are as follows:

     -c file
             Copy file to the remote machine.

     -d file
             Download file from the remote machine.

     -f script.sh
             Runs shell script on the remote host.

     -h hosts_file
             Reads hostnames from the given hosts_file.  Lines in the
             hosts_file can include commented lines beginning with a `#' and
             only one host per line is allowed.

     -i identity_file
             Selects a file from which the identity (private key) for public
             key authentication is read.  The default is ~/.ssh/id_rsa.

     -l login_name
             Specifies the user to log in as on the remote machine.  This also
             may be specified on a per-host basis on the command line.

     -o output_file
             Saves standard output and standard error to a file.

     -P port
             Port to connect to on the remote host.

     -p path
             Specifies remote or local path for files in a remote copy or
             download mode.

     -q      Quiet mode.  Supresses verbose standard output from remote
             machines.  This mode reports success or failure only.

     -r      Recursively copy entire directories.  It follows symbolic links
             encountered in the tree traversal.

     -s      Sorts output by lines.

     -t timeout
             Specifies the timeout (in seconds) used when connecting to the
             SSH server.  The default value is 5 seconds.

     -V      Displays the version number and exit.

     -v      Verbose mode.  Causes cs to print debugging messages from ssh(1)
             about its progress.  This is helpful in debugging connection,
             authentication, and configuration problems.  Multiple -v options
             increase the verbosity.  The maximum is 3.

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

     Run a shell script on multiple hosts:

           $ cs -f script.sh foo{1..100}

     Copy file to multiple hosts using specified remote path:

           $ cs -c file -h hosts_file -p /foo/bar

     Download recursively files from /foo/bar from multiple hosts to a speci-
     fied local path:

           $ cs -r -d /foo/bar/\* -h hosts_file -p /tmp

     Run a command on multiple hosts and sort the output:

           $ cs -1 -h hosts_file 'free -m | grep Swap' | sort -rnk4 | head

SEE ALSO
     scp(1), ssh(1), ssh_config(5)

AUTHORS
     Andy Kosela <akosela@andykosela.com>

FreeBSD 10.0                    April 22, 2015                    FreeBSD 10.0
```
