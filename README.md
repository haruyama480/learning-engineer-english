# learning-engineer-english

日本のエンジニアが、英語で業務を進められるようになるための学習リポジトリです。

最初の教材は、メルカリの [Engineer Vocabulary List](https://github.com/mercari/engineer-vocabulary-list) を語彙単位で使える形にしたものです。語彙・例文データと、日本語から英語への穴埋めクイズを YAML で持ち、スキーマは Go の struct で定義しています。

## 構成

```
vocabulary/mercari/
  source/csv/          元の語彙リスト（list_1.csv 〜 list_5.csv）
  source/LICENSE       元リポジトリの CC BY 4.0 ライセンス全文
  vocabulary.yaml      語彙 + 例文（各語 2〜3 文）
  quizzes.yaml         日本語提示・英語穴埋めクイズ
internal/vocabulary/   YAML に対応する Go の型と読み込み処理
cmd/genquiz/           vocabulary.yaml から quizzes.yaml を生成するコマンド
web/                   SvelteKit の選択クイズ UI
```

`vocabulary.yaml` の各エントリは、語彙（日英）と複数の例文を持ちます。例文は元リストをベースにしつつ、文構造を変えてその語の別の側面が見えるようにしています。

`quizzes.yaml` は、例文の英語側から対象語を `____` に置き換えた穴埋め問題です。問題文は日本語、解答は英語です。

## クイズの再生成

例文を直したあとは、次でクイズ YAML を作り直します。

```bash
go run ./cmd/genquiz
```

スキーマとデータの整合はテストで確認できます。

```bash
go test ./...
```

## クイズ UI

日本語の例文を見て、英語の空欄に入る語句を選ぶ Web UI です。リストの絞り込み、シャッフル、語彙ごとの1問出題、ミス復習ができます。

```bash
cd web
npm install
npm run dev
```

ブラウザで http://localhost:5173 を開きます。

## 出典

語彙は Mercari, Inc. の [engineer-vocabulary-list](https://github.com/mercari/engineer-vocabulary-list) を複製・改変したものです。元データは [CC BY 4.0](https://creativecommons.org/licenses/by/4.0/) です。
