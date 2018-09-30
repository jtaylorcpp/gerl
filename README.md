GERL 
========

Pronounced like gurl as in "hey gurl, hey"...


Gerl is an attempt to build out the remarkable parts of the Erlang/OTP for Go
while keeping in spirit with the language.

The vision is to provide a way to build, schedule, and manage both locally and globally
processes and their ability to communicate.

Gerl provides functionality for:

  - process id\'s

  - generic server (gen_servers)

  - message passing between pids

## Base Types

### Generic Servers

A Generic Server (gen_server in erlang) is a single-threaded state machine
in which the process handles synchronous and asynchronous events triggered by
messages being routed to the genserver. Each event also passes in and expects out
its state making it operate as a completely functional process.