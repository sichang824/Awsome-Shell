
# Awesome-Shell

Awesome-Shell is a powerful collection of shell script tools designed to simplify daily development and operations tasks. With Awesome-Shell, you can easily perform various common tasks such as database management, system monitoring, and automated deployment.

## Install

### coding

```sh
git clone https://e.coding.net/cloudbase-100009281119/Awesome-Shell/Awesome-Shell.git ~/.Awesome-Shell
cd  ~/.Awesome-Shell && ./core/install.sh
```

### github

```sh
git clone https://github.com/sichang824/Awsome-Shell.git ~/.Awesome-Shell
cd  ~/.Awesome-Shell && ./core/install.sh
```

## Uninstall

```sh
cd ~/.Awesome-Shell && ./core/uninstall.sh
```

## Usage

### Awesome Usage

Use core/usage.sh:

Include core/usage.sh at the end of your script.

```sh
source "${AWESOME_SHELL_ROOT}/core/usage.sh" && usage "$@"
```

```sh
  This script can be called with various commands to perform specific tasks.
  Use the following commands during development:
    - entry_<command>: Executes the specified command.
    - main: The main function if defined.
    - -h, --help: Display this help message.
```

Examples:

```sh
# Entry function to: <Messages>
# Usage: ./<SCRIPT_NAME> <COMMAND>
# Description: <Display into usage>
entry_<COMMAND>() {
    echo "This is a template"
}
```

```sh
# Main function example
# Description: Main function to execute the default behavior
main() {
    echo "This is the main function"
}
```

Note: Replace <ACTION_DESCRIPTION>, <SCRIPT_NAME>, <COMMAND>, and <DETAILED_DESCRIPTION> with appropriate values.

### Awesome Colors

## Benchmark

### Benchmark Results

We used `hyperfine` to benchmark the `bash cli.sh -h` command. Here are the results:

| Benchmark Command                | Mean Time (ms) ± σ | User Time (ms) | System Time (ms) | Range (ms)        | Runs | Notes                                                                 |
|----------------------------------|--------------------|----------------|------------------|-------------------|------|----------------------------------------------------------------------|
| `hyperfine 'bash cli.sh -h'`     | 5.6 ± 0.4          | 2.9            | 2.1              | 5.1 - 9.8         | 359  | Statistical outliers were detected. Consider re-running this benchmark on a quiet system without any interferences from other programs. It might help to use the '--warmup' or '--prepare' options. |
| `hyperfine --runs 1000 --warmup 10 'bash cli.sh -h'` | 6.2 ± 1.9          | 3.0            | 2.3              | 5.0 - 33.5        | 1000 | Statistical outliers were detected. Consider re-running this benchmark on a quiet system without any interferences from other programs. It might help to use the '--warmup' or '--prepare' options. |

From the above results, we can see that the execution time of the `bash cli.sh -h` command is very short, averaging between 5-6 milliseconds. Despite detecting statistical outliers in some tests, the overall performance is still excellent.

```sh
❯ hyperfine 'bash cli.sh -h'

Benchmark 1: bash cli.sh -h
  Time (mean ± σ):       5.6 ms ±   0.4 ms    [User: 2.9 ms, System: 2.1 ms]
  Range (min … max):     5.1 ms …   9.8 ms    359 runs

  Warning: Statistical outliers were detected. Consider re-running this benchmark on a quiet system without any interferences from other programs. It might help to use the '--warmup' or '--prepare' options.
```

```sh
❯ hyperfine --runs 1000 --warmup 10 'bash cli.sh -h'
Benchmark 1: bash cli.sh -h
  Time (mean ± σ):       6.2 ms ±   1.9 ms    [User: 3.0 ms, System: 2.3 ms]
  Range (min … max):     5.0 ms …  33.5 ms    1000 runs

  Warning: Statistical outliers were detected. Consider re-running this benchmark on a quiet system without any interferences from other programs. It might help to use the '--warmup' or '--prepare' options.
```

```sh
❯ bash cli.sh -h
Usage: cli.sh [command]

Main:  Main function to execute the default behavior
     ./cli.sh <COMMAND>

Commands:
  mongo_client  Executes the provided MongoDB command using the MongoDB client
     ./cli.sh entry_mongo_client <COMMAND>
  mongo_create_user  Creates a MongoDB user with the provided username, optional role, and optional database (default: admin)
     ./cli.sh mongo_create_user <USERNAME> [<ROLE>] [<DATABASE>]
  mongo_create_db  Creates a MongoDB database with the provided name
     ./cli.sh mongo_create_db <DATABASE_NAME>
  mongo_grant  Grants the specified role on the specified database to the specified user (dbOwner)
     ./cli.sh entry_mongo_grant <USERNAME> <ROLE> <DATABASE>
  mongo_delete_db  Deletes the specified MongoDB database
     ./cli.sh entry_mongo_delete_db <DATABASE_NAME>
  mongo_delete_user  Deletes the specified MongoDB user from the admin database
     ./cli.sh entry_mongo_delete_user <USERNAME>
  mongo_login  Logs in to the MongoDB database using the provided credentials
     ./cli.sh mongo_login <USERNAME> <PASSWORD>
  login  Logs in to the specified MySQL database using the provided credentials
     ./cli.sh login <USERNAME> <PASSWORD>
  create_db  Creates a MySQL database with the provided name
     ./cli.sh create_db <DATABASE_NAME>
  client  Connects to the specified MySQL database using the root user
     ./cli.sh client [<Options>]
  create_user  Creates a MySQL user with the provided credentials
     ./cli.sh create_user <USERNAME> <PASSWORD>
  delete_user  Deletes a MySQL user with the provided name after confirmation
     ./cli.sh delete_user <USERNAME>
  delete_db  Deletes a MySQL database with the provided name after confirmation
     ./cli.sh delete_db <DATABASE_NAME>
  grant  Grants all privileges on the specified database to the specified user
     ./cli.sh grant <DATABASE_NAME> <USERNAME>
```

## Contributing

We welcome contributions from the community! Please read our [contributing guidelines](CONTRIBUTING.md) to get started.

## FAQ

### How do I install Awesome-Shell?

Follow the installation steps mentioned above. If you encounter any issues, please check the log files for more details.

### How do I uninstall Awesome-Shell?

Run the uninstall script as mentioned in the "Uninstall" section.

### How do I contribute to Awesome-Shell?

Please read our [contributing guidelines](CONTRIBUTING.md) for detailed instructions on how to contribute.

## Contact

For any questions or feedback, please reach out to us at [zhaoanke@163.com](mailto:zhaoanke@163.com) or join our [discussion forum](https://github.com/sichang824/Awsome-Shell/discussions).
