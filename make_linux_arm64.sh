#!/bin/bash

set -e -x
uname -m

# 리눅스 및 arm64 아키텍처 확인
if [[ "$(uname -s)" != "Linux" || "$(uname -m)" != "aarch64" && "$(uname -m)" != "arm64" ]]; then
    echo "Error: This script is only for Linux (ARM64)."
    echo "Detected: $(uname -s) ($(uname -m))"
    exit 1
fi

echo "Linux ARM64 detected. Starting build..."

# 1. Go 1.24 설치 및 경로 설정 (프로젝트 밖인 상위 폴더에 설치하여 간섭 차단)
GO_INSTALL_DIR="$HOME/sdk/go1.24"
if [ ! -d "$GO_INSTALL_DIR" ]; then
    echo "🌐 Go 1.24 설치 중 (위치: $GO_INSTALL_DIR)..."
    mkdir -p "$HOME/sdk"
    curl -OL https://go.dev/dl/go1.24.0.linux-arm64.tar.gz
    tar -xzf go1.24.0.linux-arm64.tar.gz
    mv go "$GO_INSTALL_DIR"
    rm go1.24.0.linux-arm64.tar.gz
fi

# 환경 변수 설정
export GOROOT="$GO_INSTALL_DIR"
export GOPATH="$HOME/go"
export PATH="$GOROOT/bin:$GOPATH/bin:$PATH"

# 2. 최신 Fyne CLI 설치 (안내된 대로 새로운 경로 사용)
if ! command -v fyne &> /dev/null; then
    echo "📦 최신 Fyne CLI 설치 중..."
    go install fyne.io/tools/cmd/fyne@latest
fi

sudo apt-get update
sudo apt-get install -y libgl1-mesa-dev xorg-dev libx11-dev libxcursor-dev libxrandr-dev libxinerama-dev libxi-dev libxxf86vm-dev libxkbcommon-dev libwayland-dev pkg-config dpkg-dev fakeroot

APP_NAME="ScriptLauncher"
APP_ID="com.chunsam.scriptlauncher"
ICON_PATH="app.png"

echo "🚀 멀티 아키텍처 빌드 시작: $APP_NAME..."

# 현재 폴더(. / script-launcher)의 의존성만 정리
go mod tidy

# 3. 빌드 및 패키징
if [ -f "$ICON_PATH" ]; then
    echo "📦 Fyne 패키징 시작..."
    
    # 1. 일단 Fyne으로 기본 패키징 (tar.xz 생성됨)
    CGO_ENABLED=1 GOOS=linux GOARCH=arm64 \
    fyne package -os linux -id "$APP_ID" -icon "$ICON_PATH" -name "${APP_NAME}_arm64"

    # 2. 만약 .tar.xz가 나왔다면 수동으로 .deb 변환
    if [ -f "${APP_NAME}_arm64.tar.xz" ]; then
        echo "📂 .tar.xz 감지됨. .deb 패키지로 변환 중..."
        
        # 임시 작업 폴더 생성
        mkdir -p tmp_deb/DEBIAN
        tar -xf "${APP_NAME}_arm64.tar.xz" -C tmp_deb/
        
        # DEBIAN/control 파일 생성 (패키지 정보)
        cat <<EOF > tmp_deb/DEBIAN/control
Package: scriptlauncher
Version: 1.0.0
Section: utils
Priority: optional
Architecture: arm64
Maintainer: chunsam <chunsam-dotcom@github.com>
Description: GUI Script Launcher
EOF

        # .deb 파일 생성
        fakeroot dpkg-deb --build tmp_deb "${APP_NAME}_arm64.deb"
        
        # 정리
        rm -rf tmp_deb "${APP_NAME}_arm64.tar.xz"
    fi

    echo "✅ 빌드 완료!"
    #ls -l | grep ".deb"
else
    echo "⚠️ $ICON_PATH 파일이 없습니다."
fi