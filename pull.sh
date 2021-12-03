#!/bin/sh
# $1 ローカルリポジトリのパス
# $2 リモートの最新コミットID
# $3 対象ブランチ(master固定)

# ローカルリポジトリのパスへ移動
cd $1

# リモートの最新のコミットIDを取得
remote=$(git ls-remote origin $3 | cut -f 1)

# もらった最新と称するコミットIDが最新でなければ終了
if [ $remote != $2 ]; then
    echo "コミットID不一致のため終了"
    exit 1
fi

# ローカルの最新コミット取得
current=$(git rev-parse HEAD)

# 既にローカルが最新であれば終了
if [ $remote = $current ]; then
    echo "コミットIDがすでに最新のため終了"
    exit 0
fi

# 最新取得
git pull
error=$?

if [ $error -eq 1 ]; then
    echo “Error”
    exit 1
fi

# docker build
git pull
error=$?

if [ $error -eq 1 ]; then
    echo "git pullでエラー発生"
    exit 1
fi

docker-compose -f /home/ubuntu/git/backlog-sandbox/docker-compose.yml up -d --build --force-recreate
error=$?

if [ $error -eq 1 ]; then
    echo "docker-composeでエラー発生"
    exit 1
fi

echo "更新完了"
exit 0