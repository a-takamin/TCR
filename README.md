# TCR
Takamin Container Registry

## 環境構築
[mise](https://mise.jdx.dev/getting-started.html) を利用しているので、以下のコマンドで必要なツールをインストールできます。
```sh
$ mise install
```

環境変数を設定します。
※ VSCODE でデバッグモードを利用すると環境変数が読み込まれなかったため、`launch.json` にも直接書いています。
※ そのためデバッグモードのみ使う場合は不要です。
```sh
export AWS_ACCESS_KEY_ID=fake
export AWS_SECRET_ACCESS_KEY=fakefake
```

マニフェストの永続化層はコスト対策のために DynamoDB を使っているため、ローカル開発では DynamoDB Local を使います。
※ LocalStack でもよいのですが、DynamoDB Local を使ってみたかったため。
```sh
$ docker compose up
```
で起動後、`localhost:8001` で DynamoDB Local のコンソールに入れます。

Blob の永続化層には S3 を使っており、ローカルでは S3 互換の `minIO` を使っています。
※ LocalStack でもよいのですが使ってみたかったため。

## ディレクトリ構成
綺麗なディレクトリ構成求む。

```
├─ /internal
│   ├─ /apperrors   # エラー
│   ├─ /client      # 外部 API クライアント（S3 Client など）
│   ├─ /handler     # ハンドラー
│   ├─ /interface   # インターフェース
│   │   └─ /persister # 永続化層に求めるインターフェース 
│   ├─ /model       # ドメインモデル
│   ├─ /repository  # 永続化層
│   └─ /service     # サービス層
│       ├─ /domain    # ドメインロジック
│       ├─ /usecase   # ユースケースの実現
│       └─ /utils     # 便利関数
│
└─ /local-env # ローカル開発環境用
```

### handler
- リクエストのバリデーション
- 適切なユースケースの呼び出し
- 適切なレスポンス（ステータス、ヘッダー、ボディ）の作成

### service/domain
- ドメインロジックの実行

### service/usecase
- ドメインロジックを駆使したビジネスロジックの実行
- レスポンスに必要なデータの返却（エラー、ボディなど）
