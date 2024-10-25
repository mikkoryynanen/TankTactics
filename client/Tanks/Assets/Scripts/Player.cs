using GameServerConnector;
using UnityEngine;

public class Player : MonoBehaviour
{
    [SerializeField] ServerConnector serverConnector;

    [SerializeField] float movementSpeed = 10;
    [SerializeField] float turnSpeed = 10;
    [SerializeField] float interpolationSpeed = .1f;
    [SerializeField] bool serverOnlyMovement = false;

    void Update()
    {
        if (serverConnector == null || serverConnector.ClientId == "")
        {
            return;
        }

        var horizontal = Input.GetAxisRaw("Horizontal");
        var vertical = Input.GetAxisRaw("Vertical");

        serverConnector.AddClientState(
            (sbyte)horizontal,
            (sbyte)vertical,
            (int)MessagesType.ClientState
        );

        if (!serverOnlyMovement)
        {
            LocalMove(horizontal, vertical);
        }
        CorrectPosition();
    }

    void LocalMove(float horizontal, float vertical)
    {
        transform.position = new Vector2(
            transform.position.x + horizontal * Time.deltaTime * movementSpeed,
            transform.position.y + vertical * Time.deltaTime * movementSpeed
        );

        var turnInputAxis = Input.GetAxisRaw("Turn");
        transform.Rotate(new Vector3(
         0, 0,
             turnInputAxis * Time.deltaTime * turnSpeed
        ));
    }

    void CorrectPosition()
    {
        var serverState = serverConnector.GetServerState();
        if (serverState != null)
        {
            var newPosition = new Vector2(serverState.posx, serverState.posy);
            transform.position = Vector3.Lerp(transform.position, newPosition, interpolationSpeed);
        }
    }
}
