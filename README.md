GERL 
========

Pronounced like gurl ... as in "hey gurl, hey"...


Gerl is an attempt to build out the remarkable parts of the Erlang/OTP for Go
while keeping in spirit with the language.

The vision is to provide a way to build, schedule, and manage both locally and globally
avaialble processes and their ability to communicate.

Gerl provides functionality for:

  - process id\'s

  - generic server (gen_servers)

  - message passing between pids

This is mainly done by using:

  - channels

  - go routines

  - grpc

## Basic Concepts

### Process ID (Pid)

Processes IDs (pid) is the main abstraction for communicating
with a running process. The pid contains channels for bidirectional communication 
from the running process and, under the hood, handles the GRPC implemtation. 

### Generic Server (genserver)

Generic servers, genservers, are a concept directly pulled from Erlang/OTP. A genserver
is a process that has a pre-specified set of functionality; mainly a *call* and *cast*.


*call* is a bidirectional action in which a client sends a message to and expects
a result back from a genserver. The genserver has a specific function dedicated
to handling *call* actions.

*cast* is a unidirectional action in which the client sends a message to a genserver
and moves on.

The genserver client builds the GRCP client necessary to make the calls and needs the 
address of the pid of the genserver to send messages back and forth.