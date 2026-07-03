Absolutely. Given your profile, I actually **wouldn't** build "just a reverse proxy." That's a good project, but it doesn't tell the strongest story.

Your portfolio is already becoming:

* Custom `malloc`
* HTTP server in C++
* `poll()` → `epoll()`
* Linux internals

The next project should look like something a backend engineer at Cloudflare, Caddy, Traefik, NGINX, or an infrastructure startup might build.

I'd frame it as:

> **EdgeProxy** — A production-inspired reverse proxy and load balancer written in Go.

The goal isn't to recreate NGINX. The goal is to learn how modern backend infrastructure is built while producing something that looks substantial on a resume.

---

# agenda.md

```md
# EdgeProxy
A production-inspired reverse proxy and load balancer in Go.

Goal:
Understand HTTP internals, concurrency, networking,
load balancing, observability and backend architecture.

Rule:
Never read ahead.
Each milestone introduces exactly the concepts needed.
Every milestone should produce a working system.
```

---

# Milestone 0 — Learn enough Go

Goal:

Don't become a Go expert.

Become productive.

Build

```
hello.go

variables

functions

structs

goroutines

maps

slices

http server
```

Resources

* Tour of Go
* Go by Example
* Effective Go (skim)

Done when

```
GET /

returns

Hello World
```

---

# Milestone 1 — Basic HTTP server

Instead of reverse proxying...

Build a tiny web server.

Learn

```
net/http

http.Handler

ResponseWriter

Request
```

Add

```
logging

multiple routes

JSON response
```

Now you understand Go's HTTP model.

---

# Milestone 2 — HTTP client

Reverse proxies spend half their lives being clients.

Build

```
CLI

edge fetch https://example.com
```

Print

```
status

headers

body

latency
```

Learn

```
http.Client

Transport

Timeouts

Connection reuse
```

---

# Milestone 3 — First reverse proxy

Browser

↓

Proxy

↓

Backend

↓

Browser

No modifications.

Just forwarding.

Learn

```
httputil.ReverseProxy

Director

Request lifecycle
```

Read

Go source code.

Seriously.

The standard library reverse proxy is only a few hundred lines.

---

# Milestone 4 — Header manipulation

Add

```
X-Forwarded-For

X-Forwarded-Host

X-Forwarded-Proto

Request-ID
```

Inject

```
Server

Edge-Version

Trace-ID
```

Learn

Why proxies modify headers.

Read

RFC 7239

Forwarded headers.

---

# Milestone 5 — Structured logging

Replace

```
fmt.Println
```

with

```
JSON logs
```

Every request

```
method

path

status

latency

backend

client ip

trace id
```

Read

12 Factor Apps

Logging

---

# Milestone 6 — Configuration

Move

```
backend

port

timeouts
```

into

```
config.yaml
```

Support

```
multiple backends
```

Now your project starts looking real.

---

# Milestone 7 — Round robin load balancing

Instead of

```
backend A
```

Support

```
A

B

C
```

Round robin.

Learn

```
mutex

atomic

goroutines
```

---

# Milestone 8 — Passive health checking

Backend returns

```
500
```

Too many times?

Temporarily remove it.

Learn

```
failure counting

atomic operations

shared state
```

---

# Milestone 9 — Active health checking

Every 5 seconds

```
GET /health
```

Recover dead servers automatically.

Learn

```
goroutines

tickers

background workers
```

---

# Milestone 10 — Graceful shutdown

Ctrl+C

Should

* stop accepting requests

* finish in-flight requests

* close cleanly

Learn

```
context

signal.Notify

Shutdown()
```

---

# Milestone 11 — Timeouts

Support

```
read timeout

write timeout

idle timeout

backend timeout
```

Now you'll understand

Why proxies exist.

---

# Milestone 12 — Rate limiter

Implement

Token Bucket

Don't use a package.

Build it.

Learn

```
mutex

maps

time

cleanup goroutines
```

---

# Milestone 13 — Reverse proxy middleware

Implement middleware

```
Logging

Recovery

Rate limiting

Authentication

Compression
```

Understand Go middleware patterns.

---

# Milestone 14 — Retry logic

Backend fails?

Retry.

Implement

```
max retries

backoff

retry only idempotent requests
```

Read

HTTP idempotency.

---

# Milestone 15 — Circuit breaker

One backend keeps failing.

Stop sending traffic.

Half-open.

Closed.

Open.

Classic systems interview topic.

---

# Milestone 16 — Metrics

Expose

```
/metrics
```

Track

```
requests

latency

active connections

backend health

errors

retries
```

Read

Prometheus exposition format.

---

# Milestone 17 — Request tracing

Generate

```
Trace-ID
```

Pass

Across

Every backend.

Makes debugging distributed systems possible.

---

# Milestone 18 — Response cache

Cache

```
GET
```

responses.

Support

TTL

Eviction

Cache-Control

Great systems exercise.

---

# Milestone 19 — Compression

If client supports

```
gzip
```

Compress.

Learn

```
Accept-Encoding

Content-Encoding
```

---

# Milestone 20 — Static file serving

Proxy

OR

Serve

```
index.html

css

images
```

Learn

Conditional requests

ETags

Last Modified.

---

# Milestone 21 — WebSocket proxying

Huge learning milestone.

Support

```
Upgrade

Connection

Hijacking
```

Learn

Why WebSockets break many proxies.

---

# Milestone 22 — HTTP/2

Read

How Go handles it.

You probably won't implement it.

But you'll understand

multiplexing

streams

flow control.

---

# Milestone 23 — Docker

Containerize

```
Proxy

Backend

Prometheus

Grafana
```

Now benchmark everything.

---

# Milestone 24 — Benchmark

Use

```
wrk

hey

ab
```

Measure

```
latency

throughput

memory

CPU
```

Find bottlenecks.

Optimize.

---

# Milestone 25 — README

Write

Architecture

```
Browser

↓

Proxy

↓

Balancer

↓

Middleware

↓

Transport

↓

Backend
```

Explain

* request lifecycle
* concurrency model
* retries
* health checks
* caching
* middleware
* benchmarks
* future improvements

---

# Stretch Goals

## Sticky sessions

Hash

```
Client IP
```

---

## Consistent hashing

For caches.

---

## Least connections

Instead of

Round Robin.

---

## Hot reload config

Reload

```
config.yaml
```

Without restarting.

---

## TLS termination

HTTPS

↓

HTTP backend

---

## Dynamic backend discovery

Read backends

From

```
JSON

API

etcd

Consul
```

---

## Admin API

```
GET /admin/backends

POST /admin/backend

DELETE /admin/backend
```

---

## Dashboard

Live

```
traffic

latency

health

errors
```

---

# Resources

## Go

* Tour of Go
* Effective Go
* Go source (`net/http`, `httputil`)

## HTTP

* RFC 9110 (Semantics)
* RFC 9112 (HTTP/1.1)

## Books

* Designing Data-Intensive Applications (selected chapters)
* The Go Programming Language
* Network Programming with Go

## Blogs

* Cloudflare Engineering
* Caddy Server blog
* Traefik blog
* NGINX architecture docs
* Uber Engineering
* Grab Engineering

---

# Resume Outcome

By the end, this project should be something you can confidently describe as:

> **EdgeProxy** — A production-inspired reverse proxy and load balancer written in Go featuring configurable routing, round-robin load balancing, active/passive health checks, retry and circuit-breaker mechanisms, middleware pipeline, rate limiting, response caching, graceful shutdown, Prometheus metrics, request tracing, WebSocket support, Dockerized deployment, and performance benchmarking.

## One recommendation

Unlike your C web server—which is valuable because you built the networking stack yourself—**lean into Go's ecosystem** here. Don't avoid `net/http` or `httputil.ReverseProxy` just because they abstract things away. Read their source code, understand how they work, and then extend or customize them where appropriate. Recruiters looking at Go backend roles are generally more interested in seeing idiomatic Go, concurrency patterns, middleware design, observability, resilience, and clean architecture than a manual reimplementation of HTTP parsing.

Taken together with your custom allocator and C server, this creates a compelling narrative: you understand systems from the lowest levels up, and you can also build production-style backend infrastructure using modern tooling.

