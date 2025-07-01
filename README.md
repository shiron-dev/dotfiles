# Brew Manager

Homebrewパッケージ管理のための統合ツール（Go言語版）。YAML設定ファイルを使用してパッケージのインストール、同期、管理を行います。

## 概要

このツールはHomebrewパッケージ管理を効率化するためのGolangで作成されたCLIツールです。以下の機能を提供します：

- **グループ・タグベース管理**: パッケージをグループやタグで分類して管理
- **プロファイル機能**: 事前定義された設定でインストール
- **同期機能**: インストール済みパッケージとYAMLファイルの同期
- **Brewfile変換**: 既存のBrewfileをYAML形式に変換
- **YAML検証**: 設定ファイルの構文チェック

## インストール

```bash
# リポジトリをクローン
git clone <repository-url>
cd dotfiles-feat-ansible/scripts/brew-management

# ビルド
go build -o brew-manager

# インストール（オプション）
go install
```

## 基本的な使用方法

```bash
# ヘルプを表示
./brew-manager --help

# すべてのパッケージをインストール
./brew-manager install

# 特定のグループのみインストール
./brew-manager install --groups core,development

# プロファイルを使用してインストール
./brew-manager install --profile developer

# 現在のインストール状況を同期
./brew-manager sync --auto-detect

# YAML設定ファイルを検証
./brew-manager validate
```

## コマンド一覧

### install
グループ・タグ機能付きYAML設定からパッケージをインストール

```bash
# 基本的なインストール
./brew-manager install

# グループでフィルタリング
./brew-manager install --groups core,development

# タグでフィルタリング
./brew-manager install --tags essential,productivity

# プロファイルを使用
./brew-manager install --profile developer

# 特定のタイプのみインストール
./brew-manager install --groups core --casks-only

# ドライランモード
./brew-manager install --dry-run --groups core

# 利用可能なグループ/タグ/プロファイルを表示
./brew-manager install --list-groups
./brew-manager install --list-tags
./brew-manager install --list-profiles
```

### install-simple
シンプルなYAML設定からパッケージをインストール

```bash
# シンプル形式のYAMLからインストール
./brew-manager install-simple

# 特定のタイプのみ
./brew-manager install-simple --brews-only
```

### sync
インストール済みパッケージをグループ・タグ機能付きYAMLに同期

```bash
# 基本的な同期
./brew-manager sync

# 自動グループ・タグ検出
./brew-manager sync --auto-detect

# インタラクティブモード
./brew-manager sync --interactive

# ソート機能付き
./brew-manager sync --sort

# バックアップ作成
./brew-manager sync --backup

# ドライランモード
./brew-manager sync --dry-run
```

### sync-simple
インストール済みパッケージをシンプルなYAMLに同期

```bash
# シンプル形式で同期
./brew-manager sync-simple

# バックアップ作成
./brew-manager sync-simple --backup
```

### validate
YAML設定ファイルをスキーマで検証

```bash
# すべてのYAMLファイルを検証
./brew-manager validate --all

# 特定のファイルを検証
./brew-manager validate packages.yml

# 詳細出力で検証
./brew-manager validate --verbose --all
```

### convert
BrewfileをYAML形式に変換

```bash
# シンプル形式に変換
./brew-manager convert Brewfile packages.yml

# グループ形式に変換（自動検出付き）
./brew-manager convert --grouped Brewfile packages-grouped.yml
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

## 設定オプション

### グローバルオプション
- `--verbose, -v`: 詳細な出力を有効化
- `--dry-run, -d`: 実際の処理を行わずに表示のみ

### インストールオプション
- `--groups, -g`: 指定したグループのみインストール
- `--tags, -t`: 指定したタグを持つパッケージのみインストール
- `--exclude-groups`: 指定したグループを除外
- `--exclude-tags`: 指定したタグを持つパッケージを除外
- `--profile, -p`: プロファイルを使用
- `--taps-only`: tapのみインストール
- `--brews-only`: brew formulaeのみインストール
- `--casks-only`: caskのみインストール
- `--mas-only`: Mac App Store appsのみインストール

### 同期オプション
- `--backup, -b`: 変更前にバックアップを作成
- `--sort, -s`: パッケージをアルファベット順にソート
- `--show-only`: 変更せずに不足パッケージのみ表示
- `--default-group`: 新規パッケージのデフォルトグループ
- `--default-tags`: 新規パッケージのデフォルトタグ
- `--interactive, -i`: 各パッケージのグループ/タグを対話的に設定
- `--auto-detect, -a`: パッケージ名から自動的にグループ/タグを検出

## 使用例

### 開発環境のセットアップ

```bash
# 開発者プロファイルでインストール
./brew-manager install --profile developer

# または段階的に
./brew-manager install --groups core
./brew-manager install --groups development --exclude-tags experimental
```

### パッケージ管理

```bash
# 新しくインストールしたパッケージを同期
./brew-manager sync --auto-detect

# 未分類パッケージをインタラクティブに分類
./brew-manager sync --interactive
```

### 既存環境の移行

```bash
# 1. 既存のBrewfileを変換
./brew-manager convert Brewfile packages.yml

# 2. 現在の環境を同期
./brew-manager sync-simple --backup

# 3. グループ・タグ機能を使用開始
./brew-manager convert --grouped Brewfile packages-grouped.yml
./brew-manager sync --auto-detect
```

## 開発

### ビルド

```bash
go build -o brew-manager
```

### テスト

```bash
go test ./...
```

### パッケージ構造

```
brew-manager/
├── cmd/            # CLI コマンド定義
├── pkg/
│   ├── types/      # 型定義
│   ├── utils/      # ユーティリティ関数
│   ├── yaml/       # YAML操作
│   ├── brew/       # Homebrew操作
│   ├── validate/   # YAML検証
│   ├── sync/       # 同期機能
│   └── convert/    # 変換機能
├── go.mod
├── go.sum
├── main.go
└── README.md
```

## ライセンス

MIT License
