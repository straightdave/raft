Raft
====

Coding exercise trying to implement Raft consistency algorithm.

# Learning materials

## 1 - [https://raft.github.io/](https://raft.github.io/)
The Demo on the page is critical for me to understand the Raft.

## 2 - [Paper](https://raft.github.io/raft.pdf)
More detailed information.

# Exercise applications

```
├── cmd
│   ├── client
│   │   └── main.go
│   └── server
│       └── main.go
```

## Server
```
$ ./server
2020/04/29 09:51:24 Starts with 0 peer(s): []
2020/04/29 09:51:24 Serving TCP connections at :8765
2020/04/29 09:51:24 192.168.1.76:8765 becomes FOLLOWER
2020/04/29 09:51:24 192.168.1.76:8765 becomes CANDIDATE {term=1}
2020/04/29 09:51:24 192.168.1.76:8765 becomes LEADER {term=1}
```

## Client
```
$ ./client
127.0.0.1:8765 0> set a 1
OK
127.0.0.1:8765 1> get a
1
127.0.0.1:8765 2>
```
