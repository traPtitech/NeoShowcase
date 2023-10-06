# Architecture

Currently only available in Japanese.

## 概要

NeoShowcaseはtraP部員（ユーザー）が利用できるPaaS（Platform as a Service）です。
ここでのPaaSとは、ユーザーが動かしたいアプリケーションのコードと少しの設定を持ち込むだけで、PaaSがアプリケーションをユーザーの代わりに動かすサービスをいいます。
ユーザーはGit上で自身のアプリケーションのコードを管理し、GitHubやGiteaなどのサービスにそれをpushすることで、NeoShowcaseが更新を自動で検知し、ユーザーの設定に従ってビルドを行いアプリケーションを動かします。

NeoShowcaseの大きな設計思想の一つに"reconciliation loop"というものがあります（reconciliation pattern, synchronization, eventual consistencyなどとも）。
これはKubernetesのcontrollerの処理方法に大きく影響を受けています（実際、NeoShowcase自身も一種のcontrollerと呼べるはずです）。
対比されるものに"イベント駆動"（event-driven）な設計があり、これはユーザーのアクションや外部のイベントなどにより内部のプロセスが起動し、何らかの処理を行い、またさらに必要に応じてイベントを発火するものです。
しかし、通信中にイベントが失われたり、イベントによって起動された処理が失敗したりクラッシュしたりすると、それのリトライ処理や失敗の処理の方法が複雑になったりします。
対して"reconciliation loop"では、あるプロセスは一定時間ごとにシステムの状態を監視し、現在の状態が望むものでなかった場合のみ、望む状態になるように処理 "reconciliation" を行います。
イベント駆動に対して少し計算量は多くなる傾向にありますが、複雑なリトライ処理を考える必要はなく、各プロセスは自身が担当する状態のみを理想状態に持っていくことだけを管理すればよくなるため、処理が簡単になります。

各コンポーネントは自身の"reconciliation loop"を持っています。

- controller/repository_fetcher: 3分毎またはイベント発生時に、DBからリポジトリ一覧とアプリ一覧を取得、アプリそれぞれが指定するgit refに対応する最新のcommit hashを取得、DBにその値を保存します。
    - これはArgoCDの更新方法に影響を受けています
    - > Argo CD polls Git repositories every three minutes to detect changes to the manifests. https://argo-cd.readthedocs.io/en/stable/operator-manual/webhook/
- controller/continuous_deployment/build_registerer: 3分毎またはイベント発生時に、まだビルドが行われていない（DBにビルド情報がqueueされていない）アプリをリストアップし、必要なビルド情報をDBに"queued"として保存します。
- controller/continuous_deployment/build_starter: 3分毎またはイベント発生に、接続されているbuilderに次にqueueされたビルドを行うよう指示します。
- controller/continuous_deployment/build_crash_detector: 1分毎に、builderがクラッシュまたは応答しなくなった場合にそれを検知し、該当ビルドを失敗として記録します。
- controller/continuous_deployment/deployments_synchronizer: 3分毎またはイベント発生時に、"起動しているべき"アプリそれぞれに対して最新のビルドが終了したかチェックし、その場合はbackendとssgenにsynchronizationを依頼します。
- controller/backend: dockerまたはk8sのシステムに接続し、実際にコンテナの起動やネットワークの管理を行います。"起動しているべき"アプリの設定一覧を受け取り、実際のシステムの状態がそれと一致するようにコンテナの起動/削除やルーティングを行います。また、ss-genが設定した配信サーバーへルーティングも行います。
- ss-gen: 3分毎またはイベント発生時に、静的サイトのファイルをStorageからダウンロードし、配信可能なように配置します。

以上のように、各コンポーネントは自身が担当する状態のみを監視しそれに専念することで、全体としてはfailureに強いシステムを構築できます。

## コンポーネント

### traefik-forward-auth

https://github.com/traPtitech/traefik-forward-auth

[traefik proxy](https://doc.traefik.io/traefik/) のforward auth middlewareを利用して、ユーザー認証を行います。

基本的に、

- 認証済みなら200 OK
    - softの場合は `/_oauth/login` でログイン可能
- 未認証なら設定されたOAuth/OIDCの認証を行うために307 Temporary Redirect
    - OAuth2リクエストではprompt=noneを最初に試すため、認可画面は（ルートドメインごとに）最初の一度だけしか現れない

のみを行うHTTPサーバーです。
細かい挙動はREADMEを見てください。

### Gateway (ns-gateway)

ユーザーがフロントエンド(dashboard)から直接操作する部分です。
HTTP/1.1上で既存のproxy認証を利用しつつ、型付きの安全な通信を行うため、[Connect · Simple, reliable, interoperable. A better gRPC.](https://connect.build/) を使用しています。

Gatewayというと多数のmicroserviceをまとめるAPI Gatewayがよくありますが、そこまで複雑なAPIでもないため、Gateway自身が全てのAPI操作を担っています。

リクエストを受け取り、ControllerやDB、Storageから各種必要な情報を読取ったり、書き込みます。
Controllerに向けてイベントも発火しますが、このイベントが万が一抜け落ちてもcontroller内部のreconciliation loopによりシステムは自動的に自身の状態を回復します。

### Controller (ns-controller)

NeoShowcaseのコアとなる部分です。
DBに記述された状態を各サブコンポーネントが監視し、望む状態へと持っていき、最終的にアプリのデプロイを行います。

重要なサブコンポーネントの機能は上の記述を参照してください。

### Builder (ns-builder)

Controllerからビルドの指示を受け取り、実際にOCI Image(docker image)のビルドを行います。

現在、ビルド方法は5種類存在します。

- Runtime (buildpack): [Cloud Native Buildpacks · Cloud Native Buildpacks](https://buildpacks.io/) を用いてruntimeアプリのビルドを行います。一般的なアプリであればzero configでビルドすることも可能です。herokuでも使われているやつです。
- Runtime (command): ベースイメージ、ビルド時と起動時のコマンド(シェルスクリプト)をそれぞれ直接記述する方式です。
- Runtime (dockerfile): Dockerfileを指定してビルドする方式です。上２つよりカスタマイズ性が高くなります。
- Static (command): ベースイメージ、ビルド時のコマンド(シェルスクリプト)、ビルド成果物(artifact)が生成されるパスをそれぞれ直接指定し、静的サイトをビルドする方式です。
- Static (dockerfile): Dockerfileを指定して静的サイトをビルドする方式です。上のcommand方式よりカスタマイズ性が高くなります。

それぞれのビルド方法に従ってビルドを行い、生成されたイメージをregistryにpushします。

### Static-Site Generator (ns-ssgen)

静的サイトのビルド成果物ファイルを配置し、配信を行うように設定します。

apache httpd, nginx, caddyなどの静的配信プロセスに対して設定を行うように拡張できます。

### Migrator (ns-migrate)

データベースのマイグレーションを行います。
Goのコードすら無く、[sqldef](https://github.com/k0kubun/sqldef) を実行するスクリプトからのみなります。

マイグレーション時はまず新旧バージョン両方にcompatibleなスキーマを定義し、先にsqldefを実行してスキーマを変更します。
その後、手動もしくはコード内から必要なデータを補完していくことで、ゼロダウンタイムでの移行が可能になります。

もっとも、NeoShowcaseのアプリ自体はcontrollerの介入が無くても動き続けるため、NeoShowcase自身のHigh Availabilityを保証しなくて良い場合はスキーマがbackwards-compatibleなマイグレーションを行う必要は無いです。

## 各種バックエンドとの対応

NeoShowcaseは特定のクラウドに依存しないよう、traefikをベースに設計されています。
各種クラウドのIngress Controllerに対応させることも理論上は可能ですが、実装が多くなって辛いと思います。

|           | Docker(traefik)         | K8s(traefik)                          | K8s(Cloud)     |
|-----------|-------------------------|---------------------------------------|----------------|
| ルーティング    | traefik docker provider | traefik CRD provider                  | Ingress (未実装)  |
| 証明書取得     | traefik Let's encrypt   | traefik Let's encrypt or cert-manager | クラウドによる        |
| 部員認証      | traefik middleware      | traefik middleware                    | クラウドによる        |
| ネットワークの分離 | docker network          | NetworkPolicy                         | クラウドによる        |
| コンテナ      | docker container        | StatefulSet など                        | StatefulSet など |
