# IF3230 - Sistem Paralel dan Terdistribusi 2024
> Tugas Besar Consensus Protocol: Raft

## About
Simple Key Value Store utilizing Raft Consensus Protocol

## Contributors
Kelompok SPG

| NIM | Nama |
| --- | --- |
| 13521004 | Henry Anand Septian Radityo |
| 13521007 | Matthew Mahendra |
| 13521010 | Salman Hakim Alfarisi |
| 13521015 | Hidayatullah Wildan Ghaly |
| 13521024 | Ahmad Nadil |

<img src="docs/spg.jpg" alt="Kelompok 1 K03" width="1280" height="720">

## Features
### Main Features
- Reliable key value store using Raft Consensure protocol
- Handle connection in bad network condition (Tested using Clumsy)

### Bonus Features
- Implementation in Go
- Local Network / Internet Demo
- Unit Test
- Transaction
- Web Client

## How to Run
1. [Backend README](backend/README.md)
2. [Frontend README](frontend/README.md)


Alternatively, you can run the main program to switch between server and client mode
```bash
python main.py
```

## Folder Structure
```bash
📦backend
 ┣ 📂cmd
 ┃ ┣ 📂client
 ┃ ┃ ┗ 📜main.go
 ┃ ┣ 📂server
 ┃ ┃ ┗ 📜main.go
 ┣ 📂internal
 ┃ ┣ 📂client
 ┃ ┃ ┗ 📜client.go
 ┃ ┣ 📂handlers
 ┃ ┃ ┗ 📜handlers.go
 ┃ ┣ 📂raft
 ┃ ┃ ┣ 📜election.go
 ┃ ┃ ┣ 📜heartbeat.go
 ┃ ┃ ┣ 📜raft.go
 ┃ ┃ ┗ 📜vote.go
 ┃ ┗ 📂store
 ┃ ┃ ┣ 📜store.go
 ┃ ┃ ┗ 📜store_test.go
 ┣ 📂log
 ┃ ┗ 📜log.go
 ┣ 📜.env
 ┣ 📜.gitignore
 ┣ 📜README.md
 ┣ 📜go.mod
 ┗ 📜go.sum
```
