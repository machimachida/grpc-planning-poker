### grpcの利点

- HTTP2による高速通信
- Protocol BufferなどのIDL(Interface Definition Language)
- 単方向通信・サーバからのストリーミング・クライアントからのストリーミング・双方向ストリーミングをサポート

#### ストリーミング通信を手軽に扱えるのでプランニングポーカーの実装に使えるのでは

---

### HTTP2めちゃ早い

- http://www.http2demo.io/
- テキストベースからバイナリベースな規格になり効率的に
- 一つのTCPコネクションを効率的に利用できるように

---

### IDL(Interface Definition Language)

- リクエスト・レスポンスの形式を記述するための言語
  - Swagger書いてる感覚にちょっと近い
- gRPCでは(主に)Protocol Buffersが使われる
- Protocol Buffersは定義から言語ごとにインターフェースを生成してくれるので、開発が楽になる

---

### 実際に定義ファイルを見てみよう

