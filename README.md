# twitter_block2mute

## 使い方

### backend 側の設定

1. `./backend/.env.sample` ファイルを .env ファイルにリネームする。
2. 設定は以下の通り。適宜修正する。

```
# GLOBAL SETTING
RELEASE_MODE="debug" # リリース時は"release"
USE_FRAMEWORK_LOGGER="true" # GinのLoggerを使うかどうか。falseでZapのLoggerのみとなる。

# LOGGER
LOGGER_LEVEL="debug" # ログレベル。 debug / warn / error
LOGGER_OUTPUT_PATHS="stdout" # ログの出力パス。複数指定できる。標準出力は stdout
LOGGER_ERROR_OUTPUT_PATHS="stderr" # エラーログの出力パス。複数指定できる。エラー出力は stderr
FRAMEWORK_LOGGER_OUTPUT_PATHS="stdout" # GinのLoggerのパス

# ROUTING
ROUTING_PORT=":8080" # サーバーのポート
ALLOW_ORIGINS="http://localhost:3000,http://127.0.0.1:3000" # 許可するorigin。フロントエンドのパスを入れる。
ALLOW_HEADERS="Access-Control-Allow-Credentials,Access-Control-Allow-Headers,Content-Type,Content-Length,Accept-Encoding,Authorization,accessToken,X-CSRF-Token,Set-Cookie,Cookie"
ALLOW_METHODS="POST,GET,OPTIONS,PUT,DELETE"
TRUSTED_PROXIES=""
EXPOSE_HEADERS="X-CSRF-Token"
CORS_MAX_AGE="86400" # CORS問い合わせ有効期限
CSRF_SECURE="false" # リリース時はtrueにすること。

# SESSION
SESSION_NAME="mysession" # セッション名
SESSION_SECRET="secret" # セッションシークレット

# DB
DB_TYPE="mysql" # 利用するDBミドルウェア
DB_HOST="localhost" # DBのホスト
DB_PORT="3306" # DBのポート
DB_USERNAME="username" # DBのユーザー名
DB_PASSWORD="password" # DBのパスワード
DB_DBNAME="db_name" # DBの名前

# TWITTER
GOTWI_API_KEY="twitter_api_key" # ツイッターのAPIキー
GOTWI_API_KEY_SECRET="twitter_api_key_secret" # ツイッターのシークレット
TWITTER_CALLBACK_URL="/auth/auth_callback" # 特に変更の必要なし
```

3. （スキップ可能）テストの必要があれば、`./backend/.env.test.sample` ファイルを .env.test ファイルにリネームする。
4. （スキップ可能）設定内容は 2.に同じ。
5. `cd backend`
6. `go run main.go`

### frontend 側の設定

1. `./frontend/env`以下のファイルを下記のようにリネームする

```
env.local.json # ローカル用
env.development.json # development環境用
env.test.json # development環境用
env.production.json # development環境用
```

2. 設定は以下の通り、適宜修正する。

```
{
  "NEXT_PUBLIC_API_URL_BASE": "http://172.28.154.46:8080" // backendのベースURL
}
```

3. 開発用なら以下のコマンドで実行可能。

```
 cd frontend
 npm run dev
```

4. 各ビルドをデプロイする場合は next.js プロダクトのデプロイ方法に従う。

   →[公式](https://nextjs-ja-translation-docs.vercel.app/docs/deployment)

   →[参考](/frontend/README.md)

## 操作

- front 側の URL にアクセスし、画面の指示に従うこと。
