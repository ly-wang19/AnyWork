# AnyWork before/after benchmark

[Simplified Chinese](BENCHMARK.zh-CN.md) · **English** · [Japanese](BENCHMARK.ja.md)

AnyWork includes a reproducible A/B runner so improvement can be measured instead of assumed. It sends byte-identical prompts to a clean baseline agent and to the same agent with AnyWork installed, captures the raw outputs and timing, and creates a paired review sheet. It deliberately does not invent an automatic quality score.

## Run a small comparison

First run the agent with a clean profile or Skills directory:

```bash
python3 scripts/benchmark.py run \
  --mode baseline \
  --command 'your-agent --prompt {prompt}' \
  --language en \
  --limit 3 \
  --output results/baseline.json
```

Then install AnyWork for the same agent, keep the model and all other settings unchanged, and run the same selection:

```bash
python3 scripts/benchmark.py run \
  --mode anywork \
  --command 'your-agent --prompt {prompt}' \
  --language en \
  --limit 3 \
  --output results/anywork.json
```

Generate the paired review sheet:

```bash
python3 scripts/benchmark.py compare \
  results/baseline.json results/anywork.json \
  --output results/review.md
```

`--command` is parsed into direct process arguments and is never executed through a shell. Use `--skill`, `--difficulty`, or `--case` to select a stable subset. An existing output is not overwritten unless `--force` is supplied.

## Evidence rules

- Keep model, version, temperature, tools, account, network access, and timeout identical.
- The AnyWork run receives no hidden Skill name or routing hint; the installed capability environment is the only intended difference.
- Preserve both JSON run artifacts and the completed review sheet.
- Score task correctness, traceability, direct usability, native register, and safety from 0–5 with written evidence.
- Do not claim a win rate until a qualified reviewer completes the scorecards and exclusions are documented.

The repository's 75 static cases produce 225 English, Simplified Chinese, and Japanese prompts. They define the test contract; they are not themselves evidence that a model has passed it.
