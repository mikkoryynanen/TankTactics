using System;
using System.Collections.Concurrent;
using System.Net.WebSockets;
using System.Text;
using System.Threading;
using System.Threading.Tasks;
using Unity.Mathematics;
using Unity.VisualScripting;
using UnityEngine;

public class ServerConnector : MonoBehaviour
{
    ConcurrentQueue<ClientState> _clientStateQueue = new();
    ConcurrentQueue<ServerState> _serverStateQueue = new();

    private ClientWebSocket _ws;

    async void Start()
    {
        Debug.Log("Starting server connection...");

        _ws = new ClientWebSocket();
        try
        {
            await _ws.ConnectAsync(new Uri("ws://localhost:8080/c"), CancellationToken.None);
            Debug.Log("Connected to server");

            var readTask = Task.Run(() => ReadMessages(_ws));
            var sendTask = Task.Run(() => SendClientState(_ws));
            await Task.WhenAll(readTask, sendTask);
        }
        catch (System.Exception e)
        {
            Debug.LogError($"Server connection failed. {e.Message}");
        }

        await OnDisconnectAsync();
    }

    private async void OnApplicationQuit()
    {
        await OnDisconnectAsync();
    }

    // Called when the GameObject is destroyed (e.g., if manually removed or scene changes)
    private async void OnDestroy()
    {
        await OnDisconnectAsync();
    }

    private async Task OnDisconnectAsync()
    {
        if (_ws != null && _ws.State == WebSocketState.Open)
        {
            await _ws.CloseAsync(WebSocketCloseStatus.NormalClosure, "Closing", CancellationToken.None);
            Debug.Log("Connection to server closed");
        }
    }

    public void AddClientState(ClientState state)
    {
        // Only add non identical messages to queue
        // if (_clientStateQueue.TryPeek(out var peekState))
        // {
        //     if (peekState.InputX == state.InputX && peekState.InputY == state.InputY)
        //     {
        //         Debug.Log("Identical client state. Do not send to server");
        //         return;
        //     }
        // }
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
        while (ws.State == WebSocketState.Open)
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

        Debug.Log("Websocket connection closed");
    }

    async Task ReadMessages(ClientWebSocket ws)
    {
        bool metadataReceived = false;
        var buffer = new byte[1024];

        while (ws.State == WebSocketState.Open || ws.State == WebSocketState.CloseReceived)
        {
            var result = await ws.ReceiveAsync(new ArraySegment<byte>(buffer), CancellationToken.None);

            if (result.MessageType == WebSocketMessageType.Close)
            {
                break;
            }

            var rawMessage = Encoding.UTF8.GetString(buffer, 0, result.Count);
            try
            {
                if (!metadataReceived)
                {
                    var playerMetadata = JsonUtility.FromJson<PlayerMetadata>(rawMessage);
                    Debug.Log($"Received clientId {playerMetadata.ClientId}");
                    GlobalOptions.PlayerClientId = playerMetadata.ClientId;
                    metadataReceived = true;
                }
                else
                {
                    var serverState = JsonUtility.FromJson<ServerState>(rawMessage);
                    _serverStateQueue.Enqueue(serverState);
                }
            }
            catch (Exception e)
            {
                Debug.LogError($"Failed to parse message {e.Message}");
            }
        }

        Debug.Log("Websocket connection closed");
    }
}
