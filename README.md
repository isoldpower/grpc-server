# User Manual

## Project structure

## Command Line Interface (CLI)
The project is driven by the CLI utilizing the Cobra library for easier set up and maintenance. 
All the commands specifications and other information related **exclusively to CLI** is stored in
the `./cmd` project folder. <br>
In fact, this folder is a **view layer** for the CLI, while all the processes occur in 
the microservice-related folders. See [structure section](#project-structure) for more information.

### Supported Commands
Here is the list of supported commands

#### Global flags
To get detailed information on the commands, you can always use the
following structure:
```
go run main.go [command] --help
```

#### Help Command
This command utilizes the default help command defined by `Cobra library`. <br>
To run this command, use the following line:
```bash
go run main.go --help
```

#### Run Command
This command runs the whole application as a single app in the
right order so that all the integrations through APIs are being run correctly. <br> 
To run this command, use the following line:
```bash
go run main.go run
```

#### Service-specific commands
Those commands target at specific microservice and run commands at a service.
For example, by using `go run main.go orders run`, you will run only `orders` microservice
and ignore the possible dependence on other microservices. <br>
Global pattern:
```bash
go run main.go [service] [command]
```
Example:
```bash
go run main.go orders migrate
```

## Future plans:
1) minimize goroutines memory leak
2) add title to HTTP and gRPC servers
3) set up database connection