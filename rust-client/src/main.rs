use tokio::{
    io::{self, AsyncBufReadExt, AsyncWriteExt, BufReader, stdin},
    net::TcpStream,
    time::{Duration, sleep},
};

#[tokio::main]
async fn main() -> io::Result<()> {
    let stream = loop {
        match TcpStream::connect("127.0.0.1:3000").await {
            Ok(s) => break s,
            Err(_) => {
                println!("Server unreachable!, Retrying in 2 seconds....");
                sleep(Duration::from_secs(2)).await;
            }
        }
    };

    println!("Connected to the Go server....");

    let (reader, mut writer) = stream.into_split();
    let mut server_reader = BufReader::new(reader).lines();
    let mut stdin_reader = BufReader::new(stdin()).lines();

    loop {
        tokio::select! {
            Ok(Some(user_input)) = stdin_reader.next_line() => {
                if user_input.trim() == "exit" {
                    println!("Goodbye...");
                    break;
                }

                writer.write_all(user_input.as_bytes()).await?;
                writer.write_all(b"\n").await?;
            }

            Ok(Some(server_msg)) = server_reader.next_line() => {
                println!("Server: {}", server_msg);
            }

            else => {
                println!("Disconnected from the Go server....");
                break;
            }
        }
    }

    Ok(())
}

