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
        if(GlobalOptions.PlayerClientId == "") return;

        var horizontal = Input.GetAxisRaw("Horizontal");
        var vertical = Input.GetAxisRaw("Vertical");

        serverConnector.AddClientState(new ClientState
        {
            InputX = (sbyte)horizontal,
            InputY = (sbyte)vertical,
            Type = 0,   // type of the state being sent
            ClientId = GlobalOptions.PlayerClientId
        });

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
