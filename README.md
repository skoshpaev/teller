# Teller
Teller is a Go-based application designed to facilitate real-time messaging through the use of Server-Sent Events (SSE). It provides a robust and secure platform for publishing and subscribing to messages, utilizing JWT (JSON Web Token) authentication for secure communication. Teller is ideal for scenarios where real-time updates and efficient communication between a server and multiple clients are crucial.

## Key Features:
* Real-Time Communication: Utilizes Server-Sent Events (SSE) to push real-time updates from the server to connected clients, ensuring that users receive timely and accurate information as it happens.

* JWT Authentication: Implements secure authentication using JWTs, allowing only authorized clients to publish or subscribe to specific channels.

* Scalable and Lightweight: Built with Go, Teller is designed to be lightweight and highly performant, capable of handling numerous concurrent connections with minimal resource usage.

* Configurable and Extendable: Teller offers easy configuration options for server settings like port number and JWT secret, and it can be extended to fit various real-time communication needs.

## Running the Teller Application

#### 1. Building the Application
Before running the Teller application, you need to build the executable. From the root of your project, run the following commands:

`go build -o bin/Teller/main`

This command compiles your Go code and generates the executable at bin/Teller/main.
#### 2. Running the Application
The Teller application requires one mandatory parameter (jwt-secret) and supports two optional parameters (port and jwt-ttl).
Mandatory Parameter:

`--jwt-secret`: The secret key used to sign JWT tokens. This parameter must be provided; otherwise, the application will not start.

Optional Parameters:

`--port`: The port number on which the server will listen. If not specified, the default port 8080 will be used.

`--jwt-ttl`: The time-to-live (TTL) for JWT tokens, specified in hours. If not specified, the default TTL is 1 hour.


##### Example 1: Running with the mandatory jwt-secret
To run the application with only the mandatory jwt-secret parameter:

`bin/Teller/main --jwt-secret=your_jwt_secret`


This command will start the server on the default port 8080, with JWT tokens having a TTL of 1 hour.
##### Example 2: Running with all parameters
To specify the port and JWT TTL:

`bin/Teller/main --jwt-secret=your_jwt_secret --port=9090 --jwt-ttl=2`

In this example:
- The server will start on port 9090.
- JWT tokens will have a TTL of 2 hours.

#### 3. Handling Common Errors3. Handling Common Errors
ERROR: MISSING JWT-SECRET

If you attempt to run the Teller application without specifying the **--jwt-secret** parameter, the application will fail to start, and you will see the following error message:

`JWT Secret Key must be provided`

**Solution**: Ensure that you provide the --jwt-secret parameter when starting the application.

ERROR: INVALID JWT TTL VALUE

If you provide an invalid value for the --jwt-ttl parameter, the application will either use the default value (1 hour) or produce an error depending on how the application is configured to handle such cases.

**Solution:** Ensure that the --jwt-ttl value is a positive integer representing the number of hours.

#### 4. Checking the Server Status4. Checking the Server Status
Once the server is running, you can check its status by trying to access it via a web browser or a command-line tool like curl:
`curl http://localhost:8080`
If the server is running correctly, you should receive a response, or in the case of a properly configured route, you may see a message or data returned from the server.


## Subscribing to Channels
#### Overview of the Subscription Process
The subscription process in Teller involves the following steps:

1. Client Sends a Subscription Request: The client sends an HTTP GET request to the /subscribe endpoint, specifying the channel they wish to subscribe to. The request must include a valid JWT token in the Authorization header.
2. Server Streams Messages: Upon a successful subscription, the server streams messages to the client as they are published to the specified channel.
3. Real-Time Updates: The client receives real-time updates in the form of SSE messages as long as the connection remains open.

##### Subscribing to a Channel Using JavaScript
To subscribe to a channel using JavaScript, follow these steps:

- Obtain a JWT Token: First, obtain a JWT token from your backend (this step depends on your authentication flow).

- Create a Subscription to the Channel:

- Use the following JavaScript code to subscribe to a channel and listen for updates:

```javascript
async function subscribeToChannel(channel) {
    const token = 'your_jwt_token'; // Replace with your actual JWT token
    const url = new URL(`http://localhost:8080/subscribe`);
    url.searchParams.append('channel', channel);

    const eventSource = new EventSource(url, {
        headers: {
            'Authorization': `Bearer ${token}`
        }
    });

    eventSource.onmessage = function(event) {
        console.log('New message:', event.data);
    };

    eventSource.onerror = function(event) {
        console.error('Subscription error:', event);
    };
}

// Call the function to subscribe to a specific channel
subscribeToChannel('test-channel');
```

- Expected Messages:
 - New Message: When a new message is published to the channel, you will see an output in the console like:
`data: {"key":"value2"}`

##### Subscribing to a Channel Using curlSubscribing to a Channel Using curl
You can also subscribe to a channel using curl from the command line. This is useful for testing or simple scripts.

```bash
curl -N -H "Accept: text/event-stream" \
     -H "Authorization: Bearer your_jwt_token" \
     "http://localhost:8080/subscribe?channel=test-channel"

```

## Publishing Messages in Teller

Teller allows you to publish messages to specific channels, which are then delivered in real-time to all clients subscribed to those channels. Below are detailed instructions on how to publish messages using both JavaScript and `curl`.

### 1. Publishing a Message Using JavaScript

To publish a message to a specific channel from a JavaScript application, follow these steps:

#### Step 1: Obtain a JWT Token

Before you can publish a message, you need a valid JWT token. This token should be obtained from your backend (depending on your authentication flow).

#### Step 2: Publish the Message

Use the following JavaScript code to publish a message to a channel:

```javascript
async function publishMessage(channel, message) {
    const token = 'your_jwt_secret'; // Replace with your actual JWT token
    const url = 'http://localhost:8080/publish';

    const payload = {
        channel: channel,
        message: message
    };

    const response = await fetch(url, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify(payload)
    });

    if (response.ok) {
        console.log('Message published successfully');
    } else {
        console.error('Failed to publish message', await response.text());
    }
}

// Example usage
publishMessage('test-channel', { key: 'value' });
```

#### Expected Response

- **Success (`200 OK`)**: The message is published, and the server responds with:

  ```json
  "Message received successfully"
  ```

### 2. Publishing a Message Using `curl`

You can also publish a message to a channel using `curl`, which is particularly useful for quick testing or automation scripts.

#### Step 1: Run the Following `curl` Command

```bash
curl -X POST \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer your_jwt_token" \
     -d '{"channel": "test-channel", "message": {"key": "value"}}' \
     http://localhost:8080/publish
```

Replace `your_jwt_token` with a valid JWT token and adjust the `channel` and `message` fields as needed.

#### Expected Response

- **Success (`200 OK`)**: The message is published, and the server responds with:

  ```json
  "Message received successfully"
  ```

### 3. Server-Side Implementation Details

Here’s how the Teller application handles message publication:

- **Endpoint**: `/publish`
- **HTTP Method**: `POST`
- **Request Payload**:
  - `channel`: The channel name to which the message is being published (required).
  - `message`: The message content in JSON format (required).

- **JWT Authentication**: The request must include a valid JWT token in the `Authorization` header. The server validates this token before processing the message.

- **Success Response**:
  - `200 OK`: Message published successfully.

### 4. Example Workflow

1. **Publish a Message**: Use JavaScript or `curl` to publish a message to a specific channel (e.g., `test-channel`).
2. **Subscribe to the Channel**: Ensure that clients are subscribed to the channel to receive the message in real-time.
3. **Observe the Interaction**: The published message will be broadcast to all clients subscribed to the `test-channel`.

