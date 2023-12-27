#!/bin/bash

# プロジェクトルートに移動
cd $(dirname $0)/..

DATADIR="$(pwd)/.testdb"

if [ -e $DATADIR ]; then
    # すでに mysqld が起動しているとき、再起動する
    PID_FILE=$(ls $DATADIR | grep .pid)
    if [ -n "$PID_FILE" ]; then
        echo "MySQL を再起動します..."
        kill $(cat "$DATADIR/$PID_FILE")
        while [ -e "$DATADIR/$PID_FILE" ]; do
            sleep 1
        done
    fi
fi

# MySQL を初期化
# https://dev.mysql.com/doc/refman/8.0/ja/postinstallation.html
# すでに datadir がある場合失敗するので OK
mysqld --datadir="$DATADIR" --log-error="$DATADIR/mysql.log" --initialize-insecure
# undo_00{1,2} を消しておかないと起動になぜか失敗する
rm $DATADIR/undo_001 $DATADIR/undo_002

# MySQL の起動
mysqld --datadir="$DATADIR" --log-error="$DATADIR/mysql.log" --socket="$DATADIR/mysql.sock" &
MYSQL_PID=$!

# Ctrl-C で MySQL も止まるようにする
function handlestop() {
    echo "MySQL を終了しています..."
    kill $MYSQL_PID
    while [ -e "$DATADIR/mysql.sock" ]; do
        sleep 1
    done
}
trap 'handlestop' 1 2 3 15

echo "socket: $DATADIR/mysql.sock, pid: $MYSQL_PID"
wait $MYSQL_PID
