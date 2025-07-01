# Brew Management Tools

Homebrewパッケージ管理のための統合ツールセット。YAML設定ファイルを使用してパッケージのインストール、同期、管理を行います。

## 概要

このツールセットは以下の機能を提供します：

- **グループ・タグベース管理**: パッケージをグループやタグで分類して管理
- **プロファイル機能**: 事前定義された設定でインストール
- **同期機能**: インストール済みパッケージとYAMLファイルの同期
- **Brewfile変換**: 既存のBrewfileをYAML形式に変換

## 統合スクリプト

### `brew-manager.sh`

すべての機能を統合したメインスクリプトです。

```bash
# 基本的な使用方法
./brew-manager.sh <command> [options]

# 利用可能なコマンド
./brew-manager.sh --help
```

## Schema Validation

すべてのYAML設定ファイルはJSON Schemaによる検証をサポートしています：

- **packages-grouped.yml**: `schemas/packages-grouped.schema.json`を使用
- **packages.yml**: `schemas/packages-simple.schema.json`を使用

スキーマの機能：
- 構文検証
- 型チェック
- 必須フィールド検証
- 名前・タグのパターン検証
- 対応エディターでの自動補完

### スキーマファイル

- `data/brew/schemas/packages-grouped.schema.json`: グループ機能付き設定用スキーマ
- `data/brew/schemas/packages-simple.schema.json`: シンプル設定用スキーマ

## コマンド一覧

### validate
YAML設定ファイルをJSON Schemaで検証

```bash
# すべてのYAMLファイルを検証
./brew-manager.sh validate

# 特定のファイルを検証
./brew-manager.sh validate packages.yml

# 詳細出力で検証
./brew-manager.sh validate --verbose --all

# 特定のスキーマを指定
./brew-manager.sh validate --schema packages-grouped.schema.json packages-grouped.yml
```

**オプション:**
- `-h, --help`: ヘルプ表示
- `-v, --verbose`: 詳細なエラーメッセージを表示
- `-a, --all`: データディレクトリ内のすべてのYAMLファイルを検証
- `--schema SCHEMA`: 特定のスキーマファイルを使用

### install
グループ・タグ機能付きYAML設定からパッケージをインストール

```bash
# すべてのパッケージをインストール
./brew-manager.sh install

# 特定のグループのみインストール
./brew-manager.sh install --groups core,development

# タグでフィルタリング
./brew-manager.sh install --tags essential,productivity

# プロファイルを使用
./brew-manager.sh install --profile developer

# 特定のタイプのみインストール
./brew-manager.sh install --groups core --casks-only

# ドライランモード（実際にはインストールしない）
./brew-manager.sh install --dry-run --groups core
```

### install-simple
シンプルなYAML設定からパッケージをインストール

```bash
# シンプル形式のYAMLからインストール
./brew-manager.sh install-simple

# 特定のタイプのみ
./brew-manager.sh install-simple --brews-only
```

### sync
インストール済みパッケージをグループ・タグ機能付きYAMLに同期

```bash
# 基本的な同期
./brew-manager.sh sync

# 自動グループ・タグ検出
./brew-manager.sh sync --auto-detect

# インタラクティブモード
./brew-manager.sh sync --interactive

# ソート機能付き
./brew-manager.sh sync --sort
```

### sync-simple
インストール済みパッケージをシンプルなYAMLに同期

```bash
# シンプル形式で同期
./brew-manager.sh sync-simple

# バックアップ作成
./brew-manager.sh sync-simple --backup
```

### convert
BrewfileをYAML形式に変換

```bash
# Brewfileを変換
./brew-manager.sh convert Brewfile packages.yml
```

### list-*
利用可能なグループ、タグ、プロファイルを表示

```bash
# グループ一覧
./brew-manager.sh list-groups

# タグ一覧
./brew-manager.sh list-tags

# プロファイル一覧
./brew-manager.sh list-profiles
```

## YAML設定ファイル形式

### グループ・タグ形式 (`packages-grouped.yml`)

```yaml
metadata:
  version: "2.1"
  supports_groups: true
  supports_tags: true

groups:
  core:
    description: "Essential development tools"
    priority: 1
    packages:
      - name: homebrew/core
        type: tap
        tags: [essential]
      - name: git
        type: brew
        tags: [essential, version-control]
        description: "Distributed version control system"

profiles:
  developer:
    description: "Full development environment"
    groups: [core, development]
    exclude_tags: [experimental]
```

### シンプル形式 (`packages.yml`)

```yaml
taps:
  - homebrew/core
  - homebrew/services

brews:
  - git
  - yq
  - jq

casks:
  - visual-studio-code
  - docker

mas_apps:
  - name: "Xcode"
    id: 497799835
```

## ファイル構成

- `brew-manager.sh` - メイン統合スクリプト
- `install-brew-grouped.sh` - グループ・タグ機能付きインストーラー
- `install-brew-from-yaml.sh` - シンプルインストーラー
- `sync-brew-grouped.sh` - グループ・タグ機能付き同期ツール
- `sync-brew-to-yaml.sh` - シンプル同期ツール
- `convert-brewfile-to-yaml.sh` - Brewfile変換ツール

## 使用例

### 開発環境のセットアップ

```bash
# 開発者プロファイルでインストール
./brew-manager.sh install --profile developer

# または段階的に
./brew-manager.sh install --groups core
./brew-manager.sh install --groups development --exclude-tags experimental
```

### パッケージ管理

```bash
# 新しくインストールしたパッケージを同期
./brew-manager.sh sync --auto-detect

# 未分類パッケージをインタラクティブに分類
./brew-manager.sh sync --interactive
```

### 既存環境の移行

```bash
# 1. 既存のBrewfileを変換
./brew-manager.sh convert Brewfile packages.yml

# 2. 現在の環境を同期
./brew-manager.sh sync-simple --backup

# 3. グループ・タグ機能を使用開始
./brew-manager.sh sync --auto-detect
```

## 必要条件

- **Homebrew**: パッケージ管理システム
- **yq**: YAML処理（自動インストールされます）
- **mas**: Mac App Store アプリ管理（オプション）

## トラブルシューティング

### パスエラー
スクリプトは相対パスで設定ファイルを参照します。正しいディレクトリ構造を維持してください：

```
dotfiles/
├── data/brew/
│   ├── packages.yml
│   └── packages-grouped.yml
└── scripts/brew-management/
    ├── brew-manager.sh
    └── ...
```

### 権限エラー
スクリプトに実行権限を付与してください：

```bash
chmod +x scripts/brew-management/*.sh
``` 
