# Dotfiles Setup

To set up your environment using these dotfiles, run the following command in your terminal:

```bash
bash <(curl -s https://raw.githubusercontent.com/shiron-dev/dotfiles/refs/heads/main/scripts/setup.bash)
```

## Profiles

仕事用 PC と趣味用 PC の差分は `profiles/` で分けています。
セットアップ時にどちらのマシンかを選び、以降はそのプロファイルの設定が
共通設定 (`config/`) の上に重なります。

```bash
./scripts/profile.bash list        # work / private
./scripts/profile.bash set work    # このマシンのプロファイルを決める
./scripts/profile.bash get         # 現在のプロファイル
```

詳細は [profiles/README.md](profiles/README.md) を参照してください。
