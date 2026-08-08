# AnyWork 導入前後の比較評価

[简体中文](BENCHMARK.zh-CN.md) · [English](BENCHMARK.md) · **日本語**

AnyWork には、改善を思い込みではなく測定するための再現可能な A/B 実行ツールがあります。バイト単位で同一のプロンプトを、クリーンなベースライン・エージェントと AnyWork 導入後の同じエージェントに送り、生の出力と時間を保存し、ペアレビューシートを作成します。自動品質スコアは意図的に作りません。

## 小規模な比較を実行

まず、クリーンなプロファイルまたは Skills ディレクトリでエージェントを実行します。

```bash
python3 scripts/benchmark.py run \
  --mode baseline \
  --command 'your-agent --prompt {prompt}' \
  --language ja \
  --limit 3 \
  --output results/baseline.json
```

次に同じエージェントに AnyWork を導入し、モデルとその他の設定を変えずに同じ選択を実行します。

```bash
python3 scripts/benchmark.py run \
  --mode anywork \
  --command 'your-agent --prompt {prompt}' \
  --language ja \
  --limit 3 \
  --output results/anywork.json
```

ペアレビューシートを生成します。

```bash
python3 scripts/benchmark.py compare \
  results/baseline.json results/anywork.json \
  --output results/review.md
```

`--command` はシェルを介さず、直接のプロセス引数として解析されます。`--skill`、`--difficulty`、`--case` で安定した部分集合を選択できます。`--force` を明示しない限り、既存の出力は上書きしません。

## 証拠ルール

- モデル、バージョン、temperature、ツール、アカウント、ネットワーク権限、タイムアウトを同一に保ちます。
- AnyWork 実行に隠し Skill 名やルーティングヒントは与えず、導入された能力環境だけを意図した差分にします。
- 2 つの JSON 実行成果と完成済みレビューシートを保存します。
- タスク正確性、追跡可能性、直接利用性、母語の自然さ、安全性をそれぞれ 0–5 で評価し、根拠を記録します。
- 適格なレビュアーが評価を完了し、除外を記録するまで勝率を主張しません。

リポジトリの 75 の静的ケースから、英語、簡体字中国語、日本語の 225 プロンプトを生成できます。これらはテスト契約であり、モデルの合格証拠ではありません。
