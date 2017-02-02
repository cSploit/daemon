cSploit daemon (metasploit RPC feature)
======================================

[![Build Status](https://travis-ci.org/cSploit/daemon.svg?branch=feature/msf)](https://travis-ci.org/cSploit/daemon)

This is a feature of the cSploit project.
It manage the interactions with the Metasploit framework through an implementation of the MSFRPC API.

[Official MSFRPC documentation](https://rapid7.github.io/metasploit-framework/api/Msf/RPC.html)

Env setup
---------

To use this work a correct configuration of your golang installation is needed.
To do so, please refer to the README from the branch *develop*

Testing
-------

To run tests simply run:

```bash
go test -v {module_you_want_to_test}
e.g: go test -v AuthMSF_test.go if you want to test the calls related to the authentication
```
