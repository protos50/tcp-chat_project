use std::net::TcpStream;
use std::io::{Write, BufReader, BufRead};
use tauri::command;
use std::sync::{Arc, Mutex};
use std::thread;
use tauri::Manager;

struct AppState {
    stream: Option<TcpStream>,
}

static mut APP_STATE: AppState = AppState { stream: None };

#[command]
fn connect_to_server(address: String, port: u16) -> Result<(), String> {
    let server_address = format!("{}:{}", address, port);
    match TcpStream::connect(&server_address) {
        Ok(stream) => {
            unsafe {
                APP_STATE.stream = Some(stream);
            }
            Ok(())
        }
        Err(e) => Err(format!("Error connecting to server: {}", e)),
    }
}

#[command]
fn disconnect_from_server() -> Result<(), String> {
    unsafe {
        if let Some(ref mut stream) = APP_STATE.stream {
            if let Err(e) = stream.shutdown(std::net::Shutdown::Both) {
                return Err(format!("Error disconnecting from server: {}", e));
            }
            APP_STATE.stream = None;
            Ok(())
        } else {
            Err("Not connected to server.".to_string())
        }
    }
}

#[command]
fn is_connected() -> bool {
    unsafe {
        APP_STATE.stream.is_some()
    }
}

#[command]
fn send_message_to_server(message: String) -> Result<(), String> {
    unsafe {
        if let Some(ref mut stream) = APP_STATE.stream {
            if let Err(e) = stream.write_all(format!("{}\n", message).as_bytes()) {
                return Err(format!("Error sending message: {}", e));
            }
        } else {
            return Err("Not connected to server.".to_string());
        }
    }
    Ok(())
}

fn main() {
    tauri::Builder::default()
        .setup(|app| {
            let app_handle = app.handle(); // Obtener el manejador de la aplicación
            let stream_clone = Arc::new(Mutex::new(TcpStream::connect("localhost:8080").unwrap()));

            // Hilo para recibir mensajes
            let stream_clone_clone = Arc::clone(&stream_clone);
            let app_handle_clone = app_handle.clone(); // Clonar el manejador para el hilo
            thread::spawn(move || {
                let mut buffer = String::new();
                loop {
                    let stream_lock = stream_clone_clone.lock().unwrap();
                    let mut reader = BufReader::new(&*stream_lock);

                    if let Ok(_) = reader.read_line(&mut buffer) {
                        if !buffer.is_empty() {
                            // Emitir el mensaje recibido
                            app_handle_clone.emit_all("message-received", buffer.trim()).unwrap();
                        }
                        buffer.clear();
                    }
                }
            });

            Ok(())
        })
        .invoke_handler(tauri::generate_handler![
            send_message_to_server,
            connect_to_server,
            disconnect_from_server,
            is_connected
        ])
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}
