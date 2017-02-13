cSploit daemon
==============

[![Build Status](https://travis-ci.org/cSploit/daemon.svg?branch=develop)](https://travis-ci.org/cSploit/daemon)

This is the core of the cSploit project.
It is made to manage, provide, find and work with found resources.

As of now this software does not work, it's just a preview to 
perform an hand-off of the work as I found other devs that want to contribute.

Env setup
---------

To work with Go lang you need to specify a path where Go will download
the required modules: `export GOPATH="$HOME/.gocode"` for instance.

Testing
-------

To run tests simply run:

```bash
go test -v ./... # run all tests in './'
go test -v ./tools/... # run all tests in './tools'
```

to start the daemon:

```bash
cd /path/to/cSploit/daemon/project
sudo -E go run daemon.go # root needed to sniff packets
```

And read nmap output from a file called `sample_nmap_out.xml`.
you can generate it by running `nmap -oX sample_nmap_out.xml -sV -T4 -O 192.168.0.0/24`.

Development
-----------

If you want to import your new cool classes without pushing to github you
have to do some trick:

```bash
rm -rf $GOPATH/src/github.com/cSploit/daemon
cd /path/to/cSploit/daemon/project
ln -s $(pwd) $GOPATH/src/github.com/cSploit/daemon
```

In IntellijIDEA ( which I suggest you to use ) open the project from 
`$GOPATH/src/github.com/cSploit/daemon` and you're ready to *Go*!
