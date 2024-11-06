import websocket
import threading
import time

# WebSocket endpoint URL
ws_url = "ws://127.0.0.1:3001/ws"

# Callback to handle received .flac data
def on_message(ws, message):
    # Write incoming flac data to a file
    with open("output3.flac", "ab") as flac_file:
        flac_file.write(message)
    print("Received FLAC data chunk.")

# Callback for WebSocket open
def on_open(ws):
    def send_audio():
        # Open the .wav file in binary mode
        with open("file_example_WAV_2MG.wav", "rb") as wav_file:
            chunk_size = 8092  # Define chunk size for sending
            while True:
                chunk = wav_file.read(chunk_size)
                if not chunk:
                    break
                # Send the .wav chunk over WebSocket
                ws.send(chunk, opcode=websocket.ABNF.OPCODE_BINARY)
                time.sleep(0.1)  # Adjust sleep to control stream rate

        # Close the WebSocket after sending all data
        ws.close()

    # Start a new thread for sending audio data
    threading.Thread(target=send_audio).start()

# Initialize WebSocket connection
ws = websocket.WebSocketApp(ws_url, on_message=on_message, on_open=on_open)

# Run the WebSocket client
ws.run_forever()
