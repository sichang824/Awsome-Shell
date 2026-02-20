
# Contributing to Awesome Shell

我们非常欢迎社区贡献者的参与！无论是修复bug、添加功能、改进文档，还是提供反馈和建议，都对项目的发展至关重要。

## 如何贡献

### 报告问题

1. 确保问题尚未在[Issue Tracker](https://github.com/sichang824/Awsome-Shell/issues)中报告。
2. 创建一个新的Issue，并提供尽可能详细的信息，包括：
   - 问题描述
   - 重现步骤
   - 预期行为和实际行为
   - 环境信息（操作系统、Shell版本等）

### 提交代码

1. Fork项目仓库到你的GitHub账户。
2. 克隆你Fork的仓库到本地：

   ```sh
   git clone https://github.com/<your-username>/Awsome-Shell.git
   ```

3. 创建一个新的分支：

   ```sh
   git checkout -b my-feature-branch
   ```

4. 在新的分支上进行开发，并确保代码符合项目的编码规范。
5. 提交代码并推送到你的Fork仓库：

   ```sh
   git add .
   git commit -m "描述你的更改"
   git push origin my-feature-branch
   ```

6. 创建一个Pull Request，并描述你的更改内容和目的。

### 编码规范

- 请确保代码风格一致，并遵循现有代码的格式。
- 提交前请运行所有测试并确保它们通过。
- 如果添加了新功能，请更新相应的文档。

### 代码审查

所有的Pull Request都会经过代码审查。请耐心等待反馈，并根据建议进行修改。代码审查的目的是确保代码质量和项目的一致性。

### 测试

在提交代码之前，请确保所有现有测试通过，并为新功能添加相应的测试。可以使用以下命令运行测试：

```sh
# 运行所有测试
./run_tests.sh
```

### 文档

如果你的更改影响了项目的使用方式，请相应地更新文档。文档文件位于`docs/`目录下。

## 社区

- 如果有任何问题或建议，可以在[讨论区](https://github.com/sichang824/Awsome-Shell/discussions)与我们交流。

感谢你的贡献！

Awesome Shell 团队
