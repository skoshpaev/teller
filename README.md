![version](https://img.shields.io/badge/version-0.7.1-yellow.svg)
![Go Version](https://img.shields.io/github/go-mod/go-version/skoshpaev/teller)

# Teller

### Get Started in Three Easy Steps

#### 1. Run the Server – It's Super Easy

Start your Teller server with just one command. All you need is a JWT secret, and you're good to go:

```bash
bin/Teller/main --jwt-secret=your_jwt_secret
```

By default, the server runs on port `8080`, but you can easily change it with a simple flag. Your server is now ready to handle real-time connections!

#### 2. Subscribe to a Channel – It's as Simple as 1-2-3

Subscribing to real-time updates is a breeze. With just a few lines of JavaScript, you're connected and ready to receive live messages:

```javascript
const eventSource = new EventSource('http://localhost:8080/subscribe?channel=test-channel', {
    headers: { 'Authorization': 'Bearer your_jwt_token' }
});

eventSource.onmessage = function(event) {
    console.log('New message:', event.data);
};
```

That's it! You're now subscribed to `test-channel` and will receive updates as soon as they're available.

#### 3. Publish a Message – No Hassle, No Fuss

Need to send a message? Teller makes it effortless. Just use `curl` or your favorite HTTP client:

```bash
curl -X POST \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer your_jwt_token" \
     -d '{"channel": "test-channel", "message": {"key": "value"}}' \
     http://localhost:8080/publish
```

OR LIKE IN MERCURE

```bash
curl -X POST \
     -H "Content-Type: application/x-www-form-urlencoded" \
     -H "Authorization: Bearer your_jwt_token" \
     -d "topic=/test-channel&data=Your any message content" \
     "http://localhost:8080/publish"
```

Your message is instantly broadcasted to all subscribers of `test-channel`. Real-time communication has never been this easy!

### Try Teller, it works

Teller is designed for developers who value simplicity and efficiency. Whether you're working on a small project or need a reliable real-time messaging solution, Teller delivers with minimal overhead and maximum performance. Give it a try and see how effortlessly you can add real-time capabilities to your application.



**Teller** is a lightweight Go-based application designed to streamline real-time messaging through the use of Server-Sent Events (SSE). It provides a secure and efficient platform for both publishing and subscribing to messages, leveraging JWT (JSON Web Token) authentication to ensure that only authorized clients can interact with the system.

Teller is ideal for scenarios where real-time updates are critical, and it excels in environments where quick setup, ease of testing, and minimal deployment complexity are desired. Its simplicity makes it particularly suited for Single Page Application (SPA) development and local frontend testing, offering a straightforward, yet powerful, real-time communication tool.

## How It Works

Teller is designed to facilitate real-time messaging between a server and multiple clients through a streamlined process. Below is a simplified explanation of how the core components interact.

### 1. The Server

- **Handles Client Connections**: The server listens for incoming client connections on a specified port. Clients can subscribe to specific channels using Server-Sent Events (SSE) to receive real-time updates.
  
- **Manages Subscriptions**: When a client subscribes to a channel, the server registers the client and creates a communication channel to send updates. The server ensures that only authenticated clients, verified through JWT tokens, can subscribe to or publish on channels.

- **Distributes Messages**: The server distributes messages to all clients subscribed to a particular channel as soon as those messages are published.

### 2. The Publisher

- **Publishes Messages**: The publisher is any client or service that sends a message to a specific channel via the `/publish` endpoint. The publisher must include a valid JWT token for authentication.

- **Specifies the Channel**: Along with the message content, the publisher specifies the target channel to which the message should be sent. This ensures that only clients subscribed to that channel will receive the message.

### 3. The Server Processes and Distributes

- **Receives the Message**: Once the server receives a message from the publisher, it validates the JWT token to ensure the publisher is authorized to send messages.

- **Broadcasts to Subscribers**: After validation, the server broadcasts the message to all clients currently subscribed to the specified channel. These clients receive the message in real-time via SSE.

- **Monitors Connections**: The server continuously monitors active connections and handles disconnections, ensuring that the messaging system remains robust and efficient.

This simple yet powerful workflow makes Teller an ideal choice for scenarios requiring low-latency updates and secure, efficient communication between a server and multiple clients.


### Key Features

- **Real-Time Communication**: Teller utilizes Server-Sent Events (SSE) to push real-time updates from the server to connected clients. This ensures that users receive timely and accurate information as soon as it becomes available, without the overhead of more complex protocols.

- **JWT Authentication**: Teller implements secure JWT-based authentication, allowing only authorized clients to publish or subscribe to specific channels. This ensures that communication remains secure and restricted to trusted users.

- **Scalable and Lightweight**: Built with Go, Teller is designed to be lightweight and highly performant. It can handle numerous concurrent connections with minimal resource usage, making it a perfect fit for resource-constrained environments.

- **Easy Configuration and Deployment**: Teller requires minimal setup and can be easily configured at startup with simple command-line options, such as setting the port and JWT secret. Its simplicity means you don’t need to deploy complex containerized environments to get it running, making it an excellent choice for quick deployments and local development.

- **Ideal for SPA and Frontend Development**: Teller’s straightforward SSE implementation is perfect for testing and developing Single Page Applications (SPAs). It provides an easy-to-use solution for developers needing to simulate or work with real-time data in their frontend applications without the need for heavy infrastructure.

- **Extensible and Flexible**: While simple out of the box, Teller is easily extendable to fit a variety of real-time communication needs. You can adapt it to more complex scenarios as your project evolves, all while keeping the core experience easy and efficient.

Teller is the go-to tool for developers looking for a simple yet powerful solution for real-time messaging, especially in environments where quick setup, minimal overhead, and ease of use are paramount.

## Running the Teller Application

#### 1. Building the Application
Before running the Teller application, you need to build the executable. From the root of your project, run the following commands:

`go build -o bin/Teller/main ./cmd/Teller`

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

#### 4. Checking the Server Status
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

### Possible Errors When Subscribing to a Channel

- **Invalid Request Method (405 Method Not Allowed)**:
  - Occurs when a request to `/subscribe` uses a method other than GET.

- **Missing or Invalid Authorization Header (401 Unauthorized)**:
  - Triggered if the `Authorization` header is missing, malformed, or does not start with "Bearer ".

- **Invalid Token (401 Unauthorized)**:
  - Happens when the JWT token is invalid, expired, or signed with an unexpected algorithm.

- **Missing Channel Parameter (400 Bad Request)**:
  - Occurs if the `channel` parameter is missing in the request URL.

- **Server Error (500 Internal Server Error)**:
  - Occurs if the server encounters an unexpected issue, such as problems with streaming support.

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

### Possible Errors When Publishing Messages

- **Invalid Request Method (405 Method Not Allowed)**:
  - Checked by verifying `r.Method` against `http.MethodPost`.

- **Missing or Invalid Authorization Header (401 Unauthorized)**:
  - Triggered if the `Authorization` header is missing, malformed, or does not start with "Bearer ".

- **Invalid Token (401 Unauthorized)**:
  - Happens when the JWT token is invalid, expired, or signed with an unexpected algorithm.

- **Invalid JSON (400 Bad Request)**:
  - Triggered if the JSON in the request body cannot be parsed.

- **Invalid Channel or Message (400 Bad Request)**:
  - Occurs if the `channel` field is empty or if the `message` field is not valid JSON.

- **Server Error (500 Internal Server Error)**:
  - Occurs if the server encounters an unexpected issue, such as a problem with streaming support.


## Load Testing

For load testing use [k6](https://k6.io/).

### Install k6 (Mac):

```bash
brew install k6
```

### Running tests

```bash
k6 run load_test.js
# ```