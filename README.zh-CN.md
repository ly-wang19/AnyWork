# AnyWork

**任何工作，任何智能体，任何机器。**

**简体中文** · [English](README.md) · [日本語](README.ja.md)

[能力地图](CAPABILITIES.zh-CN.md) · [对比评测](docs/BENCHMARK.zh-CN.md) · [路线图](ROADMAP.zh-CN.md) · [质量标准](docs/QUALITY.zh-CN.md) · [支持矩阵](docs/SUPPORT.zh-CN.md) · [贡献指南](CONTRIBUTING.zh-CN.md)

AnyWork 是采用 Apache-2.0 许可证的 AI 智能体工作能力层：包含 24 个原子 Skill，覆盖调研、分析、决策、写作、数据、会议、运营与持续改进；再加 1 个编排器，把复杂需求自动变成正确、经过检查的能力链。

当前仓库是**实验性 Alpha**，不能把“世界第一”当成已经成立的宣传语。AnyWork 要把这个目标变成可检验的工程事实：每项能力、语言、智能体与机器声明，都必须指向可复现证据。

## AnyWork 能做什么

| 工作领域 | 可以获得的结果 | Skill 数量 |
|---|---|---:|
| 定义与规划 | 澄清目标、制定计划、排列工作优先级 | 3 |
| 调研与证据 | 搜集来源、提取信息、核验主张、综合证据 | 4 |
| 分析与决策 | 分析数据、诊断根因、建立情景、比较方案、提出最终决策 | 5 |
| 内容与沟通 | 交付简报、文档、电子表格、演示文稿、消息和本地化内容 | 6 |
| 会议与运营 | 筹备会议、整理决策行动、创建 SOP、自动化重复工作 | 4 |
| 质量与改进 | 审查交付物并开展证据化复盘 | 2 |
| 端到端编排 | 选择、执行、检查并交接最小必要的多 Skill 能力链 | 1 |

可以直接安装 **7 个工作包**之一——通用工作核心、管理与领导力、产品与运营、研究与咨询、市场增长、人力与招聘、完整能力包——也可以自行组合 24 个原子 Skill。详见[完整能力地图、25 项能力和端到端工作流示例](CAPABILITIES.zh-CN.md)。

## 北极星契约

- **任何工作：**提供通用原子技能，以及可组合的工作包与岗位包。
- **任何智能体：**维护一份保守的标准技能源，并为不同宿主提供有文档依据的适配器。
- **任何机器：**提供原生二进制、确定性安装与更新、恢复、诊断和可恢复卸载。
- **三语同一标准：**中文、英文、日文遵守同一结构契约与正式版门槛。
- **默认可信：**强制记录来源、许可证、所需权限和风险等级。

Alpha 的真实状态：24 个原子 Skill 和 1 个编排器全部为 `experimental`；中、英、日内容均为机器起草，尚待人工签署审校；跨智能体行为等价与完整原生机器矩阵尚未验证。详见[声明证据账本](evidence/claims.json)与[质量门槛](docs/QUALITY.zh-CN.md)。

默认文档化适配覆盖 Codex、Claude Code、Gemini CLI、GitHub Copilot 和 OpenCode。Cursor 为实验适配，不进入 `all`。静态评测契约包含 75 个案例、225 个三语提示变体，但不声称这些提示已经在各模型上完成实跑。

## 一条命令开始

macOS 或 Linux：

```bash
curl -fsSL https://raw.githubusercontent.com/ly-wang19/AnyWork/v0.3.0-alpha.1/install.sh | sh -s -- --setup --agent all --lang zh-CN
```

Windows PowerShell：

```powershell
& ([scriptblock]::Create((irm https://raw.githubusercontent.com/ly-wang19/AnyWork/v0.3.0-alpha.1/install.ps1))) -Setup -Agent all -Language zh-CN
```

版本锁定的脚本会识别机器、下载匹配的原生归档、校验已发布的 SHA-256，然后安装 CLI 并执行 `anywork setup`。如果你的安全规则不允许直接执行远程脚本，请先审查 [install.sh](install.sh) 或 [install.ps1](install.ps1)。脚本不会修改 shell 配置，也不会默认授予工具权限。

然后直接看一个真实的端到端任务：

```bash
anywork demo research --lang zh-CN
```

把显示的任务粘贴到 Codex、Claude Code 或其他已安装宿主。CLI 也支持 `--lang en` 和 `--lang ja`；未指定时跟随操作系统语言。如果只需较小工作包，运行 `anywork install essential`。`anywork-cli.py` 保留为 Python 标准库参考实现与测试基准。

## 仓库结构

```text
skills/       聚焦单一任务的标准 Skills
registry/     工作包，以及来源、语言审校、能力、风险和版本元数据
adapters/     不同智能体的安装规则
evals/        多语言行为与质量评测
evidence/     机器可读的声明状态与证据
*.go          无第三方依赖的原生 AnyWork CLI
src/          Python 参考实现
schemas/      机器可读的契约
scripts/      可复现发布构建
tests/        注册表、证据、三语和安装器回归测试
CAPABILITIES*.md  24 个原子 Skill、1 个编排器、7 个工作包和工作流
ROADMAP*.md   从 Alpha 到稳定生态的证据门槛路线图
```

## 路线图

`v0.3` 的“最后一公里”体验已加入一键初始化、自动多 Skill 编排、可运行 Demo 和保持提示完全一致的[安装前后对比评测](docs/BENCHMARK.zh-CN.md)。剩余证据门槛是三语人工审校、模型实跑结果和签名发行。详见[完整的证据门槛路线图](ROADMAP.zh-CN.md)。

## 当前状态

AnyWork 尚未正式发布。目录结构检查、本地安装生命周期测试、三系统 CI、六目标交叉编译和 GitHub 公共仓库已经验证；母语人工审校、完整模型实跑评测、每个系统与架构上的原生执行和签名发布仍是发布门槛，不能被描述成已完成。

## 许可证

项目采用 Apache-2.0。第三方目录条目保留其原始许可证与署名；许可证不兼容的项目只提供链接，不复制内容。
