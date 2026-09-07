# エンジニア英語 穴埋めクイズ

`../vocabulary/mercari/` の YAML を読み、日本語例文から英語の空欄に入る語句を選ぶ SvelteKit アプリです。

## 開発

```bash
npm install
npm run dev
```

http://localhost:5173 でクイズを開始できます。

## GitHub Pages への公開

`main` ブランチの `docs/` が GitHub Pages の配信元です。公開用ビルドは JS・CSS・favicon を 1 つの `index.html` にまとめ、リポジトリ直下の `docs/index.html` へコピーします。

公開 URL は https://haruyama480.github.io/learning-engineer-english/ です。事前に `gh auth login` で GitHub CLI へログインしてください。

### ビルドして gh で反映する

リポジトリの `web/` で次を実行します。

```bash
npm install
npm run publish:pages
```

`publish:pages` は次をまとめて行います。

1. `BASE_PATH=/learning-engineer-english` 付きで本番ビルドする
2. 単一ファイルになった `build/index.html` を `../docs/index.html` にコピーする
3. `gh api` で `main` ブランチの `docs/index.html` を作成または更新する

ビルドだけする場合は `npm run build:pages` です。

`gh api` は GitHub 上に直接コミットします。ローカルを合わせるときは、リポジトリルートで次を実行してください。未追跡の `docs/index.html` があると `git pull` が失敗することがあります。

```bash
cd ..
git fetch origin
git checkout origin/main -- docs/index.html
```

### 手動で gh を叩く場合

`npm run build:pages` のあと、リポジトリルートで Contents API を使っても同じ反映ができます。ファイルが大きいので、base64 は引数ではなく JSON ファイル経由で渡します。

```bash
cd ..   # learning-engineer-english/

python3 - <<'PY'
import base64, json, pathlib, subprocess, tempfile, os

html = pathlib.Path('docs/index.html').read_bytes()
payload = {
    'message': 'Publish quiz UI to GitHub Pages',
    'content': base64.b64encode(html).decode(),
    'branch': 'main',
}
try:
    payload['sha'] = subprocess.check_output(
        ['gh', 'api', 'repos/{owner}/{repo}/contents/docs/index.html', '--jq', '.sha'],
        text=True,
    ).strip()
except subprocess.CalledProcessError:
    pass

fd, path = tempfile.mkstemp(suffix='.json')
os.close(fd)
pathlib.Path(path).write_text(json.dumps(payload))
try:
    subprocess.check_call(
        ['gh', 'api', '--method', 'PUT', 'repos/{owner}/{repo}/contents/docs/index.html', '--input', path]
    )
finally:
    os.remove(path)
PY
```

Pages の状態確認は次です。

```bash
gh api repos/{owner}/{repo}/pages --jq '{url: .html_url, status: .status, source: .source}'
```
