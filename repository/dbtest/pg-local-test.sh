#!/usr/bin/env bash
# 本地真 PostgreSQL 回归:自动拉起 postgres:16 容器 -> 跑 dbtest -> 删容器。
# 用法: ./repository/dbtest/pg-local-test.sh
set -euo pipefail

NAME=ckman-pgtest
PORT=${DBTEST_PG_PORT:-5433}
PASS=${DBTEST_PG_PASS:-Test@1234}
DB=${DBTEST_PG_DB:-ckman_db}
IMAGE=${DBTEST_PG_IMAGE:-postgres:16}

cleanup() { docker rm -f "$NAME" >/dev/null 2>&1 || true; }
trap cleanup EXIT

cleanup
echo ">> 启动 $IMAGE (容器 $NAME, 端口 $PORT)"
docker run -d --name "$NAME" -e POSTGRES_PASSWORD="$PASS" -e POSTGRES_DB="$DB" -p "$PORT":5432 "$IMAGE" >/dev/null

echo -n ">> 等待就绪"
for _ in $(seq 1 30); do
    if docker exec "$NAME" pg_isready -U postgres -d "$DB" >/dev/null 2>&1; then echo " ok"; break; fi
    echo -n .; sleep 1
done

# 仓库根目录(脚本位于 repository/dbtest/ 下)
ROOT=$(cd "$(dirname "$0")/../.." && pwd)
cd "$ROOT"

export DBTEST_PG_HOST=127.0.0.1 DBTEST_PG_PORT="$PORT" \
       DBTEST_PG_USER=postgres DBTEST_PG_PASS="$PASS" DBTEST_PG_DB="$DB"

echo ">> 跑 dbtest (postgres)"
go test -tags dbintegration -v ./repository/dbtest/ -run TestBackends/postgres
