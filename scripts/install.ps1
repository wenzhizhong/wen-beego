# install.ps1 - WenBeego 安装脚本 (最终版：保留缩进、独立密钥、UTF8)
$ErrorActionPreference = "Stop"

# 切换到项目根目录
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$ProjectRoot = Split-Path -Parent $ScriptDir
Set-Location $ProjectRoot

function Write-Info { Write-Host "[INFO] $args" -ForegroundColor Green }
function Write-Warn { Write-Host "[WARN] $args" -ForegroundColor Yellow }
function Write-ErrorMsg { Write-Host "[ERROR] $args" -ForegroundColor Red; exit 1 }

$ConfigFile = "conf\app.yaml"
$BackupFile = "conf\app.yaml.bak"

# ----------------------------- 1. 检查 Go -----------------------------
Write-Info "检测 Go 环境..."
if (-not (Get-Command go -ErrorAction SilentlyContinue)) {
    Write-ErrorMsg "Go 未安装，请先安装 Go 1.16+"
}
$goVer = & go version
Write-Info "Go 版本: $goVer"

# ----------------------------- 2. 安装 bee -----------------------------
Write-Info "安装 bee 工具..."
if (-not (Get-Command bee -ErrorAction SilentlyContinue)) {
    & go install github.com/beego/bee/v2@latest
    $goPath = & go env GOPATH
    $env:Path += ";$goPath\bin"
}
Write-Info "bee 已就绪"

# ----------------------------- 3. 下载依赖 -----------------------------
Write-Info "下载 Go 模块依赖..."
& go env -w GOPROXY=https://goproxy.cn,direct
& go mod download
& go mod tidy

# ----------------------------- 4. 备份配置 -----------------------------
if (-not (Test-Path $ConfigFile)) {
    Write-ErrorMsg "配置文件 $ConfigFile 不存在"
}
Copy-Item $ConfigFile $BackupFile -Force
Write-Info "已备份配置至 $BackupFile"

# ----------------------------- 5. 生成四个不同的随机密钥 -----------------------------
function Generate-RandomKey {
    -join ((48..57) + (65..90) + (97..122) | Get-Random -Count 32 | ForEach-Object { [char]$_ })
}

$newKeyCommon = Generate-RandomKey
$newKeyAdminPlat = Generate-RandomKey
$newKeyAdminMchnt = Generate-RandomKey
$newKeyApi = Generate-RandomKey

Write-Info "自动生成随机加密密钥 (branca)..."

$lines = Get-Content $ConfigFile -Encoding UTF8
$outLines = @()
$inBranca = $false
$inCommon = $false
$inAdminPlat = $false
$inAdminMchnt = $false
$inApi = $false

foreach ($line in $lines) {
    # 进入 branca 块
    if ($line -match '^\s*branca\s*:') {
        $inBranca = $true
        $outLines += $line
        continue
    }
    if ($inBranca) {
        # 进入 common 子块
        if ($line -match '^\s*common\s*:') {
            $inCommon = $true
            $outLines += $line
            continue
        }
        # 在 common 块内找 key 行
        if ($inCommon -and $line -match '^(\s*)key\s*:') {
            $indent = $matches[1]   # 保留原缩进（一般是两个空格）
            $outLines += "${indent}key: $newKeyCommon"
            $inCommon = $false      # 退出 common 子块
            continue
        }
        # 进入 admin_plat 子块
        if ($line -match '^\s*admin_plat\s*:') {
            $inAdminPlat = $true
            $outLines += $line
            continue
        }
        if ($inAdminPlat -and $line -match '^(\s*)key\s*:') {
            $indent = $matches[1]
            $outLines += "${indent}key: $newKeyAdminPlat"
            $inAdminPlat = $false
            continue
        }
        # 进入 admin_mchnt 子块
        if ($line -match '^\s*admin_mchnt\s*:') {
            $inAdminMchnt = $true
            $outLines += $line
            continue
        }
        if ($inAdminMchnt -and $line -match '^(\s*)key\s*:') {
            $indent = $matches[1]
            $outLines += "${indent}key: $newKeyAdminMchnt"
            $inAdminMchnt = $false
            continue
        }
        # 进入 api 子块
        if ($line -match '^\s*api\s*:') {
            $inApi = $true
            $outLines += $line
            continue
        }
        if ($inApi -and $line -match '^(\s*)key\s*:') {
            $indent = $matches[1]
            $outLines += "${indent}key: $newKeyApi"
            $inApi = $false
            continue
        }
    }
    $outLines += $line
}

# 写回文件，使用绝对路径和 UTF8 without BOM
$fullPath = (Resolve-Path $ConfigFile).Path
[System.IO.File]::WriteAllLines($fullPath, $outLines, [System.Text.UTF8Encoding]::new($false))
Write-Info "已更新 branca 密钥（四个独立随机值，保留原缩进）"

# ----------------------------- 6. 读取 PostgreSQL 配置 -----------------------------
function Get-PostgresConfig {
    $lines = Get-Content $ConfigFile -Encoding UTF8
    $inPgsql = $false
    $pgHostName = $null
    $pgPortNum = $null
    $pgUserName = $null
    $pgPassWord = $null
    $pgDbName = $null
    foreach ($line in $lines) {
        if ($line -match '^\s*pgsql\s*:') { $inPgsql = $true; continue }
        if ($inPgsql -and $line -match '^\s*host\s*:\s*(.*)') { $pgHostName = $matches[1].Trim('"') }
        if ($inPgsql -and $line -match '^\s*port\s*:\s*(.*)') { $pgPortNum = $matches[1].Trim('"') }
        if ($inPgsql -and $line -match '^\s*user\s*:\s*(.*)') { $pgUserName = $matches[1].Trim('"') }
        if ($inPgsql -and $line -match '^\s*password\s*:\s*(.*)') { $pgPassWord = $matches[1].Trim('"') }
        if ($inPgsql -and $line -match '^\s*dbname\s*:\s*(.*)') { $pgDbName = $matches[1].Trim('"') }
        if ($inPgsql -and $line -match '^\S' -and $line -notmatch '^\s*pgsql') { break }
    }
    return @{
        host = $pgHostName
        port = $pgPortNum
        user = $pgUserName
        password = $pgPassWord
        dbname = $pgDbName
    }
}

function Test-Port {
    param($targetHost, $targetPort, $timeout=2)
    try {
        $tcp = New-Object System.Net.Sockets.TcpClient
        $async = $tcp.BeginConnect($targetHost, $targetPort, $null, $null)
        if ($async.AsyncWaitHandle.WaitOne($timeout*1000)) {
            $tcp.EndConnect($async)
            $tcp.Close()
            return $true
        }
        $tcp.Close()
        return $false
    } catch { return $false }
}

$pgConfig = Get-PostgresConfig
Write-Info "检测 PostgreSQL 连接..."
if ($pgConfig.host -and $pgConfig.port) {
    if (Test-Port $pgConfig.host $pgConfig.port) {
        Write-Info "[OK] PostgreSQL 端口 $($pgConfig.port) 可达"
    } else {
        Write-Warn "PostgreSQL 端口 $($pgConfig.port) 不可达"
    }
} else {
    Write-Warn "无法解析 PostgreSQL 配置"
}

# ----------------------------- 7. 读取 Redis 配置 -----------------------------
function Get-RedisConfig {
    $lines = Get-Content $ConfigFile -Encoding UTF8
    $inRedis = $false
    $redisHostName = $null
    $redisPortNum = $null
    $redisPassWord = $null
    foreach ($line in $lines) {
        if ($line -match '^\s*redis\s*:') { $inRedis = $true; continue }
        if ($inRedis -and $line -match '^\s*host\s*:\s*(.*)') { $redisHostName = $matches[1].Trim('"') }
        if ($inRedis -and $line -match '^\s*port\s*:\s*(.*)') { $redisPortNum = $matches[1].Trim('"') }
        if ($inRedis -and $line -match '^\s*password\s*:\s*(.*)') { $redisPassWord = $matches[1].Trim('"') }
        if ($inRedis -and $line -match '^\S' -and $line -notmatch '^\s*redis') { break }
    }
    return @{
        host = $redisHostName
        port = $redisPortNum
        password = $redisPassWord
    }
}

$redisConfig = Get-RedisConfig
Write-Info "检测 Redis 连接..."
if ($redisConfig.host -and $redisConfig.port) {
    if (Test-Port $redisConfig.host $redisConfig.port) {
        Write-Info "[OK] Redis 端口 $($redisConfig.port) 可达"
    } else {
        Write-Warn "Redis 端口 $($redisConfig.port) 不可达"
    }
} else {
    Write-Warn "无法解析 Redis 配置"
}

# ----------------------------- 8. 读取 RabbitMQ 配置 -----------------------------
function Get-RabbitMQConfig {
    $lines = Get-Content $ConfigFile -Encoding UTF8
    $inQueue = $false
    $inRabbit = $false
    $rmqHostName = $null
    $rmqPortNum = $null
    $rmqUserName = $null
    $rmqPassWord = $null
    foreach ($line in $lines) {
        if ($line -match '^\s*queue\s*:') { $inQueue = $true; continue }
        if ($inQueue -and $line -match '^\s*rabbitmq\s*:') { $inRabbit = $true; continue }
        if ($inRabbit -and $line -match '^\s*host\s*:\s*(.*)') { $rmqHostName = $matches[1].Trim('"') }
        if ($inRabbit -and $line -match '^\s*port\s*:\s*(.*)') { $rmqPortNum = $matches[1].Trim('"') }
        if ($inRabbit -and $line -match '^\s*user\s*:\s*(.*)') { $rmqUserName = $matches[1].Trim('"') }
        if ($inRabbit -and $line -match '^\s*password\s*:\s*(.*)') { $rmqPassWord = $matches[1].Trim('"') }
        if ($inRabbit -and $line -match '^\S' -and $line -notmatch '^\s*rabbitmq') { break }
    }
    return @{
        host = $rmqHostName
        port = $rmqPortNum
        user = $rmqUserName
        password = $rmqPassWord
    }
}

$rmqConfig = Get-RabbitMQConfig
Write-Info "检测 RabbitMQ 连接..."
if ($rmqConfig.host -and $rmqConfig.port) {
    if (Test-Port $rmqConfig.host $rmqConfig.port) {
        Write-Info "[OK] RabbitMQ 端口 $($rmqConfig.port) 可达"
    } else {
        Write-Warn "RabbitMQ 端口 $($rmqConfig.port) 不可达"
    }
} else {
    Write-Warn "无法解析 RabbitMQ 配置"
}

# ----------------------------- 9. 数据库迁移 -----------------------------
Write-Info "执行数据库迁移，请耐心等待结果..."
if ($pgConfig.host -and $pgConfig.port -and $pgConfig.user -and $pgConfig.dbname) {
    $conn = "host=$($pgConfig.host) port=$($pgConfig.port) user=$($pgConfig.user) password=$($pgConfig.password) dbname=$($pgConfig.dbname) sslmode=disable"
    & bee migrate -driver=postgres -conn="$conn"
    if ($LASTEXITCODE -eq 0) {
        Write-Info "[OK] 数据库迁移成功"
    } else {
        Write-Warn "数据库迁移失败"
    }
} else {
    Write-Warn "数据库配置不完整，跳过迁移"
}

# ----------------------------- 10. 编译测试 -----------------------------
Write-Info "尝试编译项目，请耐心等待结果..."
if (Test-Path "cmd/http/main.go") {
    Push-Location "cmd/http"
    & go build -o "$env:TEMP\wenbeego_test.exe"
    if ($LASTEXITCODE -eq 0) {
        Remove-Item "$env:TEMP\wenbeego_test.exe" -Force -ErrorAction SilentlyContinue
        Write-Info "[OK] 编译成功"
    } else {
        Write-Warn "编译失败，请检查代码"
    }
    Pop-Location
} elseif ((Get-ChildItem -Filter "*.go" | Measure-Object).Count -gt 0) {
    & go build -o "$env:TEMP\wenbeego_test.exe"
    if ($LASTEXITCODE -eq 0) {
        Remove-Item "$env:TEMP\wenbeego_test.exe" -Force -ErrorAction SilentlyContinue
        Write-Info "[OK] 编译成功"
    } else {
        Write-Warn "编译失败，请检查代码"
    }
} else {
    Write-Warn "未找到 Go 源文件，跳过编译"
}

# ----------------------------- 11. 默认凭证警告 -----------------------------
$warnings = @()
if ($pgConfig.password -eq "postgres") { $warnings += "PostgreSQL 密码仍为默认值 'postgres'，建议立即修改" }
if ($pgConfig.port -eq "5432") { $warnings += "PostgreSQL 端口为默认 5432，生产环境建议修改" }
if ($rmqConfig.user -eq "guest" -and $rmqConfig.password -eq "guest") { $warnings += "RabbitMQ 用户名/密码为默认 'guest/guest'，建议立即修改" }
if ([string]::IsNullOrEmpty($redisConfig.password) -or $redisConfig.password -eq '""') { $warnings += "Redis 未设置密码，生产环境建议配置密码" }

if ($warnings.Count -gt 0) {
    Write-Host "`n========== 安全警告 ==========" -ForegroundColor Yellow
    foreach ($w in $warnings) { Write-Host "WARNING: $w" -ForegroundColor Yellow }
    Write-Host "================================`n" -ForegroundColor Yellow
}

Write-Info "安装完成！可以使用 'bee run' 启动项目"