# YAML-based Homebrew Package Management

BrewfileのようにYAMLファイルからHomebrew packages（tap、brew、cask、MAS）をインストールするためのスクリプト群です。

## 概要

このツールセットには以下のコンポーネントが含まれています：

- `packages.yml` - YAML形式のパッケージ定義ファイル
- `install-brew-from-yaml.sh` - YAMLファイルからパッケージをインストールするスクリプト
- `convert-brewfile-to-yaml.sh` - 既存のBrewfileをYAML形式に変換するスクリプト
- `sync-brew-to-yaml.sh` - 現在インストール済みのパッケージとYAMLファイルを同期するスクリプト
- `packages-grouped.yml` - group/tag機能をサポートする拡張YAML形式のパッケージ定義ファイル
- `install-brew-grouped.sh` - group/tag機能をサポートするインストーラースクリプト
- `sync-brew-grouped.sh` - group/tag機能をサポートする同期スクリプト

## YAML形式（Group/Tag機能）

```yaml
# Group/Tag機能をサポートする拡張形式
metadata:
  version: "2.0"
  supports_groups: true
  supports_tags: true

# グループ定義
groups:
  core:
    description: "Essential development tools"
    priority: 1
  development:
    description: "Development tools and environments"
    priority: 2
  productivity:
    description: "Productivity and office applications"
    priority: 3

# group/tag付きパッケージ定義
brews:
  - name: git
    group: core
    tags: [essential, version-control]
    description: "Distributed version control system"
  - name: docker
    group: development
    tags: [container, development]

# インストールプロファイル
profiles:
  minimal:
    description: "Minimal development setup"
    groups: [core]
    tags: [essential]
  developer:
    description: "Full development environment"
    groups: [core, development]
    exclude_tags: [experimental]

# 技術スタック
stacks:
  web-development:
    description: "Web development stack"
    include:
      brews: [git, node, yarn]
      casks: [visual-studio-code, google-chrome]
      tags: [web, frontend, backend]
```

## 使用方法

### 1. YAMLファイルからパッケージをインストール

```bash
# デフォルトのYAMLファイル（packages.yml）から全てのパッケージをインストール
./scripts/install-brew-from-yaml.sh

# カスタムYAMLファイルから全てのパッケージをインストール
./scripts/install-brew-from-yaml.sh my-packages.yml

# ドライランモード（実際にはインストールしない）
./scripts/install-brew-from-yaml.sh --dry-run

# 詳細出力モード
./scripts/install-brew-from-yaml.sh --verbose
```

### 2. 特定のカテゴリのみインストール

```bash
# tapのみ
./scripts/install-brew-from-yaml.sh --taps-only

# brew formulaeのみ
./scripts/install-brew-from-yaml.sh --brews-only

# caskのみ
./scripts/install-brew-from-yaml.sh --casks-only

# Mac App Storeアプリのみ
./scripts/install-brew-from-yaml.sh --mas-only
```

### 3. 特定のカテゴリをスキップ

```bash
# Mac App Storeアプリをスキップ
./scripts/install-brew-from-yaml.sh --skip-mas

# caskをスキップ
./scripts/install-brew-from-yaml.sh --skip-casks

# tapをスキップ
./scripts/install-brew-from-yaml.sh --skip-taps

# brew formulaeをスキップ
./scripts/install-brew-from-yaml.sh --skip-brews
```

### 4. 既存のBrewfileをYAML形式に変換

```bash
# デフォルトのBrewfileを変換
./scripts/convert-brewfile-to-yaml.sh

# カスタムBrewfileを変換
./scripts/convert-brewfile-to-yaml.sh my-Brewfile output.yml

# 既存ファイルを強制上書き
./scripts/convert-brewfile-to-yaml.sh --force

# 詳細出力
./scripts/convert-brewfile-to-yaml.sh --verbose
```

### 5. 現在インストール済みのパッケージをYAMLファイルに同期

```bash
# 現在インストール済みで未記録のパッケージを表示
./scripts/sync-brew-to-yaml.sh --show-only

# ドライランで何が追加されるかを確認
./scripts/sync-brew-to-yaml.sh --dry-run

# 未分類のパッケージをYAMLファイルに追加（バックアップ付き）
./scripts/sync-brew-to-yaml.sh --backup

# パッケージをアルファベット順でソートして追加
./scripts/sync-brew-to-yaml.sh --sort

# カスタムYAMLファイルと同期
./scripts/sync-brew-to-yaml.sh my-packages.yml
```

### 6. Group/Tag機能を使用したインストール

```bash
# 利用可能なグループを表示
./scripts/install-brew-grouped.sh --list-groups

# 利用可能なタグを表示
./scripts/install-brew-grouped.sh --list-tags

# 利用可能なプロファイルを表示
./scripts/install-brew-grouped.sh --list-profiles

# 特定のグループのみインストール
./scripts/install-brew-grouped.sh --groups core,development

# 特定のタグを持つパッケージのみインストール
./scripts/install-brew-grouped.sh --tags essential,productivity

# 特定のタグを除外してインストール
./scripts/install-brew-grouped.sh --exclude-tags experimental

# プロファイルを使用してインストール
./scripts/install-brew-grouped.sh --profile developer

# 技術スタックを使用してインストール
./scripts/install-brew-grouped.sh --stack web-development

# ドライランで確認
./scripts/install-brew-grouped.sh --groups development --dry-run
```

### 7. Group/Tag機能を使用した同期

```bash
# 自動検出でgroup/tagを割り当てて同期
./scripts/sync-brew-grouped.sh --auto-detect --backup

# インタラクティブでgroup/tagを手動設定
./scripts/sync-brew-grouped.sh --interactive

# デフォルトのgroup/tagを指定して同期
./scripts/sync-brew-grouped.sh --default-group system --default-tags "utility,cli"

# 表示のみ（変更なし）
./scripts/sync-brew-grouped.sh --show-only
```

## オプション一覧

### install-brew-from-yaml.sh

| オプション | 説明 |
|------------|------|
| `-h, --help` | ヘルプメッセージを表示 |
| `-v, --verbose` | 詳細出力を有効化 |
| `-d, --dry-run` | 実際にインストールせずに何がインストールされるかを表示 |
| `-t, --taps-only` | tapのみインストール |
| `-b, --brews-only` | brew formulaeのみインストール |
| `-c, --casks-only` | caskのみインストール |
| `-m, --mas-only` | Mac App Storeアプリのみインストール |
| `--skip-taps` | tapのインストールをスキップ |
| `--skip-brews` | brew formulaeのインストールをスキップ |
| `--skip-casks` | caskのインストールをスキップ |
| `--skip-mas` | Mac App Storeアプリのインストールをスキップ |

### sync-brew-to-yaml.sh

| オプション | 説明 |
|------------|------|
| `-h, --help` | ヘルプメッセージを表示 |
| `-v, --verbose` | 詳細出力を有効化 |
| `-d, --dry-run` | 実際に追加せずに何が追加されるかを表示 |
| `-b, --backup` | YAML ファイルの変更前にバックアップを作成 |
| `-s, --sort` | カテゴリ内でパッケージをアルファベット順にソート |
| `--show-only` | 未分類パッケージを表示するのみ（ファイルを変更しない） |

### install-brew-grouped.sh

| オプション | 説明 |
|------------|------|
| `-h, --help` | ヘルプメッセージを表示 |
| `-v, --verbose` | 詳細出力を有効化 |
| `-d, --dry-run` | 実際にインストールせずに何がインストールされるかを表示 |
| `--list-groups` | 利用可能なグループ一覧を表示 |
| `--list-tags` | 利用可能なタグ一覧を表示 |
| `--list-profiles` | 利用可能なプロファイル一覧を表示 |
| `--list-stacks` | 利用可能なスタック一覧を表示 |
| `-g, --groups GROUPS` | 指定したグループのパッケージのみインストール（カンマ区切り） |
| `-t, --tags TAGS` | 指定したタグを持つパッケージのみインストール（カンマ区切り） |
| `--exclude-groups GROUPS` | 指定したグループを除外（カンマ区切り） |
| `--exclude-tags TAGS` | 指定したタグを持つパッケージを除外（カンマ区切り） |
| `-p, --profile PROFILE` | 事前定義されたプロファイルを使用 |
| `-s, --stack STACK` | 事前定義されたスタックを使用 |

### sync-brew-grouped.sh

| オプション | 説明 |
|------------|------|
| `-h, --help` | ヘルプメッセージを表示 |
| `-v, --verbose` | 詳細出力を有効化 |
| `-d, --dry-run` | 実際に追加せずに何が追加されるかを表示 |
| `-b, --backup` | YAML ファイルの変更前にバックアップを作成 |
| `-s, --sort` | カテゴリ内でパッケージをアルファベット順にソート |
| `--show-only` | 未分類パッケージを表示するのみ（ファイルを変更しない） |
| `--default-group GROUP` | 新しいパッケージに割り当てるデフォルトグループ |
| `--default-tags TAGS` | 新しいパッケージに割り当てるデフォルトタグ（カンマ区切り） |
| `--interactive` | 各パッケージのgroup/tag割り当てを対話的に設定 |
| `--auto-detect` | パッケージ名に基づいてgroup/tagを自動検出 |

### convert-brewfile-to-yaml.sh

| オプション | 説明 |
|------------|------|
| `-h, --help` | ヘルプメッセージを表示 |
| `-v, --verbose` | 詳細出力を有効化 |
| `-f, --force` | 既存の出力ファイルを強制上書き |

## 前提条件

- **Homebrew**: 必須
- **yq**: YAML解析用（自動インストールされます）
- **mas**: Mac App Storeアプリのインストール用（オプション）

## 実行例

```bash
# 1. 既存のBrewfileをYAMLに変換
./scripts/convert-brewfile-to-yaml.sh --force

# 2. ドライランで確認
./scripts/install-brew-from-yaml.sh --dry-run

# 3. 実際にインストール
./scripts/install-brew-from-yaml.sh --verbose

# 4. 特定のカテゴリのみインストール
./scripts/install-brew-from-yaml.sh --brews-only --verbose

# 5. インストール済みパッケージとの同期
./scripts/sync-brew-to-yaml.sh --backup --sort

# 6. Group/Tag機能を使用したインストール
./scripts/install-brew-grouped.sh --profile developer --dry-run

# 7. Group/Tag機能を使用した同期
./scripts/sync-brew-grouped.sh --auto-detect --backup
```

## エラーハンドリング

- インストールに失敗したパッケージは赤色で表示されます
- 個別のパッケージの失敗は全体のプロセスを停止しません
- `--verbose`オプションでより詳細なエラー情報を確認できます

## 利点

1. **柔軟性**: 特定のカテゴリのみインストール可能
2. **安全性**: ドライランモードで事前確認可能
3. **可読性**: YAMLファイルでの管理が直感的
4. **互換性**: 既存のBrewfileから簡単に移行可能
5. **カスタマイズ性**: カテゴリごとのスキップやフィルタリングが可能
6. **同期機能**: インストール済みパッケージを自動で未分類カテゴリに追加
7. **バックアップ機能**: 変更前の自動バックアップで安全性を確保
8. **Group/Tag機能**: パッケージを論理的にグループ化・タグ付けして管理
9. **プロファイル機能**: 用途別の事前定義セットでのインストール
10. **スタック機能**: 技術スタック別でのパッケージ管理
11. **自動検出機能**: パッケージ名に基づくgroup/tagの自動割り当て 
