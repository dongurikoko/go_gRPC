gRPCという存在を知り、勉強してみたかったため、go_lesson1と同様のAPIをgRPCで実装した。
internalはlesson1と同様の構成になっているが、handler層の内容がgRPC仕様になっている。

### gRPCとは?
* https://www.xlsoft.com/jp/blog/blog/2022/05/25/post-29393-post-29393/

### インストール
* protocコマンドをインストール  
  `brew install protobuf`
* GoでgRPCを扱うためのパッケージと、protocコマンドがGoのコードを生成するのに利用するパッケージをインストール

   ```
   go mod init mygrpc
   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
   ```

### protoファイルからコード生成
まず最初にprotoファイルを作るので、スキーマファースト、スキーマ騒動開発と呼ばれる。  
protoファイルからコードを自動生成するには、
```
$ cd api
$ protoc --go_out=../pkg/grpc --go_opt=paths=source_relative \
	--go-grpc_out=../pkg/grpc --go-grpc_opt=paths=source_relative \
	todo.proto
```

インストール時のエラーの対処
- https://qiita.com/pugiemonn/items/67eac48a5254682849db

### 動作確認：  
1.MySQL、phpMyAdmin立ち上げ  
`docker compose up`

2.環境変数設定
```
  export MYSQL_USER=root          
  export MYSQL_PASSWORD=go-lesson2
  export MYSQL_HOST=127.0.0.1
  export MYSQL_PORT=3301
  export MYSQL_DATABASE=go_lesson2_api_gRPC
```

3.サーバーの起動  
`go run cmd/server/main.go`

(4.確認
* サーバー内に実装されているサービス一覧の確認  
`grpcurl -plaintext localhost:8080 list`
 
* あるサービスのメソッド一覧の確認   
例: `grpcurl -plaintext localhost:8080 list todo.TodoService`  
結果:
   ```
   todo.TodoService.Create
   todo.TodoService.Delete
   todo.TodoService.Get
   todo.TodoService.Update
   ```
)

5. メソッドの呼び出し
* Create  : titleを指定して作成する。作成できた場合はそのidを返す。  
コマンド: `grpcurl -plaintext -d '{"title": "test"}' localhost:8080 todo.TodoService.Create`  
レスポンス:  {"id": "1"}

  エラー時の一例：
   ```
   grpcurl -plaintext -d '{"title": ""}' localhost:8080 todo.TodoService.Create 
   ERROR:
     Code: Internal
     Message: Service Unavailable: title is empty
   ```

* Get : クエリがある時はtitleが部分一致しているものを返す、クエリがないときは全部返す  
コマンド : `grpcurl -plaintext -d '{}' localhost:8080 todo.TodoService.Get`   
レスポンス :   
   ```
   {
     "todos": [
       {
         "id": "1",
         "title": "test",
         "createdAt": "2024-01-31T17:22:33Z",
         "updatedAt": "2024-01-31T17:22:33Z"
       },
       {
         "id": "3",
       .......略........
     ]
   }
   ```

   コマンド : `grpcurl -plaintext -d '{"query": "テスト"}' localhost:8080 todo.TodoService.Get`   
   レスポンス :    
   ```
      {
        "todos": [
          {
            "id": "3",
            "title": "テスト",
            "createdAt": "2024-01-31T17:33:10Z",
            "updatedAt": "2024-01-31T17:33:10Z"
          },
          {
             "id": "4",
             "title": "これはテストです！",
	     .......略......
        ]
      }
   ```
* Update : idとtitleを指定して更新する  
コマンド:  
`grpcurl -plaintext -d '{"id": 5, "title":"おやすみ"}' localhost:8080 todo.TodoService.Update`     
レスポンス : {"id": "5" }

* Delete : idを指定して削除する  
コマンド: `grpcurl -plaintext -d '{"id": 2}' localhost:8080 todo.TodoService.Delete`  
レスポンス: {}    
エラー時の一例：
   ```
   grpcurl -plaintext -d '{"id": 10}' localhost:8080 todo.TodoService.Delete
   ERROR:
     Code: NotFound
     Message: Todo not found: failed to delete todo in Delete: sql: no rows in result set
   ```

参考記事：  
gRPCについて：  
https://www.xlsoft.com/jp/blog/blog/2022/05/25/post-29393-post-29393/    
https://zenn.dev/hsaki/books/golang-grpc-starting  
https://github.com/TiraTom/gin-study  
https://qiita.com/asdf22/items/f8609abde6d439c5cfbe  
https://zenn.dev/mrmt/articles/38fb13d9890629  
