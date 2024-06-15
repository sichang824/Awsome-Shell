# Awesome-Shell

Awesome-Shell是一个功能强大的Shell脚本工具集合，旨在简化日常开发和运维工作。通过使用Awesome-Shell，您可以轻松执行各种常见任务，例如数据库管理、系统监控和自动化部署等。

## AI 助手

Awsome Shell现在接入了 Coze AI Bot，你可以通过 AI 助手快速的解决问题，并实现你的需求

[Awsome Shell - Coze AI Bot](https://www.coze.com/s/Zs8MoWf1c/)

## 安装

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

## 卸载

```sh
cd ~/.Awesome-Shell && ./core/uninstall.sh
```

## 使用

### Awesome Usage

使用 core/usage.sh:

在你的脚本末尾中引入core/usage.sh

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

示例：

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

注意: 用适当的值替换 <ACTION_DESCRIPTION>, <SCRIPT_NAME>, <COMMAND>, 和 <DETAILED_DESCRIPTION>。

### Awesome 颜色

## 基准测试

### 基准测试结果

我们使用`hyperfine`对`bash cli.sh -h`命令进行了基准测试。以下是测试结果：

| 基准命令                          | 平均时间 (ms) ± σ | 用户时间 (ms) | 系统时间 (ms)   | 范围 (ms)        | 运行次数 | 备注                                                                 |
|----------------------------------|-------------------|----------------|------------------|-----------------|------|----------------------------------------------------------------------|
| `hyperfine 'bash cli.sh -h'`     | 5.6 ± 0.4         | 2.9            | 2.1              | 5.1 - 9.8       | 359  | 检测到统计异常值。建议在没有其他程序干扰的安静系统上重新运行此基准测试。使用 '--warmup' 或 '--prepare' 选项可能会有所帮助。 |
| `hyperfine --runs 1000 --warmup 10 'bash cli.sh -h'` | 6.2 ± 1.9      | 3.0            | 2.3              | 5.0 - 33.5      | 1000 | 检测到统计异常值。建议在没有其他程序干扰的安静系统上重新运行此基准测试。使用 '--warmup' 或 '--prepare' 选项可能会有所帮助。 |

从以上结果可以看出，`bash cli.sh -h`命令的执行时间非常短，平均在5-6毫秒之间。尽管在某些测试中检测到统计异常值，但总体性能表现仍然非常优秀。

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

## 贡献

我们欢迎社区的贡献！请阅读我们的[贡献指南](CONTRIBUTING.md)以开始。

## 常见问题

### 如何安装Awesome-Shell？

请按照上面提到的安装步骤进行操作。如果遇到任何问题，请检查日志文件以获取更多详细信息。

### 如何卸载Awesome-Shell？

按照“卸载”部分提到的步骤运行卸载脚本。

### 如何为Awesome-Shell做贡献？

请阅读我们的[贡献指南](CONTRIBUTING.md)以获取详细的贡献说明。

## 联系方式

如有任何问题或反馈，请通过[zhaoanke@163.com](mailto:zhaoanke@163.com)联系我们，或加入我们的[讨论论坛](https://github.com/sichang824/Awsome-Shell/discussions)。
