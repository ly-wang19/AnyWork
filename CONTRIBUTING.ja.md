# AnyWork への貢献

[简体中文](CONTRIBUTING.zh-CN.md) · [English](CONTRIBUTING.md) · **日本語**

一つのタスクに集中し、出典が明確で、安全かつテスト可能な仕事能力を貢献してください。プロンプトの寄せ集めや、互換ライセンスのない公開コンテンツをコピーしないでください。

## Pull Request の前に

1. 一つの明確な仕事成果を選び、小文字 kebab-case の Skill 名を使用する。
2. ワークフローのロジックは一つの `SKILL.md` に保ち、トリガー説明と評価で英語、簡体字中国語、日本語をカバーする。
3. `registry/catalog.json` に出典、ライセンス、バージョン、能力、リスク、副作用、正確な言語レビュー状況を記録する。
4. 改変した第三者コンテンツは、固定された上流 revision、作者、ライセンス、帰属表示、変更点を `THIRD_PARTY.yml` に記録する。
5. 3言語の評価ケースを追加し、再現可能な処理には自動テストを追加する。
6. 次を実行する：

   ```bash
   go test ./...
   go run . doctor --lang ja
   python3 anywork-cli.py doctor --lang ja
   PYTHONPATH=src python3 -m unittest discover -s tests -v
   ```

7. 各 commit に `Signed-off-by: Name <email>` を付け、[DCO](DCO) の Developer Certificate of Origin 1.1 を証明する。

自動翻訳は下書きとして利用できますが、安定版には3言語すべての適格な人手レビューが必要です。AI支援を使った貢献でも、権利、事実、安全性、テストの責任は貢献者が負います。

## 第三者コンテンツの境界

- `original`：Apache-2.0 と DCO に基づくオリジナル貢献。
- `adapted`：互換ライセンスと完全な来歴がある場合のみコピー・改変。
- `curated`：メタデータと上流リンクのみで、内容はコピーしない。

ライセンスなし、source-available、非商用、改変禁止、その他非互換の素材をコア配布物へコピーしてはいけません。
