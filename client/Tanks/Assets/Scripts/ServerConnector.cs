using System;
using System.Collections.Concurrent;
using System.Linq;
using System.Net.WebSockets;
using System.Resources;
using System.Text;
using System.Threading;
using System.Threading.Tasks;
using UnityEngine;

public class ServerConnector : MonoBehaviour
{
    ConcurrentQueue<ClientState> _clientStateQueue = new();
    ConcurrentQueue<ServerState> _serverStateQueue = new();

    bool _isConnected = false;

    async void Start()
    {
        Debug.Log("Starting server connection...");

        var ws = new ClientWebSocket();
        try
        {
            await ws.ConnectAsync(new Uri("ws://localhost:8080/c"), CancellationToken.None);
            Debug.Log("Connected to server");

            _isConnected = true;

            var readTask = Task.Run(() => ReadMessages(ws));
            var sendTask = Task.Run(() => SendClientState(ws));
            await Task.WhenAll(readTask, sendTask);
        }
        catch (System.Exception e)
        {
            Debug.LogError($"Server connection failed. {e.Message}");
        }

        await ws.CloseAsync(WebSocketCloseStatus.NormalClosure, "Closing", CancellationToken.None);
        Debug.Log("Connection to server closed");
    }

    public void AddClientState(ClientState state)
    {
        // Only add non identical messages to queue
        if (_clientStateQueue.TryPeek(out var peekState))
        {
            if (peekState.InputX == state.InputX && peekState.InputY == state.InputY)
            {
                return;
            }
        }
        _clientStateQueue.Enqueue(state);
        // var json = JsonUtility.ToJson(state);
        // Debug.Log($"Adding state: {json}");
    }

    public ServerState GetServerState()
    {
        return _serverStateQueue.TryDequeue(out ServerState state) ? state : null;
    }

    async Task SendClientState(ClientWebSocket ws)
    {
        while (_isConnected)
        {
            while (_clientStateQueue.TryDequeue(out var state))
            {
                // Debug.Log($"Sending state: {JsonUtility.ToJson(state)}");
                var bytes = Encoding.UTF8.GetBytes(JsonUtility.ToJson(state));
                await ws.SendAsync(
                    new ArraySegment<byte>(bytes),
                    WebSocketMessageType.Text,
                    true,
                    CancellationToken.None);
            }
        }
    }

    async Task ReadMessages(ClientWebSocket ws)
    {
        var buffer = new byte[1024];
        while (true)
        {
            var result = await ws.ReceiveAsync(new ArraySegment<byte>(buffer), CancellationToken.None);

            if (result.MessageType == WebSocketMessageType.Close)
            {
                break;
            }

            var rawMessage = Encoding.UTF8.GetString(buffer, 0, result.Count);
            Debug.Log($"rawMessage: {rawMessage}");
            var serverState = JsonUtility.FromJson<ServerState>(rawMessage);
            Debug.Log($"room state x: {serverState.posx} y: {serverState.posy}");

            _serverStateQueue.Enqueue(serverState);
        }
    }
}
