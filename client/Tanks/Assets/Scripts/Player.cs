using UnityEngine;

public class Player : MonoBehaviour
{
    [SerializeField] ServerConnector serverConnector;

    [SerializeField] float movementSpeed = 10;
    [SerializeField] float turnSpeed = 10;

    void FixedUpdate()
    {
        var horizontal = Input.GetAxisRaw("Horizontal");
        var vertical = Input.GetAxisRaw("Vertical");

        serverConnector.AddClientState(new ClientState
        {
            InputX = (sbyte)horizontal,
            InputY = (sbyte)vertical,
            Type = 0,   // type of the state being sent
            ClientId = "client" // TODO This will be the UUID we receive from the server
        });

        var serverState = serverConnector.GetServerState();
        if (serverState != null)
        {
            transform.position = new Vector2(serverState.posx, serverState.posy);
        }

        // Move(horizontal, vertical);
    }

    void Move(float horizontal, float vertical)
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
}
