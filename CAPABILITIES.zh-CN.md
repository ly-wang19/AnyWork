# AnyWork 能力地图

**简体中文** · [English](CAPABILITIES.md) · [日本語](CAPABILITIES.ja.md)

AnyWork 为 AI 智能体提供一套可重复执行的工作方法，把模糊任务转成经过检查、可以直接使用的交付物。当前 Alpha 包含 **24 个原子 Skill、1 个端到端编排器和 7 个可安装工作包**。原子 Skill 可单独使用；编排器会为复杂成果选择并执行最小必要能力链。

## 能力总览

| 工作领域 | 可以获得的结果 | 包含的 Skill |
|---|---|---|
| 定义与规划 | 可检验的任务约定、可执行计划或经过排序的工作列表 | [澄清结果](skills/clarify-outcome/SKILL.md)、[制定工作计划](skills/plan-work/SKILL.md)、[排列工作优先级](skills/prioritize-work/SKILL.md) |
| 调研与证据 | 可追溯的来源清单、结构化信息、事实核验结论或证据综合 | [搜集来源](skills/gather-sources/SKILL.md)、[提取结构化信息](skills/extract-structured-info/SKILL.md)、[核验主张](skills/verify-claims/SKILL.md)、[综合多方资料](skills/synthesize-sources/SKILL.md) |
| 分析与决策 | 可复算分析、根因报告、透明情景、方案比较或最终决策建议 | [分析数据](skills/analyze-data/SKILL.md)、[诊断根因](skills/diagnose-root-cause/SKILL.md)、[建立情景模型](skills/model-scenarios/SKILL.md)、[比较方案](skills/compare-options/SKILL.md)、[制定决策](skills/make-decision/SKILL.md) |
| 内容与沟通 | 决策简报、长文档、可审计工作簿、演示文稿、可直接发送的消息或本地化内容 | [撰写简报](skills/write-brief/SKILL.md)、[起草文档](skills/draft-document/SKILL.md)、[构建电子表格](skills/build-spreadsheet/SKILL.md)、[设计演示文稿](skills/design-presentation/SKILL.md)、[撰写消息](skills/write-message/SKILL.md)、[本地化内容](skills/localize-content/SKILL.md) |
| 会议与运营 | 会议设计、决策行动清单、可执行 SOP 或安全自动化方案 | [筹备会议](skills/prepare-meeting/SKILL.md)、[提取决策与行动](skills/capture-decisions-actions/SKILL.md)、[创建 SOP](skills/create-sop/SKILL.md)、[自动化重复工作](skills/automate-routine/SKILL.md) |
| 质量与改进 | 交付物审查，或包含责任人和改进实验的证据化复盘 | [审查交付物](skills/review-deliverable/SKILL.md)、[开展复盘](skills/run-retrospective/SKILL.md) |
| 端到端编排 | 一个已检查的多阶段成果，带明确交接、证据、不确定性和权限关卡 | [编排工作](skills/orchestrate-work/SKILL.md) |

## 工作包

工作包是按常见工作场景整理的起点。不同工作包会有意复用同一个原子能力，因为同一能力可以服务多个岗位。

| 工作包 ID | 适用工作 | Skill 数量 | 安装示例 |
|---|---|---:|---|
| `essential` | 日常调研、写作、协同、审查与编排 | 10 | `anywork install essential` |
| `manager-leadership` | 通用核心加决策、优先级、情景、会议与复盘 | 18 | `anywork install manager-leadership` |
| `product-operations` | 通用核心加产品发现、优先级、交付体系与持续改进 | 19 | `anywork install product-operations` |
| `research-consulting` | 通用核心加证据密集型分析、建议、报告与演示 | 17 | `anywork install research-consulting` |
| `go-to-market` | 通用核心加市场证据、客户决策、内容沟通与执行 | 21 | `anywork install go-to-market` |
| `people-recruiting` | 通用核心加结构化、公平、可审计的人力与招聘工作 | 20 | `anywork install people-recruiting` |
| `complete` | 当前版本全部官方原子 Skill 与编排器 | 25 | `anywork install complete` |

可添加 `--agent codex`、`--agent claude-code`、`--agent gemini-cli`、`--agent github-copilot` 或 `--agent opencode` 指定宿主；使用 `--scope project` 或 `--scope user` 选择项目级或用户级安装；正式写入前可先使用 `--dry-run`。

`anywork setup` 会一次为全部文档化宿主安装 `complete`。`anywork demo research`、`anywork demo meeting` 和 `anywork demo automation` 会输出可直接粘贴的三语编排任务。

## 端到端工作流示例

### 从调研到建议

`orchestrate-work` → 必要时 `clarify-outcome` → `gather-sources` → `verify-claims` → `synthesize-sources` → `compare-options` → `make-decision` → `write-brief` → `review-deliverable`

结果：形成一份可以直接决策的建议，并保留来源、不确定性、备选方案和决策理由的追溯链路。

### 从数据到高管汇报

`orchestrate-work` → `extract-structured-info` → `analyze-data` → `model-scenarios` → `build-spreadsheet` → `design-presentation` → `review-deliverable`

结果：交付可复算模型与面向受众的汇报叙事，并明确区分证据和假设。

### 从会议到责任落实

`orchestrate-work` → `prepare-meeting` → 会议发生 → `capture-decisions-actions` → `plan-work` → `prioritize-work` → `review-deliverable`

结果：先完成聚焦的会议，再形成明确决策、负责人、截止时间、依赖和验收检查。

### 从重复运营到安全自动化

`orchestrate-work` → `diagnose-root-cause` → `create-sop` → `automate-routine` → 受控运行 → `run-retrospective`

结果：形成包含异常处理、人工关卡、恢复机制和改进衡量的受控运营流程。

### 全球化内容交付

`orchestrate-work` → `write-brief` → `draft-document` → `localize-content` → `write-message` → `review-deliverable`

结果：产出可用于具体渠道的英文、简体中文和日文内容，同时保留已经批准的含义与术语。

## 能力边界

AnyWork 安装的是工作方法，不是凭据或静默集成。Skill 只能使用宿主智能体已经具备的工具与访问权限。外部写入、发消息、购买、删除等有后果的操作，仍需明确授权和可核验目标。在完成[质量门槛](docs/QUALITY.zh-CN.md)和[路线图](ROADMAP.zh-CN.md)证据之前，24 个原子 Skill 和编排器全部保持 `experimental` 状态。
