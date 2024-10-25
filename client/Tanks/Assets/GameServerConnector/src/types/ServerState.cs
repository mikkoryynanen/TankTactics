namespace GameServerConnector
{
    [System.Serializable]
    public class ServerState
    {
        public float posx;
        public float posy;
    }

    [System.Serializable]
    public class BaseServerMessage
    {
        public int type;
        public string ClientId;
    }
}