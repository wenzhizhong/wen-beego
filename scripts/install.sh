#!/bin/sh
# WenBeego 安装脚本

set -e

# 切换到项目根目录
SCRIPT_DIR=$(cd "$(dirname "$0")" && pwd)
PROJECT_ROOT=$(cd "$SCRIPT_DIR/.." && pwd)
cd "$PROJECT_ROOT"

# 颜色
if [ -t 1 ] && [ "$TERM" != "dumb" ]; then
    GREEN='\033[0;32m'; YELLOW='\033[0;33m'; RED='\033[0;31m'; NC='\033[0m'
else
    GREEN=''; YELLOW=''; RED=''; NC=''
fi

info() { printf "${GREEN}[INFO]${NC} %s\n" "$*"; }
warn() { printf "${YELLOW}[WARN]${NC} %s\n" "$*"; }
error() { printf "${RED}[ERROR]${NC} %s\n" "$*"; exit 1; }

CONFIG="conf/app.yaml"
BACKUP="conf/app.yaml.bak"

# ----------------------------- 环境检查 -----------------------------
info "检测 Go 环境..."
command -v go >/dev/null 2>&1 || error "Go 未安装"
info "Go 版本: $(go version | awk '{print $3}')"

info "安装 bee 工具..."
command -v bee >/dev/null 2>&1 || {
    go install github.com/beego/bee/v2@latest
    GOPATH=$(go env GOPATH)
    export PATH="$PATH:$GOPATH/bin"
}
info "bee 已就绪"

info "下载依赖..."
go env -w GOPROXY=https://goproxy.cn,direct
go mod download && go mod tidy

[ -f "$CONFIG" ] || error "配置文件 $CONFIG 不存在"
cp "$CONFIG" "$BACKUP"
info "已备份配置至 $BACKUP"

# ----------------------------- 生成随机密钥 -----------------------------
gen_key() {
    if command -v openssl >/dev/null 2>&1; then
        openssl rand -base64 32 2>/dev/null | tr -d '\n=' | cut -c1-32
    elif [ -r /dev/urandom ]; then
        cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 32 | head -n1
    else
        echo "$(date +%s%N)$(hostname)" | sha256sum | cut -c1-32
    fi
}
K1=$(gen_key); K2=$(gen_key); K3=$(gen_key); K4=$(gen_key)
info "生成随机密钥并更新 branca 节..."
sed -i.bak -E "/^[[:space:]]*branca:/,/^[^[:space:]]/ {
    /^[[:space:]]*common:/,/^[^[:space:]]/ {
        s/^([[:space:]]*key:).*/\1 $K1/
    }
    /^[[:space:]]*admin_plat:/,/^[^[:space:]]/ {
        s/^([[:space:]]*key:).*/\1 $K2/
    }
    /^[[:space:]]*admin_mchnt:/,/^[^[:space:]]/ {
        s/^([[:space:]]*key:).*/\1 $K3/
    }
    /^[[:space:]]*api:/,/^[^[:space:]]/ {
        s/^([[:space:]]*key:).*/\1 $K4/
    }
}" "$CONFIG"
info "密钥已更新"

# ----------------------------- 提取配置值（清除 \r） -----------------------------
get_val() {
    awk -v s="$1" -v k="$2" '
        $0 ~ "^[[:space:]]*" s ":" { in_s=1; next }
        in_s && $0 ~ "^[[:space:]]*" k ":" {
            sub(/^[[:space:]]*[^:]+:[[:space:]]*/, "")
            gsub(/^"|"$/, "")
            gsub(/\r/, "")
            print
            exit
        }
        in_s && /^[^[:space:]]/ { in_s=0 }
    ' "$CONFIG"
}

PG_HOST=$(get_val "pgsql" "host")
PG_PORT=$(get_val "pgsql" "port")
PG_USER=$(get_val "pgsql" "user")
PG_PASS=$(get_val "pgsql" "password")
PG_DB=$(get_val "pgsql" "dbname")

REDIS_HOST=$(get_val "redis" "host")
REDIS_PORT=$(get_val "redis" "port")
REDIS_PASS=$(get_val "redis" "password")

RABBIT_HOST=$(awk '/^[[:space:]]*queue:/{q=1} q && /^[[:space:]]*rabbitmq:/{r=1;next} r && /^[[:space:]]*host:/{sub(/^[[:space:]]*host:[[:space:]]*/, ""); gsub(/^"|"$/,""); gsub(/\r/,""); print; exit}' "$CONFIG")
RABBIT_PORT=$(awk '/^[[:space:]]*queue:/{q=1} q && /^[[:space:]]*rabbitmq:/{r=1;next} r && /^[[:space:]]*port:/{sub(/^[[:space:]]*port:[[:space:]]*/, ""); gsub(/^"|"$/,""); gsub(/\r/,""); print; exit}' "$CONFIG")
RABBIT_USER=$(awk '/^[[:space:]]*queue:/{q=1} q && /^[[:space:]]*rabbitmq:/{r=1;next} r && /^[[:space:]]*user:/{sub(/^[[:space:]]*user:[[:space:]]*/, ""); gsub(/^"|"$/,""); gsub(/\r/,""); print; exit}' "$CONFIG")
RABBIT_PASS=$(awk '/^[[:space:]]*queue:/{q=1} q && /^[[:space:]]*rabbitmq:/{r=1;next} r && /^[[:space:]]*password:/{sub(/^[[:space:]]*password:[[:space:]]*/, ""); gsub(/^"|"$/,""); gsub(/\r/,""); print; exit}' "$CONFIG")

info "PostgreSQL: ${PG_HOST}:${PG_PORT}"
info "Redis: ${REDIS_HOST}:${REDIS_PORT}"
info "RabbitMQ: ${RABBIT_HOST}:${RABBIT_PORT}"

# ----------------------------- 端口检测函数（优先 /dev/tcp） -----------------------------
test_port() {
    host="$1"; port="$2"; timeout=2
    if command -v bash >/dev/null 2>&1; then
        bash -c "timeout $timeout bash -c 'echo >/dev/tcp/$host/$port' 2>/dev/null" && return 0
        if ! command -v timeout >/dev/null 2>&1; then
            bash -c "echo >/dev/tcp/$host/$port" 2>/dev/null && return 0
        fi
    fi
    if command -v nc >/dev/null 2>&1; then
        nc -z "$host" "$port" 2>/dev/null && return 0
        nc -v -w $timeout "$host" "$port" </dev/null 2>&1 | grep -q "Connected" && return 0
        return 1
    fi
    if command -v telnet >/dev/null 2>&1; then
        echo "quit" | telnet "$host" "$port" 2>/dev/null | grep -q "Connected" && return 0 || return 1
    fi
    warn "未安装 bash/nc/telnet，无法检测端口"
    return 0
}

info "检测 PostgreSQL 连接..."
if [ -n "$PG_HOST" ] && [ -n "$PG_PORT" ]; then
    if test_port "$PG_HOST" "$PG_PORT"; then
        info "[OK] PostgreSQL 端口 $PG_PORT 可达"
    else
        warn "PostgreSQL 端口 $PG_PORT 不可达"
    fi
else
    warn "PostgreSQL 配置不完整"
fi

info "检测 Redis 连接..."
if [ -n "$REDIS_HOST" ] && [ -n "$REDIS_PORT" ]; then
    if test_port "$REDIS_HOST" "$REDIS_PORT"; then
        info "[OK] Redis 端口 $REDIS_PORT 可达"
    else
        warn "Redis 端口 $REDIS_PORT 不可达"
    fi
else
    warn "Redis 配置不完整"
fi

info "检测 RabbitMQ 连接..."
if [ -n "$RABBIT_HOST" ] && [ -n "$RABBIT_PORT" ]; then
    if test_port "$RABBIT_HOST" "$RABBIT_PORT"; then
        info "[OK] RabbitMQ 端口 $RABBIT_PORT 可达"
    else
        warn "RabbitMQ 端口 $RABBIT_PORT 不可达"
    fi
else
    warn "RabbitMQ 配置不完整"
fi

# ----------------------------- 数据库迁移 -----------------------------
info "执行数据库迁移，请耐心等待结果..."
if [ -n "$PG_HOST" ] && [ -n "$PG_PORT" ] && [ -n "$PG_USER" ] && [ -n "$PG_DB" ]; then
    CONN="host=$PG_HOST port=$PG_PORT user=$PG_USER password=$PG_PASS dbname=$PG_DB sslmode=disable"
    bee migrate -driver=postgres -conn="$CONN"
    if [ $? -eq 0 ]; then
        info "[OK] 数据库迁移成功"
    else
        warn "数据库迁移失败"
    fi
else
    warn "数据库配置不完整，跳过迁移"
fi

# ----------------------------- 编译测试 -----------------------------
info "编译测试，请耐心等待结果..."
if [ -f "cmd/http/main.go" ] || ls *.go >/dev/null 2>&1; then
    cd cmd/http
    go build -o /tmp/wenbeego_test
    rm -f /tmp/wenbeego_test
    info "[OK] 编译成功"
else
    warn "未找到 Go 源文件，跳过编译"
fi

# ----------------------------- 安全警告（加强 Redis 检测） -----------------------------
WARN_MSGS=""
[ "$PG_PASS" = "postgres" ] && WARN_MSGS="${WARN_MSGS}\n  ${YELLOW}⚠️${NC} PostgreSQL 密码为默认值 'postgres'，建议立即修改"
[ "$PG_PORT" = "5432" ] && WARN_MSGS="${WARN_MSGS}\n  ${YELLOW}⚠️${NC} PostgreSQL 端口为默认 5432，生产环境建议修改"
[ "$RABBIT_USER" = "guest" ] && [ "$RABBIT_PASS" = "guest" ] && WARN_MSGS="${WARN_MSGS}\n  ${YELLOW}⚠️${NC} RabbitMQ 用户名/密码为默认 'guest/guest'，建议立即修改"

# 清理 REDIS_PASS：去除首尾空格、双引号、单引号
REDIS_PASS_CLEAN=$(printf "%s" "$REDIS_PASS" | sed -e 's/^[[:space:]]*//' -e 's/[[:space:]]*$//' -e 's/^"//' -e 's/"$//' -e "s/^'//" -e "s/'$//")
if [ -z "$REDIS_PASS_CLEAN" ]; then
    WARN_MSGS="${WARN_MSGS}\n  ${YELLOW}⚠️${NC} Redis 未设置密码，生产环境建议配置密码"
fi

if [ -n "$WARN_MSGS" ]; then
    printf "\n${YELLOW}========== 安全警告 ==========${NC}\n"
    printf "%b\n" "$WARN_MSGS"
    printf "${YELLOW}================================${NC}\n\n"
else
    info "未发现明显的默认配置问题"
fi

info "安装完成！运行 'bee run' 启动项目"