#!/bin/bash

set -e -x

# 1. Go 1.24 설치 및 경로 설정 (x86_64 버전)
GO_INSTALL_DIR="$HOME/sdk/go1.24"
if [ ! -d "$GO_INSTALL_DIR" ]; then
    echo "🌐 Go 1.24 (amd64) 설치 중..."
    mkdir -p "$HOME/sdk"
    # amd64 전용 Go 바이너리 다운로드
    curl -OL https://go.dev/dl/go1.24.0.linux-amd64.tar.gz
    tar -xzf go1.24.0.linux-amd64.tar.gz
    mv go "$GO_INSTALL_DIR"
    rm go1.24.0.linux-amd64.tar.gz
fi

# 환경 변수 설정
export GOROOT="$GO_INSTALL_DIR"
export GOPATH="$HOME/go"
export PATH="$GOROOT/bin:$GOPATH/bin:$PATH"

# 2. 최신 Fyne CLI 설치
if ! command -v fyne &> /dev/null; then
    echo "📦 최신 Fyne CLI 설치 중..."
    go install fyne.io/tools/cmd/fyne@latest
fi

# x86 환경용 의존성 설치
sudo apt-get update
sudo apt-get install -y libgl1-mesa-dev xorg-dev libx11-dev libxcursor-dev libxrandr-dev libxinerama-dev libxi-dev libxxf86vm-dev libxkbcommon-dev libwayland-dev pkg-config dpkg-dev fakeroot

APP_NAME="ScriptLauncher"
APP_ID="com.chunsam.scriptlauncher"
ICON_PATH="app.png"

echo "🚀 Linux AMD64 빌드 시작: $APP_NAME..."

go mod tidy

# 3. 빌드 및 패키징
if [ -f "$ICON_PATH" ]; then
    echo "📦 Fyne 패키징 시작 (amd64)..."
    
    # GOARCH를 amd64로 설정
    CGO_ENABLED=1 GOOS=linux GOARCH=amd64 \
    fyne package -os linux -id "$APP_ID" -icon "$ICON_PATH" -name "${APP_NAME}_amd64"

    if [ -f "${APP_NAME}_amd64.tar.xz" ]; then
        echo "📂 .deb 패키지로 변환 중..."
        mkdir -p tmp_deb/DEBIAN
        tar -xf "${APP_NAME}_amd64.tar.xz" -C tmp_deb/
        
        # DEBIAN/control 파일 (Architecture: amd64)
        cat <<EOF > tmp_deb/DEBIAN/control
Package: scriptlauncher
Version: 1.0.0
Section: utils
Priority: optional
Architecture: amd64
Maintainer: chunsam <chunsam-dotcom@github.com>
Description: GUI Script Launcher
EOF

        fakeroot dpkg-deb --build tmp_deb "${APP_NAME}_amd64.deb"
        rm -rf tmp_deb "${APP_NAME}_amd64.tar.xz"
    fi

    echo "✅ AMD64 빌드 완료!"
    #ls -l grep "_amd64.deb"
else
    echo "⚠️ $ICON_PATH 파일이 없습니다."
    exit 1
fi