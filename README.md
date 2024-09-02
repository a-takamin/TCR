# TCR
Takamin Container Registry

## 環境構築
[mise](https://mise.jdx.dev/getting-started.html) を利用しているので、以下のコマンドで必要なツールをインストールできます。
```sh
$ mise install
```

永続化層はコスト対策のために DynamoDB を使っているため、ローカル開発では DynamoDB Local を使います。
```sh
$ docker compose up
```
で起動後、`localhost:8001` で DynamoDB Local のコンソールに入れます。