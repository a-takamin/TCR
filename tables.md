```mermaid
---
title: DynamoDB Table
---
erDiagram
    Manifest {
        string Name PK "リポジトリ名"
        string Digest PK "(Sort Key)ダイジェスト"
        string Tag "タグ"
        string Manifest "マニフェスト(Base64)"
    }

    ManifestTagLSI {
        string Name PK "リポジトリ名"
        string Tag PK "(Sort Key)タグ"
        string Manifest "マニフェスト(Base64)"
    }

    BlobUpload {
        string Uuid PK "アップロードごとに割り振られる一意のID"
        int ByteUploaded "アップロード済みのバイト数"
        int NextChunkNo "次のチャンク番号"
        boolean Done "すべてのチャンクがアップロードされたかどうか"
        string Digest "ダイジェスト"
    }

    BlobConcat {
      string Digest PK "ダイジェスト"
      string Status "結合処理の状況。notyet, doing, done, error"
    }
```