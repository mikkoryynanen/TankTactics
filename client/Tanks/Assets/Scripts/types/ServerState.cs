[System.Serializable]
public class ServerState
{
    public float posx;
    public float posy;
}

[System.Serializable]
public class BaseServerState
{
    public int MessageType;
    public string ClientId;
}
