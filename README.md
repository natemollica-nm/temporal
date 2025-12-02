## Temporal Go Learning Repo

### Getting Started

**_Install Temporal GO SDK_**

```shell
go get go.temporal.io/sdk
```

**_Install Temporal CLI_**

```shell
brew install temporal
```

**_Start Temporal Server_**

```shell
temporal server start-dev
```

This command starts a local Temporal Service. It starts the Web UI, creates the default Namespace, and uses an in-memory database.

* The Temporal Service will be available on `localhost:7233`.
* The Temporal Web UI will be available at http://localhost:8233.

Leave the local Temporal Service running as you work through tutorials and other projects. 

The `temporal server start-dev` command uses an **_in-memory_** database, so stopping the server will erase all your Workflows and all your Task Queues. If you want to retain those between runs, start the server and specify a database filename using the `--db-filename` option, like this:

```shell
temporal server start-dev --db-filename your_temporal.db
```

---

# Temporal - Get Address from IP - Go

This application demonstrates using Temporal by calling two APIs in sequence.
It fetches the user's IP address and then uses that address to geolocate that user.

You can use the app in two ways:

- Through a web front-end
- Through a JSON POST request

In both cases, you provide a name that's included in the greeting.

## Using the app

The app requires the Temporal Service.

Start the web server to handle API and web requests:

```bash
$ go run server/main.go
```

Now start the Temporal Worker

```bash
$ go run worker/main.go
```

Now visit `http://localhost:4000` and enter your name to run the Workflow.


You can also issue a cURL request to start the Workflow:

```bash
$ curl -X POST http://localhost:4000/api -H "Content-Type: application/json" -d '{"name":"Mike Jones"}'
```

Visit http://localhost:8233 to view the Event History in the Temporal UI.

Disable your internet connection and try again. This time you'll see the Workflow pause. Restore the internet connection and the Workflow completes.