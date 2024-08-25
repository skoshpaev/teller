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
