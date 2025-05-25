# TCP Communication: Server & Rust Client

This project demonstrates bi-directional TCP communication between a server and a Rust client. The Rust client connects to the TCP server, sends messages, and receives responses from the server in real-time.

## Overview

The project consists of two primary components:

- **TCP Server:** A TCP server (implemented in Go, for example) that listens for incoming connections, sends initial prompts (such as "Enter your name", "Enter your age", "Enter your likes") immediately upon connection, and echoes or processes messages sent by clients.
- **Rust Client:** A Rust-based client that connects to the TCP server, sends user-provided input (via standard input) over the TCP connection, and concurrently listens for messages sent by the server. The client uses the Tokio async runtime to enable concurrent reading and writing.

In this example, the Rust client uses `tokio::select!` to concurrently:
- Read input from the user (stdin) and forward it to the server.
- Listen for incoming messages from the server and print them out.

## Architecture

### TCP Server (Go)
- **Listening Port:** 3000  
- **Behavior:**  
  - Immediately prompts the connecting client for initial information (name, age, likes).
  - After receiving the initial data, the server can echo a welcome message or continue with further conversation.
  - Supports bi-directional messaging for an interactive session.

### Rust Client
- **Communication Protocol:** TCP
- **Features:**
  - Connects to the TCP server at a configurable IP address/port (default: 127.0.0.1:3000).
  - Uses asynchronous I/O (via Tokio) for simultaneous reading and writing.
  - Reads user input from stdin, sends each line as a message.
  - Displays incoming messages from the server immediately on the console.
  - Optionally, supports additional commands such as `exit` to gracefully terminate the connection.

## Prerequisites

### Server (Go)
- Go 1.24+ (or a recent version)
- Docker 

### Client (Rust)
- Rust toolchain (rustc, cargo)
- Tokio crate (with full features) included in `Cargo.toml`
- (Optional) Docker if you plan to containerize the Rust client

## Getting Started

### Running the TCP Server

 **Build & Run Locally:**  
   The project was made with interest for developers with less experience in running monorepos or running these languages, A bash script was created for this purpose `run.sh`. On your terminal, run:
   ```bash
    chmod +x run.sh
    ./run.sh
   ```

This simple command runs a docker postgres, build and runs the Go TCP server and also the runs the Rust client easily.

### Convenient Enough?
