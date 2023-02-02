# osemisan-resource-server

Osemisan プロジェクトの OAuth 保護対象リソースサーバーを提供します。

## 使い方

リポジトリをクローンして次のコマンドを実行すると、`http://localhost:9002` で起動します。

```
go run ./main.go
```

## エンドポイント

### GET `/resources`

リクエストにトークン(JWT)が含まれている場合、そのトークンに許可されたスコープのセミに関する情報を返します。

レスポンスは以下の形式の JSON です。

```json
[
  {
    "name": "アブラゼミ",
    "length": "5cm"
  },
  {
    "name": "ミンミンゼミ",
    "length": "3.5cm"
  }
]
```
