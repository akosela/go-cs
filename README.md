```
CS(1)                   FreeBSD General Commands Manual                  CS(1)

NAME
     cs -- concurrent ssh client

SYNOPSIS
     cs [-cdfqrs] [-h hosts_file] [-i identity_file] [-l login_name]
        [-o output_file] [-P port] [-p path] [-t timeout] {command | file}
        [[user@]host] ...

DESCRIPTION
     cs is a program for concurrently executing ssh(1)/scp(1) on a number of
     hosts.  It is intended to automate running remote commands or copying
     files between hosts on a network.  Public key authentication is used for
     establishing passwordless connection.

     The options are as follows:

     -c      Remote file copy mode.

     -d      Remote download mode.

     -f      Runs script file on the remote host.

     -h      Reads hosts from the given hosts_file.  Lines in the hosts_file
             can include commented lines beginning with a `#' and there can
             only be one host per line.

     -i      Selects a file from which the identity (private key) for public
             key authentication is read.  The default is ~/.ssh/id_rsa.

     -l      Specifies the user to log in as on the remote machine.  This also
             may be specified on a per-host basis on the command line.

     -o      Saves standard output and standard error to a file.

     -P      Port to connect to on the remote host.

     -p      Specifies remote or local path for files in a remote copy or
             download mode.

     -q      Quiet mode.  Supresses verbose standard output from remote
             servers.  This mode reports success or failure only.

     -r      Recursively copy entire directories.  It follows symbolic links
             encountered in the tree traversal.

     -s      Sorts output by lines.

     -t      Specifies the timeout (in seconds) used when connecting to the
             SSH server.  The default value is 5 seconds.

AUTHENTICATION
     The default method for authentication is a public key authentication
     which serves its purpose when dealing with multiple hosts.  You can read
     more about public key authentication in ssh(1).

EXIT STATUS
     The cs utility exits 0 on success, and >0 if an error occurs.

EXAMPLES
     Run commands on hosts foo and bar:

           $ cs 'uptime; uname -a' foo bar

     Run a command on multiple hosts specified in a file:

           $ cs -h hosts_file uptime

     Run a script on multiple hosts:

           $ cs -f script.sh foo{1..100}

     Copy file to multiple hosts using specified remote path:

           $ cs -h hosts_file -p foo/bar file

     Download recursively files from /foo/bar from multiple hosts to a speci-
     fied local path:

           $ cs -d -r -h hosts_file -p /tmp /foo/bar/\*

SEE ALSO
     scp(1), ssh(1), ssh_config(5)

AUTHORS
     Andy Kosela <akosela@andykosela.com>

FreeBSD 10.0                    August 14, 2014                   FreeBSD 10.0
```