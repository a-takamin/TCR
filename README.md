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