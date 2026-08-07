# AnyWork

**あらゆる仕事、あらゆるエージェント、あらゆるマシン。**

[简体中文](README.zh-CN.md) · [English](README.md) · **日本語**

AnyWork は、AI エージェントのためのオープンソースな仕事能力レイヤーです。検証済みのワークフロー、参考資料、スクリプト、テンプレート、評価をまとめ、新しいマシンでも数分で高品質な AI ワークスペースを再現できるようにします。

AnyWork はプロンプトの寄せ集めではありません。公式 Skill は実際のタスクに集中し、多言語対応、評価、出典追跡、安全なインストールを必須とします。

## プロダクトの約束

- **あらゆる仕事：**汎用スキルに加え、組み合わせ可能な職種別・業界別パックを提供します。
- **あらゆるエージェント：**標準となる一つの Skill ソースを保ち、まず Codex と Claude Code に対応します。
- **あらゆるマシン：**再現可能なインストール、更新、診断、ロールバック、アンインストールを提供します。
- **3 言語を同等に扱う：**ドキュメント、トリガー、例、CLI メッセージ、評価を英語・簡体字中国語・日本語で提供します。
- **信頼性を既定にする：**出典、ライセンス、権限、リスク情報を必須にします。

## クイックスタート

```bash
python3 anywork-cli.py list --lang ja
python3 anywork-cli.py install essential --agent all --dry-run --lang ja
python3 anywork-cli.py install essential --agent all --lang ja
python3 anywork-cli.py doctor --lang ja
```

CLI は `--lang en` と `--lang zh-CN` にも対応します。指定しない場合は OS のロケールを使用します。
Windows では `py anywork-cli.py ...` を使用してください。リポジトリ CLI に第三者ランタイム依存はありません。

## リポジトリ構成

```text
skills/       単一タスクに集中した標準 Skills
packs/        組み合わせ可能な汎用・職種別・業界別パック
registry/     出典、ライセンス、言語、リスク、バージョン情報
adapters/     エージェントごとのインストール規則
evals/        多言語の振る舞い・品質評価
src/          クロスプラットフォーム AnyWork CLI
schemas/      機械可読な契約
tests/        インストーラーとレジストリの回帰テスト
```

## 現在の状態

AnyWork は基盤構築段階です。最初の対応先は Codex と Claude Code です。公開品質基準と評価スイートは後付けの文書ではなく、プロダクトそのものです。

## ライセンス

Apache-2.0。第三者カタログ項目は元のライセンスと帰属表示を保持し、互換性がないソースはコピーせずリンクのみを掲載します。
