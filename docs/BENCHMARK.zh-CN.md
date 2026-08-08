# AnyWork 安装前后对比评测

**简体中文** · [English](BENCHMARK.md) · [日本語](BENCHMARK.ja.md)

AnyWork 提供可复现的 A/B 实跑工具，让效果改进可以被测量，而不是靠想象。工具会把字节完全一致的提示分别发给干净基线智能体和安装 AnyWork 后的同一智能体，保存原始输出与耗时，再生成成对审阅表。它刻意不伪造自动质量分数。

## 运行一次小型对比

先用干净配置或干净 Skills 目录运行智能体：

```bash
python3 scripts/benchmark.py run \
  --mode baseline \
  --command 'your-agent --prompt {prompt}' \
  --language zh-CN \
  --limit 3 \
  --output results/baseline.json
```

然后为同一智能体安装 AnyWork，保持模型和其他设置不变，运行相同案例：

```bash
python3 scripts/benchmark.py run \
  --mode anywork \
  --command 'your-agent --prompt {prompt}' \
  --language zh-CN \
  --limit 3 \
  --output results/anywork.json
```

生成成对审阅表：

```bash
python3 scripts/benchmark.py compare \
  results/baseline.json results/anywork.json \
  --output results/review.md
```

`--command` 会被解析成直接进程参数，不经过 shell 执行。可用 `--skill`、`--difficulty` 或 `--case` 选择稳定子集。除非显式传入 `--force`，否则不会覆盖现有结果。

## 证据规则

- 保持模型、版本、温度、工具、账号、网络权限和超时一致。
- AnyWork 组不接收隐藏的 Skill 名称或路由提示；安装的能力环境是唯一预期变量。
- 保留两份 JSON 运行产物和填写完成的审阅表。
- 对任务正确性、可追溯性、直接可用性、母语表达和安全性分别打 0–5 分，并写明证据。
- 在合格审阅者完成评分且排除项已记录前，不得声称胜率。

仓库的 75 个静态案例可生成 225 个中、英、日提示。它们定义测试契约，但不等于模型已经通过。
