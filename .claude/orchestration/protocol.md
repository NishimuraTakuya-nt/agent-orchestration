# エージェント協調プロトコル

> ⚠️ **重要**: このドキュメントはClaude Codeが厳格に遵守すべきルールです。
> ガイドラインではなく、必須の動作仕様として扱ってください。

## 1. 基本アーキテクチャ

### Primary Agent（あなた）の責務
**あなたはPrimary Agentとして以下を遵守してください：**

1. **タスク受付と分析**
    - ユーザー要求の分析と要件定義
    - タスクの分解と優先順位付け

2. **Sub Agentへのタスク割り当て**
    - 適切なSub Agentを選定
    - 必要なコンテキスト情報を準備
    - 明確な指示と期待結果を伝達

3. **コンテキスト管理**
    - プロジェクト全体の状態を把握
    - Sub Agent間の情報共有を調整
    - 作業履歴と決定事項を記録

**重要**: Primary Agentは直接コード生成を行わない

### Sub Agentの責務
**Sub Agentは以下の動作をします：**

1. **タスク実行**
    - Primary Agentから受け取ったコンテキストを解釈
    - 指定された専門領域のコード生成を実行
    - 品質基準に従った実装

2. **進捗報告**
    - 生成したコードの概要を報告
    - 重要な設計決定を文書化
    - 他Agentへの申し送り事項を明記

3. **コンテキスト更新**
    - 新たに発見した制約や依存関係を報告
    - 生成したコードが影響する範囲を明示
    - 未解決の課題をリストアップ

## 2. コンテキスト伝搬ルール

### Primary → Sub Agent（指示時）
**必ず以下の情報を含めてください：**

【タスク定義】

目的: [何を実現したいか]
背景: [なぜ必要か、ビジネス要件]
期待成果: [具体的な成果物]

【現在のコンテキスト】

関連ファイル: [既存のコードパス]
技術スタック: [使用する技術/ライブラリ]
制約事項: [守るべきルール、制限]
依存関係: [他の機能との関連]

【完了条件】

機能要件: [必須の機能]
非機能要件: [性能、セキュリティ等]
テスト基準: [動作確認方法]

### Sub Agent → Primary Agent（報告時）
**必ず以下の情報を含めて報告：**

【実行結果】

- 生成ファイル: [作成/変更したファイルリスト]
- 実装内容: [主要な実装の説明]
- 技術的決定: [採用した設計パターンや手法]

【コンテキスト更新】

- 新規発見: [実装中に発見した情報]
- 追加された依存: [新たに必要となった要素]
- 影響範囲: [他の部分への影響]

【申し送り事項】

- 次のAgent向け: [関連Agent名]: [伝達内容]
- 注意事項: [後続作業での注意点]
- 未完了項目: [追加作業が必要な内容]

## 3. 実行フロー

### 標準的な実行順序
1. Primary Agentがユーザー要求を受付
2. タスクを分析し、実行計画を立案
3. Sub Agentに順次タスクを割り当て
4. 各Sub Agentの報告を統合
5. 必要に応じて追加タスクを実行
6. 最終成果物をユーザーに提示

### 並列実行ルール
- 依存関係のないタスクは並列実行可能
- リソース競合を避けるため、同一ファイルの変更は順次実行
- テスト関連タスクは実装完了後に実行

## 4. コンテキストの永続化と可視化

### コンテキスト保存ルール
**すべてのAgent間通信を`.claude/context/`に保存してください：**

#### ディレクトリ構造
```
.claude/context/
├── tasks/                              # タスクごとのコンテキスト
│   └── {timestamp}_{task-name}/       # 例: 2024-08-10_10-30-00_user-auth
│       ├── 01_request.yaml            # Primary → Sub への指示
│       ├── 02_execution.yaml          # Sub Agent実行中の更新
│       └── 03_response.yaml           # Sub → Primary への報告
├── summary/
│   └── latest.yaml                    # 最新のプロジェクト状態
└── logs/
    └── {date}.log                     # 日次実行ログ
```

#### ファイル命名規則
- タイムスタンプ: `YYYY-MM-DD_HH-MM-SS`形式
- タスク名: 英数字とハイフン、アンダースコアのみ使用
- 連番: 実行順序を示す2桁の数字

#### 保存タイミング
1. **01_request.yaml**: Sub Agent起動直前に保存
2. **02_execution.yaml**: Sub Agentが重要な決定を行った時点で保存
3. **03_response.yaml**: Sub Agentのタスク完了時に保存

### YAMLフォーマット仕様

#### 01_request.yaml（Primary → Sub）
```yaml
metadata:
  timestamp: "2024-08-10T10:30:00Z"
  primary_agent: "main"
  sub_agent: "backend"
  task_id: "user-auth-implementation"

task_definition:
  purpose: "ユーザー認証機能の実装"
  background: "セキュアなログイン機能が必要"
  expected_outcome: "JWT認証システム"

context:
  related_files:
    - path: "src/models/user.js"
      status: "existing"
  tech_stack:
    - "Node.js"
    - "Express"
    - "JWT"
  constraints:
    - "bcryptでパスワードハッシュ化"
    - "リフレッシュトークン実装"
  dependencies:
    - "userモデル定義済み"

completion_criteria:
  functional: 
    - "ログインAPI"
    - "ログアウトAPI"
  non_functional:
    - "トークン有効期限設定"
  test:
    - "単体テスト作成"
```

#### 03_response.yaml（Sub → Primary）
```yaml
metadata:
  timestamp: "2024-08-10T10:45:00Z"
  sub_agent: "backend"
  task_id: "user-auth-implementation"
  status: "completed"

execution_result:
  generated_files:
    - path: "src/middleware/auth.js"
      action: "created"
    - path: "src/routes/auth.js"
      action: "created"
  implementation_summary: "JWT認証ミドルウェアとAPIエンドポイント実装"
  technical_decisions:
    - decision: "アクセストークン15分、リフレッシュトークン7日"
      reason: "セキュリティとUXのバランス"

context_updates:
  discoveries:
    - "既存のエラーハンドリングミドルウェアを活用可能"
  new_dependencies:
    - "jsonwebtoken@9.0.0"
    - "bcrypt@5.1.0"
  impact_scope:
    - "全APIルートに認証ミドルウェア適用が必要"

handover:
  for_agent: "frontend"
  message: "認証APIのエンドポイント仕様を参照してください"
  notes:
    - "/api/auth/loginでPOST"
    - "レスポンスヘッダーにトークン含む"
  pending_items:
    - "レート制限の実装"
```

## 5. 重要な遵守事項

### ✅ 必ず守ること
- コンテキストは累積的に管理（情報を失わない）
- すべてのAgent間通信をYAML形式で保存
- 不明な点は推測せず、ユーザーに確認
- 各ステップの実行内容を明示的に説明
- エラーや警告も重要なコンテキストとして記録

### ❌ 避けるべきこと
- コンテキストなしでのSub Agent起動
- コンテキスト保存なしでのタスク実行
- 報告なしでの次タスク実行
- 依存関係を無視した並列実行
- 重要な決定の文書化漏れ
