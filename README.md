# semantra

## 1. 要求定義

(Systems-Engineering Ready / 100 h スコープ)

項目 定義
プロダクト名 Semantra ― Semantic × Terra（空間）を統合したデジタルツイン基盤
PoC 目的 「意味」と「位置」が結び付いた地物をリアルタイムに登録・検索・可視化できる OSS ベースのクラウドサービスを 100 時間で実証する
適用ドメイン 都市計画／モビリティシミュレーション／衛星体験／3D CG オーサリング

⸻

### 1.1 ミッションとスコープ

ID ミッション (Why) スコープ内 (PoC) スコープ外
M1 空間情報を「意味を持つ構造体」として再構成し、検索・編集・可視化できる基盤を実証する GeoJSON & 3D Tiles の統合表示 / 意味+空間検索 API 生成 AI、自動プロシージャル生成
M2 完全 OSS・GitOps で社内外問わず再現できる DevOps パイプラインを示す Forgejo ＋ Woodpecker ＋ Harbor ＋ k3d 商用 SaaS CI、マルチクラウド切替
M3 Sony RDC_R0119 に対して GIS／クラウド／3D 連携力をポートフォリオで証明 PostGIS ＋ Elastic ＋ MapLibre ＋ Cesium デモ 大規模 LiDAR 処理

⸻

### 1.2 ステークホルダー要求 (Stakeholder Requirements)

SR-ID ステークホルダー 要求
SR-1 デジタルツイン開発者 意味タグと位置で瞬時に絞込検索したい
SR-2 デジタルツイン開発者 3D Tiles／宇宙ビューと地図をシームレスに扱いたい
SR-3 DevOps チーム Git Push → 自動 Build／Deploy が OSS だけで完結してほしい

⸻

### 1.3 ユーザ要求 (User Requirements)

UR-ID 機能要求 (FR)
UR-1 GeoJSON をアップロードし、地図上に即時反映できる
UR-2 キーワード・タグで全文検索し、該当地物をハイライト表示できる
UR-3 表示範囲 (BBOX) を変更すると地物が再クエリ・再描画される
UR-4 クリックした地物の属性をパネルで確認・編集できる
UR-5 2D 地図 ↔ 3D 都市モデルをワンクリックで切替えたい

⸻

### 1.4 システム要求 (System Requirements)

#### 1.4.1 機能要件 (Functional)

FR-ID 内容 関連 UR
FR-01 POST /models/:m/features で地物登録 (GeoJSON) UR-1
FR-02 GET /features?bbox=... で空間検索 (PostGIS) UR-3
FR-03 GET /search?q=... で意味検索 (Elasticsearch) UR-2
FR-04 GET /features/:id で詳細＋メタ情報取得 UR-4
FR-05 GET /tilesets/:id で 3D Tiles URI 提供 UR-5

#### 1.4.2 非機能要件 (Non-Functional)

NFR-ID 指標 目標値
NFR-01 空間/意味検索レスポンス p95 < 1.0 s
NFR-02 CI/CD 所要時間 Build→Deploy < 5 min
NFR-03 再現性 make up 一発起動 100 %

⸻

#### 1.5 制約条件 (Constraints)

• 工数: 100 h
• 完全 OSS: 商用クラウド/SaaS を使用しない
• ローカル完結: k3d 上でフルスタック動作
• データ: 公開 OSM サンプル & ダミー 3D Tiles

⸻

1.6 合格基準 (Acceptance)

テスト 成功条件
API E2E FR-01〜05 全エンドポイントが 200 OK
UX Demo Map 2D/3D 切替・検索・詳細表示が動作
CI/CD Git Push → Woodpecker → Harbor Push → k3d 更新まで自動
ドキュメント README, ER 図, OpenAPI spec がリポジトリに含まれる

🏗 2. システム構成

(Systems-Engineering View / OSS-only / 100 h スコープ)

⸻

2.1 コンテキスト・レイヤビュー

flowchart TD
User["利用者／外部システム"]-->|HTTPS(JSON)|FE["🌐 Frontend\nNext.js + MapLibre/Cesium"]
FE-->|GraphQL/REST|GW["🌊 API Gateway\nEcho (Go)"]
GW-->|Command<br/>(CREATE/UPDATE)|PG[(🗄 PostGIS)]
GW-->|Query<br/>(SEARCH)|ES[(🔎 Elasticsearch)]
subgraph GitOps Pipeline
Forgejo["📝 Forgejo\n(Git)"]-->WP["🔧 Woodpecker CI"]
WP-->Harbor["📦 Harbor"]
Harbor-->|signed image|Argo["🚀 Argo CD"]
end
Argo-->|Kustomize manifests|k3d["☸ k3d Cluster"]
k3d --> PG & ES & GW & FE

レイヤ コンポーネント 主要機能 OSS 根拠
UI Next.js / MapLibre GL / CesiumJS 地図・3D 表示・検索 MIT/Apache
BFF/API Echo (Go) + gqlgen REST・GraphQL、CQRS MIT
Spatial DB PostgreSQL + PostGIS 空間格納・BBOX 検索 PostgreSQL Licence
Semantic DB Elasticsearch OSS 全文・タグ検索 Apache-2.0
CI Woodpecker CI Build/Test/Scan Apache-2.0
Registry Harbor + Cosign 署名付き OCI 配布 Apache-2.0
CD Argo CD + Kustomize GitOps デプロイ Apache-2.0
Runtime k3d→Kubernetes オートスケール・観測性 Apache-2.0

⸻

2.2 物理デプロイ図（PoC）

graph LR
subgraph Dev Laptop
k3d[(k3d Cluster)]
k3d --> apiPod(API)
k3d --> pgPod(PostGIS)
k3d --> esPod(Elastic)
k3d --> uiPod(UI)
end
subgraph DevOps Host
Forgejo-->Woodpecker-->Harbor-->ArgoCD
end
apiPod-->|JDBC/5432|pgPod
apiPod-->|HTTPS/9200|esPod
uiPod-->|GraphQL/REST|apiPod

ローカルで全サービス稼働し、EKS/Fargate 等へそのまま昇格可能なクラスタ設計。

⸻

2.3 コンテナ詳細 (k3d/Compose)

サービス イメージ例 CPU/Memory 初期値 ポート 水平スケール戦略
api semantra/api:tag 0.3 vCPU / 256 Mi 8080 HPA (QPS)
ui semantra/ui:tag 0.2 vCPU / 128 Mi 3000 固定 1
postgis postgis/postgis:15-3.3 0.5 vCPU / 512 Mi 5432 StatefulSet 1
elastic elasticsearch:8.13.4 1 vCPU / 1 Gi 9200 単ノード

⸻

2.4 主要データフロー

✔ 地物登録 (Command Path) 1. UI → /models/:m/features (JSON/GeoJSON) 2. API Gateway → PostGIS (INSERT) 3. 成功時、API→Elasticsearch (Index) 4. API → 200 OK → UI

✔ 検索 (Query Path) 1. UI 検索バー入力 → /search?q=…&bbox=… 2. API → Elasticsearch match + geo_bounding_box 3. API → PostGIS for accurate clip (optional) 4. 結果配列 → UI → 地図ハイライト

⸻

2.5 品質属性対応

属性 施策
性能 GIST idx (PostGIS)、ES single-shard、Go Echo → p95 <1 s
スケール API & UI Deployment に HPA、DB は PoC は単体・将来はパーティション
セキュリティ RBAC 最小権限
再現性 make up → Compose → k3d、manifest は Kustomize でバージョン管理

⸻

2.6 工数アロケーション (100 h)

項目 時間 内容
インフラ設定 (k3d, Harbor, TLS) 12h Base クラスタ＋レジストリ
CI/CD パイプライン 8h Woodpecker + Argo CD
API & DB スキーマ 18h Echo scaffolding ＋ DDL
意味検索連携 (ES) 8h Index 設計＋同期実装
UI 地図＋検索 24h MapLibre 実装＋状態管理
3D Tiles 簡易表示 6h Cesium embed
テスト & ドキュメント 12h E2E + README / OpenAPI
緩衝 12h バッファ + 予備改善

🧩 3. データモデル

(Relational + NoSQL Hybrid / PostGIS × Elasticsearch)

⸻

3.1 データ設計ポリシー

原則 内容
S-S 分離 Spatial-Store ＝ PostGIS（厳密位置・トポロジ）Semantic-Store ＝ Elasticsearch（全文・タグ）
7W 拡張性 7W 属性を必ずどこかにマッピングできる列／フィールドを保持し、後付け属性も JSON で許容
一意識別子 全エンティティを UUID v7（時系列ソート可）で統一
同期保証 地物 INSERT/UPDATE/DELETE 時に DB → ES へ Tx 後非同期 イベントで整合（低結合）

⸻

3.2 リレーショナルモデル（PostGIS）

3.2.1 ER 図（論理）

erDiagram
model {
UUID id PK
TEXT name "例: building, route"
}

feature {
UUID id PK
UUID model_id FK
GEOMETRY geom SRID 4326
JSONB props
TIMESTAMPTZ obs_time
TIMESTAMPTZ created_at DEFAULT now()
TIMESTAMPTZ updated_at
}

tag {
UUID id PK
TEXT name UNIQUE
}

feature_tag {
UUID feature_id FK
UUID tag_id FK
PRIMARY KEY (feature_id, tag_id)
}

model ||--o{ feature : contains
feature ||--o{ feature_tag: "labeled"
tag ||--o{ feature_tag: "tagged"

3.2.2 主要テーブル DDL

CREATE EXTENSION IF NOT EXISTS postgis, "uuid-ossp";

CREATE TABLE model (
id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
name TEXT NOT NULL UNIQUE
);

CREATE TABLE feature (
id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
model_id UUID NOT NULL REFERENCES model(id) ON DELETE CASCADE,
geom GEOMETRY (Geometry, 4326) NOT NULL,
props JSONB,
obs_time TIMESTAMPTZ,
created_at TIMESTAMPTZ DEFAULT now(),
updated_at TIMESTAMPTZ
);

CREATE TABLE tag (
id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
name TEXT NOT NULL UNIQUE
);

CREATE TABLE feature_tag (
feature_id UUID REFERENCES feature(id) ON DELETE CASCADE,
tag_id UUID REFERENCES tag(id) ON DELETE CASCADE,
PRIMARY KEY (feature_id, tag_id)
);

-- 空間 & JSON インデックス
CREATE INDEX idx_feature_geom ON feature USING GIST(geom);
CREATE INDEX idx_feature_props ON feature USING GIN(props jsonb_path_ops);

性能試算（Mac M1・1M 件想定）
BBOX 100km² → 50 ms／タグ検索 (ES) → 30 ms／複合検索 → 80 ms

⸻

3.3 Elasticsearch インデックス

3.3.1 マッピング

PUT places
{
"mappings": {
"properties": {
"id" : { "type": "keyword" },
"model" : { "type": "keyword" },
"name" : { "type": "text", "fields": { "keyword": { "type": "keyword" } } },
"tags" : { "type": "keyword" },
"location" : { "type": "geo_point" },
"props" : { "type": "flattened" },
"obs_time" : { "type": "date" },
"created_at": { "type": "date" }
}
}
}

3.3.2 登録ドキュメント例

{
"id": "018e9c0d-31fa-f3b8-b34b-123456789abc",
"model": "building",
"name": "Shibuya Scramble Square",
"tags": ["landmark", "observation", "skyscraper"],
"location": { "lat": 35.6595, "lon": 139.7005 },
"props": { "height": 229, "floors": 47, "owner": "Tokyu" },
"obs_time": "2025-06-20T08:00:00Z",
"created_at": "2025-06-21T09:12:00Z"
}

⸻

3.4 データ同期戦略

フェーズ 処理
Write Path Echo API → PostGIS INSERT → tx commit → publish geo.event (Redis Stream)
Sync Worker Subscribes geo.event → transforms → bulk index to Elasticsearch
Delete/Update 同上。uuid キー整合で ES ドキュメント置換 or delete

非同期にすることで API レイテンシを最小化、ES を冗長化しても DB トランザクションに影響なし。

⸻

3.5 インデックス & クエリ最適化

ケース 主エンジン 実装例
BBOX 空間検索 PostGIS SELECT \* FROM feature WHERE geom && ST_MakeEnvelope(...)
“park” タグ検索 ES {"query":{"term":{"tags":"park"}}}
複合（タグ＋ BBOX ＋高さ>100） ES → PG ES bool filter → 得られた ids で WHERE id IN (...)

→ レスポンスを保ちつつ、PostGIS で最終クリップする事で精度担保。

⸻

3.6 バージョニング & 監査

オプションだが PoC で準備

CREATE TABLE feature_log (
log_id BIGSERIAL PRIMARY KEY,
feature_id UUID,
op TEXT, -- INSERT/UPDATE/DELETE
before JSONB,
after JSONB,
ts TIMESTAMPTZ DEFAULT now()
);
CREATE TRIGGER trg_feature_audit
AFTER INSERT OR UPDATE OR DELETE ON feature
FOR EACH ROW EXECUTE FUNCTION audit_feature();

🔧 4. バックエンド設計

(Clean Architecture × CQRS × Resource-Oriented API / Go 1.22)

⸻

4.1 アーキテクチャ全景

┌────────────────────────┐
│ handler (HTTP/GraphQL) │ ← Echo + gqlgen
├────────────────────────┤
│ usecase (CQRS) │ ← Command / Query Facade
├────────────────────────┤
│ domain (Entity/VO) │ ← Place, Model, Tag, BBox …  
├────────────────────────┤
│ repository(interface) │ ← PostGISRepo / ESRepo
├────────────────────────┤
│ infra (DB, Search) │ ← sqlx + lib/pq, Elastic v7
└────────────────────────┘

    •	依存逆転: domain は外部に依存しない
    •	CQRS: 書込み(Command) と 読取り(Query) でリポジトリ実装を分離し、

API ハンドラ ↔ ユースケース 間のインターフェースを明瞭化
• トランザクション境界: Command ユースケース内部のみ (pgx.Tx)

⸻

4.2 ディレクトリ構成

/backend
├── cmd/
│ └── main.go エントリ
├── handler/ プレゼンテーション層
│ ├── rest/feature_handler.go
│ └── graphql/
├── usecase/ ユースケース層
│ ├── command/
│ │ └── create_feature.go
│ └── query/
│ └── search_feature.go
├── domain/ エンティティ・VO
│ ├── model.go
│ ├── feature.go
│ └── tag.go
├── repository/ 抽象リポジトリ
│ ├── feature_repo.go
│ └── tag_repo.go
├── infra/ 実装
│ ├── postgres/
│ │ └── feature_repo_pg.go
│ └── elastic/
│ └── feature_repo_es.go
└── pkg/
├── config/ (viper)
├── logger/ (zap)
└── event/ (Redis Streams producer)

⸻

4.3 API デザイン（Resource-Oriented）

4.3.1 Base URL

/api/v1

4.3.2 Endpoints

Verb Path 説明 ユースケース 読み/書き
POST /models/{model}/features 地物登録 (GeoJSON) CreateFeature Command
GET /models/{model}/features/{id} 地物取得 GetFeature Query
GET /models/{model}/features bbox, q, tags パラメータ検索 SearchFeature Query
PATCH /models/{model}/features/{id} 部分更新 UpdateFeature Command
DELETE /models/{model}/features/{id} 削除 DeleteFeature Command
GET /tags タグ一覧 ListTags Query

    •	準拠: REST の名詞ベース / HATEOAS 不採用 (PoC)
    •	拡張: /tilesets/{id} を返すことで Cesium 連携を別リソースに分離

⸻

4.4 ユースケース詳細

UC フロー 例外 レイテンシ目標
CreateFeature 1) JSON → VO 変換 2) PG へ挿入 (Tx)3) Tx commit 後、イベント発火 geo.created 不正 JSON → 400DB Err → 500 150 ms
SearchFeature 1) ES へ Query (タグ/全文)2) 必要なら bbox post-filter3) 結果リスト返却 ES Down → フォールバック無しで 503 300 ms
UpdateFeature 1) PG RowLock → UPDATE2) イベント geo.updated 404 if not found 200 ms

イベント非同期同期: Redis Streams 消費ワーカーが ES を確実に index/update。
失敗時は DLQ Stream に退避し再試行。

⸻

4.5 ドメインオブジェクト (抜粋)

// domain/feature.go
type Feature struct {
ID uuid.UUID
Model string // building, route ...
Geometry geojson.Geometry // domain VO
Props map[string]any // semantic kv
Tags []string
ObsTime time.Time
CreatedAt time.Time
}

// Value Object
type BBox struct {
MinLon, MinLat, MaxLon, MaxLat float64
}

Geometry は独自 VO であり、infra/postgres で ST_GeomFromGeoJSON に変換。

⸻

4.6 非同期イベントスキーマ

// stream: geo.events
{
"ver" : 1,
"type": "geo.created", // or updated, deleted
"id" : "018e...", // UUID
"model": "building",
"payload": {
"geometry": "{...}", // GeoJSON
"tags": ["landmark"]
},
"ts" : "2025-06-22T03:15:11Z"
}

ElasticSync ワーカーは Bulk API で 500 件 / batch

⸻

4.8 テスト戦略

レイヤ フレームワーク カバレッジ目標
単体 (domain) testing + testify 90 %
ユースケース Go test (repo Mock) 80 %
API E2E supertest / newman CRUD 全通

🎨 5. フロントエンド設計

(Next.js 14 App Router × MUI × MapLibre/Cesium / モバイル・アクセシビリティ対応)

⸻

5.1 UI 全体コンセプト

軸 デザイン指針
One-Canvas 1 ページ＝ 1 マップ。Map ⇄ Panel の動的レイアウトで操作コンテキストを途切れさせない
Progressive 3D デフォルト 2D MapLibre → ワンクリックで Cesium 3D に昇格／復帰
Semantic-First 検索バーは「意味（タグ・キーワード）」主導。絞込結果で BBOX が自動ズーム
100 h 測度 Figma レス図・複雑グリッドを排し、MUI + Tailwind utility で高速実装

⸻

5.2 画面フロー

flowchart TD
Launch[起動 /map] --> Idle[Map View 2D]
Idle -->|🔎 Search| Result[検索結果ハイライト]
Result -->|🖱 Click Feature| Detail[サイドパネル詳細]
Detail -->|✏ Edit| EditForm[編集フォーム]
Idle -->|🌐 3D| ThreeD[Cesium View]
ThreeD -->|⤺ 2D| Idle

⸻

5.3 Atomic Design × Feature-Slice

app/
├── (routes)/
│ └── map/ ← 1 ページ
│ ├── page.tsx
│ ├── layout.tsx
│ ├── components/ (Organisms)
│ │ ├── MapView.tsx
│ │ ├── SearchBar.tsx
│ │ ├── SidePanel.tsx
│ │ └── FeatureForm.tsx
│ └── hooks/ (Feature Slice state)
│ ├── useBBox.ts
│ └── useFeatureStore.ts (Zustand)
├── ui/ (Atoms/Molecules – reusable)
│ ├── IconButton.tsx
│ ├── DataChip.tsx
│ └── Modal.tsx
├── lib/
│ ├── apiClient.ts (fetch wrapper)
│ └── maplibreInit.ts
├── types/ (Generated by OpenAPI)
└── styles/ (Tailwind base + MUI theme override)

    •	Feature-slice: map/ 内に状態・UI・hooks が完結。
    •	Zustand: グローバル selectedFeature, features, bbox を一元管理。

⸻

5.4 コア UI コンポーネント

Component 機能 技術要点
MapView MapLibre GL 初期化・BBOX 変更イベント・地物レイヤ描画 useMapLibreGL() custom hook
SearchBar キーワード / タグ chips / Model 下拉 MUI TextField + Autocomplete; debounce 300 ms
SidePanel 詳細表示・タグ編集・GeoJSON ダウンロード Accessible Dialog (role=“dialog”)
FeatureForm GeoJSON ドロップ & マップクリックで座標入力 React Hook Form + Zod validation
ViewToggle 2D/3D 切替スイッチ lazy import('cesium') で初回ロード時 250 KB

⸻

5.5 状態管理 & データフロー

SearchBar onChange
↓ (debounce)
useFeatureStore.search(q,bbox) ── fetch('/search')
↓ (set features[])
MapView ← subscribes → renders Layer
MapView onClick(id) → set selectedFeature
SidePanel ← subscribes → show detail/edit

    •	Zustand subscription で最小再描画
    •	API エラーは react-error-boundary でページ全体を守る

⸻

5.6 パフォーマンス・UX 最適化

対策 詳細
地物タイル分割 >5000 objects は supercluster でクラスタリング
3D ライブラリ遅延 dynamic(() => import('cesium'), { ssr:false })
画像 & アイコン最適化 Next.js Image ＋ svg-sprite-loader
WebGL メモリリーク監視 ResizeObserver で unmount 時 dispose()

⸻

5.7 アクセシビリティ & i18n

項目 実装
キーボード操作 Map 操作以外は Tab ナビ可能、aria-\* 属性付与
カラーブラインド対応 MUI Color Palette を WCAG AA 準拠に
言語切替 next-intl（en/ja 最低対応）

⸻

5.8 テスト

レイヤ ツール カバー
単体 Jest + Testing-Library hooks / utils
UI 視覚 Storybook + Chromatic MapView/Panel snapshots
e2e Playwright 検索 → クリック → 詳細表示フロー

⸻

⚙️ 6. 開発基盤設計

OSS-only • Zero-SaaS

⸻

6.1 アーキテクチャ原則

# 原則 実装

1 Single-Trust Chain Forgejo → Woodpecker → Harbor (Cosign 署名) → Argo CD まで SBOM／署名を継承
2 Repro-by-Make make bootstrap 1 行で Git/CI/Registry/K8s をローカル構築
3 Cluster-in-Repo すべての K8s マニフェストは infra/gitops repo で Git 管理
4 Pull-only CD Argo CD は Harbor イメージタグ と Git の両方を watch。Outbound しない
5 Observability First OpenTelemetry Collector → Prometheus → Grafana を side-stack で常駐

⸻

6.2 コンポーネント & バージョン

レイヤ OSS Version 理由
VCS Forgejo 7.0.x GitHub フル互換 + 軽量
CI Woodpecker CI 2.1.x YAML 驗明 + Kubernetes executor
Registry Harbor 2.11.x Notary v2 / Cosign integration
CD Argo CD 2.11.x GitOps de-facto
Runtime k3d (k8s v1.30) latest ローカル ≒ 本番 API
TLS mkcert + cert-manager 1.13 全内部 FQDN を自動証明
SBOM Syft + Grype latest Build 時に CycloneDX 生成／スキャン
Trace OpenTelemetry Collector 0.93 API latency 可視化
Log Loki 2.9 検索性担保

⸻

6.3 ディレクトリ & Make Targets

repo-root/
├─ .env.example
├─ Makefile # ← 主要ターゲット
├─ docker-compose.dev.yml # UI+API+DB ローカル即起動
├─ infra/
│ ├─ bootstrap/ # Forgejo / Harbor / Woodpecker stacks
│ └─ gitops/ # kustomize base/overlays
└─ apps/
├─ api/ # backend (Go)
│ └─ .woodpecker.yml
└─ ui/ # Next.js
└─ .woodpecker.yml

主な Make

Target 処理 時間目安
make bootstrap Forgejo・Harbor・Woodpecker を 3 分で起動 3 min
make k3d-up k3d + Argo CD + cert-manager 2 min
make dev Compose で UI/API/DB 起動 10 s
make otel OTEL Collector + Grafana 30 s

⸻

6.4 ブランチ戦略 & GitOps フロー

main → prod overlay (タグ push 時)
develop → staging overlay
feature/\* → preview <5h TTL ネームスペース自動生成

Woodpecker がタグ or ブランチ名を読み、image:<git-sha> をビルド →Harbor→
infra/gitops/overlays/<env>/kustomization.yaml の newTag を PR 自動更新。
Argo CD は PR マージ後に環境を同期。

⸻

6.5 Woodpecker CI 高度化 YAML ⬇︎

clone:
depth: 0 # full history for versioning

steps:

- name: sbom-scan
  image: anchore/syft:latest
  commands:

  - syft dir:. -o cyclonedx-json > sbom.json
  - grype dir:. -q --fail-on critical

- name: build-api
  image: docker:24
  privileged: true
  commands:

  - |
    echo "$HARBOR_PASS" | \
      docker login ${HARBOR_URL} -u "$HARBOR_USER" --password-stdin
  - docker build -t ${HARBOR_URL}/semanterra/api:${CI_COMMIT_SHA} .
  - cosign sign --key cosign.key ${HARBOR_URL}/semanterra/api:${CI_COMMIT_SHA}
  - docker push ${HARBOR_URL}/semanterra/api:${CI_COMMIT_SHA}

- name: update-kustomize-tag
  image: ghcr.io/kyoh86/git-ops-updater:v0.4
  when:
  branch: main
  settings:
  repo: git@forgejo.local:infra/gitops.git
  file: overlays/prod/kustomization.yaml
  image: ${HARBOR_URL}/semanterra/api
  newTag: ${CI_COMMIT_SHA}

services:
docker:
image: docker:dind
privileged: true

⸻

6.6 ローカル ⇄ 本番 差異ゼロ方針

項目 ローカル (k3d) 本番 (EKS/Fargate)
Ingress Traefik ALB Ingress
TLS mkcert ACM + cert-manager
Storage local-path EBS CSI Driver
Secrets dotenv → k8s AWS Secrets Manager sync

※Manifests は Kustomize patch で差分 200 行以下に抑制。

⸻

http://{HOST_EXTERNAL_IP}/authorize
