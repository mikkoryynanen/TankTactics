using System;
using System.Collections.Concurrent;
using System.Net.WebSockets;
using System.Text;
using System.Threading;
using System.Threading.Tasks;
using UnityEditor;
using UnityEngine;

namespace GameServerConnector
{
    [DisallowMultipleComponent]
    public class ServerConnector : MonoBehaviour
    {
        [Header("Settings")]
        [SerializeField] string serverAddress = "";
        [SerializeField] uint bufferSize = 1024;

        ConcurrentQueue<ClientState> _clientStateQueue = new();
        ConcurrentQueue<ServerState> _serverStateQueue = new();

        public string ClientId { get; private set; } = "";

        private ClientWebSocket _ws;
        private ClientState _lastClientState = null;

        async void Start()
        {
            Debug.Log("Starting server connection...");

            _ws = new ClientWebSocket();
            try
            {
                await _ws.ConnectAsync(new Uri($"ws://{serverAddress}:8080/c"), CancellationToken.None);
                Debug.Log("Connected to server");

                var readTask = Task.Run(() => ReadMessages());
                var sendTask = Task.Run(() => SendClientState());
                await Task.WhenAll(readTask, sendTask);
            }
            catch (Exception e)
            {
                Debug.LogError($"Server connection failed. {e.Message}");
            }

            await OnDisconnectAsync();
        }

        private async void OnApplicationQuit()
        {
            await OnDisconnectAsync();
        }

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

        public void AddClientState(sbyte inputX, sbyte inputY, MessageType messageType)
        {
            // Only add non identical messages to queue
            if (_lastClientState != null)
            {
                if (_lastClientState.InputX == inputX && _lastClientState.InputY == inputY)
                {
                    Debug.Log("Identical client state. Do not send to server");
                    return;
                }
            }

            var state = new ClientState
            {
                InputX = inputX,
                InputY = inputY,
                Type = (int)messageType,
                ClientId = ClientId
            };

            _clientStateQueue.Enqueue(state);

            _lastClientState = state;
        }

        public ServerState GetServerState()
        {
            return _serverStateQueue.TryDequeue(out ServerState state) ? state : null;
        }

        async Task SendClientState()
        {
            while (_ws.State == WebSocketState.Open)
            {
                while (_clientStateQueue.TryDequeue(out var state))
                {
                    await SendData(state);
                }
            }

            Debug.Log("Websocket connection closed");
        }

        async Task ReadMessages()
        {
            var buffer = new byte[bufferSize];

            while (_ws.State == WebSocketState.Open || _ws.State == WebSocketState.CloseReceived)
            {
                var result = await _ws.ReceiveAsync(new ArraySegment<byte>(buffer), CancellationToken.None);
                if (result.MessageType == WebSocketMessageType.Close)
                {
                    break;
                }

                var rawMessage = GetData(buffer, result.Count);

                try
                {
                    var baseMessage = JsonUtility.FromJson<BaseServerMessage>(rawMessage);
                    if (baseMessage.type == 0)
                    {
                        this.ClientId = baseMessage.ClientId;
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

        async Task SendData<T>(T data)
        {
            var bytes = Encoding.UTF8.GetBytes(JsonUtility.ToJson(data));
            await _ws.SendAsync(
                new ArraySegment<byte>(bytes),
                WebSocketMessageType.Text,
                true,
                CancellationToken.None);
        }

        string GetData(byte[] buffer, int byteCount)
        {
            return Encoding.UTF8.GetString(buffer, 0, byteCount);
        }
    }
}
