<script lang="ts">
  import { onMount } from 'svelte';
  import { invoke } from '@tauri-apps/api';
  import { listen } from '@tauri-apps/api/event';

  let ip = 'localhost';
  let port = 8080;
  let message = '';
  let messages: string[] = [];
  let connected = false;
  let connectionMessage = '';

  const connectToServer = async () => {
    try {
      await invoke('connect_to_server', { address: ip, port });
      connectionMessage = 'Conectado al servidor';
      connected = true;
      console.log(connectionMessage);

      // Escucha los mensajes del servidor
      const unlisten = await listen('message-received', (event) => {
        const serverMessage = event.payload as string;
        if (!serverMessage.startsWith('Tú:')) {
          messages = [...messages, serverMessage];
        }
      });

      return () => {
        unlisten();
      };
    } catch (error) {
      connectionMessage = 'Error al conectar al servidor: ' + error;
      console.error(connectionMessage);
    }
  };

  const disconnectFromServer = async () => {
    try {
      await invoke('disconnect_from_server');
      connectionMessage = 'Desconectado del servidor';
      connected = false;
      console.log(connectionMessage);
    } catch (error) {
      connectionMessage = 'Error al desconectar del servidor: ' + error;
      console.error(connectionMessage);
    }
  };

  const sendChatMessage = async () => {
    if (message.trim()) {
      try {
        await invoke('send_message_to_server', { message });
        messages = [...messages, `Tú: ${message}`]; 
        message = ''; 
      } catch (error) {
        console.error('Error al enviar el mensaje:', error);
      }
    }
  };

  const handleKeyPress = (event: KeyboardEvent) => {
    if (event.key === 'Enter') {
      sendChatMessage();
    }
  };
</script>

<style>
  body {
    background-color: #121212;
    color: #e0e0e0;
    font-family: Arial, sans-serif;
    margin: 0;
    padding: 0;
    display: flex;
    justify-content: center;
    align-items: center;
    height: 100vh;
  }

  .container {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 20px;
    background-color: #1e1e1e;
    border-radius: 10px;
    box-shadow: 0 0 10px rgba(0, 0, 0, 0.5);
  }

  h1 {
    text-align: center;
    color: #bb86fc;
    margin-bottom: 20px;
  }

  input[type="text"], input[type="number"] {
    background-color: #2c2c2c;
    border: 1px solid #bb86fc;
    color: #e0e0e0;
    padding: 10px;
    margin: 5px;
    border-radius: 5px;
    width: 300px;
  }

  button {
    background-color: #bb86fc;
    border: none;
    color: #121212;
    padding: 10px 20px;
    margin: 5px;
    border-radius: 5px;
    cursor: pointer;
  }

  button:disabled {
    background-color: #3a3a3a;
    cursor: not-allowed;
  }

  p {
    margin: 10px 0;
  }

  .chat-box {
    background-color: #2c2c2c;
    border: 1px solid #bb86fc;
    border-radius: 5px;
    padding: 10px;
    width: 80%;
    max-width: 600px;
    height: 400px;
    overflow-y: auto;
    margin-top: 20px;
  }

  ul {
    list-style-type: none;
    padding: 0;
  }

  li {
    background-color: #3a3a3a;
    border: 1px solid #bb86fc;
    border-radius: 5px;
    padding: 10px;
    margin: 5px 0;
  }
</style>

<div class="container">
  <h1>Chat</h1>
  <input type="text" bind:value={ip} placeholder="IP del servidor" />
  <input type="number" bind:value={port} placeholder="Puerto del servidor" />
  <button on:click={connectToServer} disabled={connected}>Conectar</button>
  <button on:click={disconnectFromServer} disabled={!connected}>Desconectar</button>
  <p>{connectionMessage}</p>

  <input type="text" bind:value={message} placeholder="Escribe un mensaje..." on:keypress={handleKeyPress} />
  <button on:click={sendChatMessage} disabled={!connected}>Enviar</button>

  <div class="chat-box">
    <ul>
      {#each messages as msg}
        <li>{msg}</li>
      {/each}
    </ul>
  </div>
</div>
