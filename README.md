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