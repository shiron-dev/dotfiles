# Profiles

マシンごとの差分（仕事用 PC / 趣味用 PC）を吸収する仕組み。

`config/` は全マシン共通、`profiles/<name>/` はそのプロファイルのマシンだけに適用される
**上書きレイヤ**です。ansible の `config` ロールが `config/` → `profiles/<name>/` の順に
シンボリックリンクを張るため、profiles 側は

- 共通に無いファイルを **追加** する
- 共通と同じパスのファイルを **差し替える**

のどちらもできます。

## アクティブなプロファイル

リポジトリ外のマーカーファイルに書かれた名前が、そのマシンのプロファイルです。

```
~/.config/dotfiles/profile
```

git 管理外なので、マシンの素性をコミットしてしまう事故が起きません。
ファイルが無い場合は `private` にフォールバックします。

```bash
./scripts/profile.bash get          # 現在のプロファイル
./scripts/profile.bash list         # 利用可能なプロファイル
./scripts/profile.bash set work     # 切り替え
./scripts/profile.bash path         # プロファイルディレクトリのパス
```

切り替えたら ansible を流し直してください。

```bash
cd scripts/ansible && ansible-playbook -i hosts.yml site.yml
```

シェルからは `$DOTFILES_PROFILE` で参照できます（`profile.zsh` が export）。

## ディレクトリ構成

`config/` とまったく同じレイアウトです。

```
profiles/<name>/<tool>/.config/<path>   ->  ~/.config/<tool>/<path>
profiles/<name>/<tool>/home/<path>      ->  ~/<path>
profiles/<name>/vscode/Application Support/<file>
                                        ->  ~/Library/Application Support/<Code|Cursor>/User/<file>
```

現在の中身:

| ファイル | 役割 |
| --- | --- |
| `<name>/git/.config/config-local` | `~/.config/git/config` の末尾から `include` される git 設定。共通設定を上書きできる |
| `work/git/.config/aureon-config` | Aureon-inc / KieiAI 配下で使う git identity（`config-local` の `includeIf` から参照） |
| `<name>/zsh/.config/profile.zsh` | `~/.zshrc` の末尾から `source` されるシェル設定。PATH の追加はここ |
| `<name>/mise/.config/conf.d/profile.toml` | mise が `config.toml` にマージする追加ツール定義 |

## 新しい設定をプロファイルに分ける

ツール側に階層化の仕組みがあるなら、それに乗せるのが第一選択です（共通ファイルは
1 つのまま、プロファイルは差分だけを持つ）。

- **git**: `config-local` に書く。`~/.config/git/config` の最後で include しているので後勝ち
- **zsh**: `profile.zsh` に書く。`~/.zshrc` の最後で source しているので PATH も後勝ち
- **mise**: `conf.d/profile.toml` に書く。mise が `config.toml` にマージする

階層化の仕組みが無いツール（VSCode の `settings.json` など）は、
`profiles/<name>/` に同じパスでファイルを置けばファイルごと差し替わります。

## Homebrew

パッケージは `data/brew/packages.yaml` の tag で分けます。

| tag | 意味 |
| --- | --- |
| `work-only` | 仕事用 PC にだけ入れる |
| `private-only` | 趣味用 PC にだけ入れる |
| (tag なし) | 両方に入れる |

```yaml
groups:
    uncategorized:
        packages:
            cask:
                - name: steam
                  tags:
                    - private-only
                - name: cloudflare-warp
                  tags:
                    - work-only
```

インストール時にプロファイルを渡します。

```bash
brew-management install --profile "$(./scripts/profile.bash get)"
brew-management install --profile work --dry-run   # 確認
```

## 新しいプロファイルを足す

`profiles/<name>/` を作れば、それだけで `profile.bash set <name>` の対象になります。
`data/brew/packages.yaml` の `profiles:` にも同名のエントリを足しておくと、
brew も同じ名前で扱えます。
