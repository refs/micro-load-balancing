# ToC

- [ToC](#toc)
  - [Premise](#premise)
  - [Preconditions:](#preconditions)
  - [Simple](#simple)
  - [Advanced](#advanced)
  - [Real Life Example](#real-life-example)
    - [Demonstration](#demonstration)

## Premise

We want to create 2 micro services, both http services and both listen to `/hello`:

- service 1
  - `/hello` => `1 here`

- service 2
  - `/hello` => `2 here`

Is it possible with go-micro to do such thing?

Answer: it appears to be possible. [As per the docs](https://m3o.com/docs/framework.html):

> Load Balancing - Client side load balancing built on service discovery. Once we have the addresses of any number of instances of a service we now need a way to decide which node to route to. We use random hashed load balancing to provide even distribution across the services and retry a different node if there’s a problem.

## Preconditions:

[micro binary installed](https://github.com/micro/micro/#install):

```
# MacOS
curl -fsSL https://raw.githubusercontent.com/micro/micro/master/scripts/install.sh | /bin/bash

# Linux
wget -q  https://raw.githubusercontent.com/micro/micro/master/scripts/install.sh -O - | /bin/bash

# Windows
powershell -Command "iwr -useb https://raw.githubusercontent.com/micro/micro/master/scripts/install.ps1 | iex"
```

## Simple
- go run main.go --msg=1

on another session:
- go run main.go --msg=2

on yer another session:
`micro list services` should output something like:

```
❯ /usr/local/bin/micro list services
go.micro.api
go.micro.api.hello
```

running `micro get service go.micro.api.hello`:

```
❯ micro get service go.micro.api.hello
service  go.micro.api.hello

version latest

ID      Address Metadata
go.micro.api.hello-3c1fed9d-95ca-4555-927e-7a6ce6c281c3 10.42.17.12:60062       registry=mdns,server=grpc,transport=grpc,broker=http,protocol=grpc
go.micro.api.hello-be14eebb-1d83-491d-80c2-dc582185ed29 10.42.17.12:60059       server=grpc,transport=grpc,broker=http,protocol=grpc,registry=mdns

Endpoint: Say.Hello

Request: {
        message_state MessageState {
                no_unkeyed_literals NoUnkeyedLiterals
                do_not_compare DoNotCompare
                do_not_copy DoNotCopy
                message_info MessageInfo
        }
        int32 int32
        unknown_fields []uint8
}

Response: {
        message_state MessageState {
                no_unkeyed_literals NoUnkeyedLiterals
                do_not_compare DoNotCompare
                do_not_copy DoNotCopy
                message_info MessageInfo
        }
        int32 int32
        unknown_fields []uint8
        msg string
}
```

## Advanced

And as you can see on the first couple lines of this output, the running nodes with its address can be seen. There is a programmatically way to get this information, bypassing the default load balancing by micro and jump straight to the necessary service.

A combination of metadata + registry can be used to query the mdns registry (for the purposes of this demo). Prepare your environment with the following commands (each command should run on a different terminal session):

```bash
## NOTE: each one of this on its own terminal session.
/usr/local/bin/micro api
go run main.go --msg=1 --ins=1
go run main.go --msg=2 --ins=2
go run main.go --msg=3

go run main.go --q
```

the output of the last should be something simmilar to:

```bash
addr: 10.42.17.12:61303
metadata: map[broker:http ocis_instance:2 protocol:grpc registry:mdns server:grpc transport:grpc]
addr: 10.42.17.12:61318
metadata: map[broker:http ocis_instance:3 protocol:grpc registry:mdns server:grpc transport:grpc]
addr: 10.42.17.12:61298
metadata: map[broker:http ocis_instance:1 protocol:grpc registry:mdns server:grpc transport:grpc]
```

## Real Life Example

Q: Suppose in oCIS, a proxy wants to route to storage #2, but there is no information for where is this service running, and the random load balancer from Micro doesn't really help us as the request can only be served by the service the url points to. How can this be solved?
A: With a solution like this, we can have a metadata value on every service and query the registry for services with such info. Let's see how this would be on this repo:

### Demonstration

```bash
## NOTE: each one of this on its own terminal session.
/usr/local/bin/micro api
go run main.go --msg=1 --ins=owncloud
go run main.go --msg=2 --ins=eos
go run main.go --msg=3 --ins=local

go run main.go --q --ins=owncloud
```

the output of the last should be something simmilar to:

```bash
addr: 10.42.17.12:61473
```

On a higher level, this functionality can be used like:

1. request to `/storage/owncloud/`
2. get the storage name from the url: `owncloud`
3. query the registry like: `go run main.go --q --ins=owncloud`
4. get node address
5. create gRPC client
6. do gRPC request