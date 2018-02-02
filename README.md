cSploit daemon
==============

[![Build Status](https://travis-ci.org/cSploit/daemon.svg?branch=develop)](https://travis-ci.org/cSploit/daemon)

This is the core of the cSploit project.
It has been made to manage, provide, find and work with found resources.

As of now this software does not work, it's just a preview to 
perform an hand-off of the work as I found other devs that want to contribute.

Just run it
-----------
> ** Coming soon! ** ( docker run ... )


Env setup
---------

To work with Go lang you need to specify a path where Go will download
the required modules: `export GOPATH="$HOME/.gocode"` for instance.

Install `libpcap-dev libc-dev gcc git go` packages,
then get the sources `go get -t -u github.com/cSploit/daemon`.

After that sources are ready to be modified or built at `$GOPATH/src/github.com/cSploit/daemon`.

**Next commands assumes that your current cirectory is that one.**

Development
-------

To build the daemon run:

```bash
go build -i .
```

To run tests simply run:

```bash
go test -v ./... # run all tests
go test -v ./tools/... # run all tests in 'tools'
```

To start the daemon:

```bash
sudo ./daemon # root needed to sniff packets
```

And read nmap output from a file called `sample_nmap_out.xml`.
you can generate it by running `nmap -oX sample_nmap_out.xml -sV -T4 -O 192.168.0.0/24`.

Fork all the things!
-----------
You can manage your fork while contributing to the project, give [this](https://splice.com/blog/contributing-open-source-git-repositories-go/) a read :wink: .
In this way you can easily make pull rquests and experiments.


Tricks
-----------

In IntellijIDEA ( which I suggest you to use ) open the project from 
`$GOPATH/src/github.com/cSploit/daemon` and you're ready to *Go*!
