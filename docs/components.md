# Components

Currently only available in Japanese.

NeoShowcaseの開発・運用に使われている要素技術の紹介

## Standalone Components

### [Buildkit](https://github.com/moby/buildkit)

Dockerイメージをビルドするツール

Dockerの代わりにDockerイメージをビルドできて、コンテナの中で動かすことができる。
[`buildkitd`](https://github.com/moby/buildkit#containerizing-buildkit)と言うデーモンを起動することでgRPC APIを用いてリモートから操作できるので、NeoShowcaseではこれをアプリビルド用のサーバーとして使う。

Buildkit内部ではLLB(Low Level Build definition format)という中間表現を用いてイメージのビルドを行っており、これを直接操作することも可能。
ただしドキュメントに乏しく、Dockerfile frontendのほうが安定しているため、ほとんどはDockerfile経由でビルドしている。
静的サイトを生成物をイメージから取り出すときなどにLLBを使っている。

### [Buildpacks](https://buildpacks.io/)

Dockerfile無しでアプリケーションのコードからDockerイメージをビルドするツール

アプリケーションが含むファイルやコードから、各"buildpack"が最適なビルド設定を検出し、Dockerイメージを作成する。
自身で"buildpack"を作成することもできる。

NeoShowcaseでは主に、buildpackの実装の一つであるpaketo-buildpacksを使用。

### [gRPC](https://grpc.io/)

Google謹製のHTTP/2を利用したRPCフレームワーク

NeoShowcaseでは、管理サーバー・ビルドサーバーなど目的別サーバー間の通信手段としてgRPCを用いる。

`Protocol Buffer 3`形式(`.proto`)で通信で使う関数の仕様を定義すると、各種言語に対応したサーバー・クライアントのコードを自動生成してくれて、開発者は関数の中身だけ実装すれば良くなる。

例: [Controllerコンポーネントの`proto`ファイル](https://github.com/traPtitech/NeoShowcase/blob/bc9694ca525c1c52530fe2b0358987e64e34900e/api/proto/neoshowcase/protobuf/controller.proto) でns-controllerとns-gatewayの通信を定義して、コード自動生成すると[こんなコード](https://github.com/traPtitech/NeoShowcase/blob/bc9694ca525c1c52530fe2b0358987e64e34900e/pkg/infrastructure/grpc/pb/controller.pb.go)を自動生成してくれて、自分たちは[関数の中身の実装](https://github.com/traPtitech/NeoShowcase/blob/bc9694ca525c1c52530fe2b0358987e64e34900e/pkg/infrastructure/grpc/controller_service.go)をすればいいだけになる。

### [Connect](https://connect.build/)

A better gRPC.

HTTP/1.1, HTTP/2のPOSTメソッドだけを用いて通信を行うプロトコル。
Connectによって生成されたサーバーまたはクライアントのコードはデフォルトでgRPC, gRPC-Web, Connectの3つのプロトコルに対応する。

Connect protocolのWebクライアントはデフォルトでapplication/jsonで通信を行うため、人間が理解しやすく、既存のRESTful APIのエコシステムにも上手くハマる。
NeoShowcaseではこの利点を生かしてtraefik forward auth middlewareに認証を委譲している。

### [traefik proxy](https://doc.traefik.io/traefik/)

モダンなリバースプロキシ

各コンポーネントの接続と、ユーザーのアプリへのルーティングに使用している。
K8s backendでは、Ingress Controllerとして使用。

### [protoc (Protocol Buffer Compiler)](https://github.com/protocolbuffers/protobuf)

`.proto`ファイルから`.go`ファイルなどを生成するときに使うコンパイラ

`make init` で入るようにしてある。
macなら`brew install protobuf`でも入るはず。

`protoc-gen-go` (protocのgoコンパイルプラグイン) が必要。
インストール: `go install google.golang.org/protobuf/cmd/protoc-gen-go`

https://grpc.io/docs/languages/go/quickstart/ も参照

### [evans](https://github.com/ktr0731/evans)

gRPC用クライアント

gRPCは当然Postmanとかcurlとかでアクセス出来ないので、デバッグするときとかには専用クライアントが必要。
これは対話的に呼び出せたりして補完とかも効くので便利。

### [sqldef](https://github.com/sqldef/sqldef)

> The easiest idempotent MySQL/PostgreSQL/SQLite3/SQL Server schema management by SQL.

.sqlファイルにテーブルやindex, foreign keyの定義を書いて、`sqldef schema.sql` すると、.sqlファイルの内容に沿うようにスキーマを変更してくれる。

マイグレーションバージョンの管理が必要なく、冪等で扱いやすい。
新・旧どちらのバージョンにも互換性のあるスキーマを定義し、sqldefでスキーマを更新してから新しいバージョンのデプロイを行うのが普通。
ただし少し凝ったスキーマの変更を行うときは、データの手動マイグレーションが必要になったり、データがうっかり落ちないように注意する必要がある。

### [sql-migrate](https://github.com/rubenv/sql-migrate) (現在不使用 -> sqldef)

DBスキーママイグレーションツール

新たなテーブルを追加したり、既存のテーブルのカラムを追加したりして、開発中にDBのテーブル構造を変えるときに、その変更手順や巻き戻し手順を書いて、DBの構造のバージョン管理をするようにするためのツール。
多分一番シンプル。マイグレーションバージョンの管理、正しいマイグレーションスクリプトの管理を自分で行う必要がある。

NeoShowcaseでは昔sql-migrateを使っていたが、冪等なツールが便利そうだったのでsqldefに移行した。

### [tbls](https://github.com/k1LoW/tbls)

RDBドキュメント自動生成ツール

実際のDBからER図やドキュメントを自動生成したり、DBのLintなどもできる。

例: [こういうの](https://github.com/traPtitech/NeoShowcase/tree/master/docs/dbschema)を自動生成する。

### [golangci-lint](https://github.com/golangci/golangci-lint)

Goコード用のLinter

Linter: コードのフォーマットを指摘してくれたり、危ないコードや不要なコードを検出してくれるツール

### [swagger / OpenAPI 3.0](https://swagger.io/specification/) (現在不使用 -> Connect)

HTTPのAPI仕様を記述するための仕様

traPの内製サービスはほぼ全てこれでAPIの仕様を決定している。
https://apis.trap.jp/

NeoShowcaseでは、昔、Webダッシュボード(管理画面)とサーバー間のAPI仕様を書くのに使っていた。
現在はprotocファイルにかかれていることが全て。

### [spectral](https://github.com/stoplightio/spectral) (現在不使用)

`swagger.yaml`用のLinter

## Go Libraries

### [sqlboiler](https://github.com/volatiletech/sqlboiler)

Go用のSQLDBのORマッパー**ジェネレーター**

他のSysAdプロジェクトでGoからMariaDBにアクセスするときには主に[Gorm](https://gorm.io/)というORMライブラリを使ってますが、NeoShowcaseではデータベースのスキーマからORMライブラリを**生成**するsqlboilerを使います。
NeoShowcase専用のORMライブラリが出来る。

[DBスキーマ](https://github.com/traPtitech/NeoShowcase/tree/master/docs/dbschema)に従ってDBを作成した後、そのDBの構造に従った構造体を[こんな感じ](https://github.com/traPtitech/NeoShowcase/blob/bc9694ca525c1c52530fe2b0358987e64e34900e/pkg/infrastructure/repository/models/applications.go)で自動生成してくれる。

参考:
[Goの新定番？ORMのSQLBoilerを使ってみる | Qiita](https://qiita.com/uhey22e/items/640a4ae861d123b15b53)
[SQLBoiler の使い方を簡単にまとめた | note](https://note.crohaco.net/2020/golang-sqlboiler/)

### [Echo](https://echo.labstack.com/)

SysAd内のデファクトスタンダードなWebサーバーライブラリ

built-inのstatic-site generator内で使用

### [logrus](https://github.com/sirupsen/logrus)

コンソールログをいい感じに出力するようにするやつ

Goの標準logライブラリと互換性があるのでimport文変えるだけで使える。
NeoShowcase自身のログの出力にはこれを使う。

### [cobra](https://github.com/spf13/cobra) / [viper](https://github.com/spf13/viper)

cobraはGoのコマンドラインツール作成支援ライブラリ
viperは設定ファイル取り扱いライブラリ

同じ作者のライブラリで連携している。
NeoShowcaseでは、[`cmd`](https://github.com/traPtitech/NeoShowcase/tree/master/cmd)以下で、実際の実行バイナリのコマンドを作るのに使う。

### [Wire](https://github.com/google/wire)

DI(Dependency Injection)のためのコードを自動生成してくれるライブラリ

参考: https://github.com/google/wire/tree/main/_tutorial

### [docker/client](https://github.com/moby/moby)

GoからDockerを操作するための公式ライブラリ

[`backend/dockerimpl`](https://github.com/traPtitech/NeoShowcase/tree/master/pkg/infrastructure/backend/dockerimpl)中で使ってる。

### [Hub](https://github.com/leandro-lugaresi/hub) (現在不使用)

PubSub型の内部イベントバスライブラリ。

コード内でイベントのPublish / Subscribeができる。
イベントバス使うとコード依存が疎結合になってメンテしやすくなる。

任意のコントロールフローのスパゲッティ化を容易にしてしまうため、使いすぎには注意。
必要ない場合は使わない方が吉かも。
