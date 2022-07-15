# Quick start

----

### TCP Server parameters:

#### Port:
    required:true yaml:"port" env_var_name:"SERVER_TCP_PORT" default:"6677"
#### Accept delay in milliseconds:
    name in config yml:"accept_delay_millisecond" env_var_name:"SERVER_ACCEPT_DELAY_MILLISECOND" default:"10"
#### Connection timeout in seconds
    name in config yml:"conn_timeout_second" env_var_name:"SERVER_CONN_TIMEOUT" default:"5"
#### Read timeout in seconds:
    name in config yml:"read_timeout_second" env_var_name:"SERVER_READ_TIMEOUT" default:"3"

---

### TCP Client parameters:

#### Port 
    required:true yaml:"port" env_var_name:"CLIENT_TCP_PORT" default:"6677"
#### Dail timeout in seconds
    name in config yml:"dial_timeout_second" env_var_name:"CLIENT_DIAL_TIMEOUT_SECOND" default:"3"
#### Connection timeout in seconds
    name in config yml:"conn_timeout_second" env_var_name:"CLIENT_CONN_TIMEOUT_SECOND" default:"5"
#### Max dail count in byte
    name in config yml:"max_dail_count_byte" env_var_name:"CLIENT_MAX_DAIL_COUNT_BYTE" default:"10"
#### Delay between dail in milliseconds:
    name in config yml:"delay_dail_millisecond" env_var_name:"CLIENT_DELAY_DAIL_MILLISECOND" default:"10"
#### Read timeout in seconds:
    name in config yml:"read_timeout_second" env_var_name:"CLIENT_READ_TIMEOUT_SECOND" default:"3"

---

### Protocol

Every message must be written and read as a tuple (*n, **content).

*n - 4 bytes, represents a length of content

**content - slice of bytes with length n


---

### POW algorithm

When a client connects to a server, the server responds
to the client with a challenge. The challenge is a
tuple (n*, bt16**). Upon receiving the request,
the client brute force finds bt16+attemptsCount***,
which contains more or equal number of leading zeros (n)
in the sha256 hash. If such a hash is found before the
timeout, then return the number of attempts as
bytes to the server.

*n - number of leading zeros, byte

**bt16 - randomly generated 16 bytes

***attemptsCount - uint32 (starting from 0 to max value uint32)

---
### Makefile
To run server:

```sh
$ make run_server
```

To run client:

```sh
$ make run_client
```

To build server image:

```sh
$ make build_server_image
```

To build client image:

```sh
$ make build_client_image
```

To run server docker:

```sh
$ make run_server_docker
```

To run client docker:

```sh
$ make run_client_docker
```

