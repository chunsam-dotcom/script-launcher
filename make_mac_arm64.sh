#!/bin/bash

# 1. OS가 macOS(Darwin)인지 확인
OS_TYPE=$(uname -s)
# 2. 아키텍처가 arm64(Apple Silicon)인지 확인
ARCH_TYPE=$(uname -m)

if [[ "$OS_TYPE" != "Darwin" || "$ARCH_TYPE" != "arm64" ]]; then
    echo "Error: This script is only for Mac with Apple Silicon (M1, M2, M3...)."
    echo "Current System: $OS_TYPE ($ARCH_TYPE)"
    exit 1
fi

# 여기서부터 실제 스크립트 내용 시작
echo "Apple Silicon Mac detected. Proceeding..."

APP_NAME="ScriptLauncher"
APP_ID="com.chunsam.scriptlauncher"
ICON_PATH="app.png"

go install fyne.io/tools/cmd/fyne@latest

echo "🚀 빌드 시작: $APP_NAME..."

go mod tidy

if [ -f "$ICON_PATH" ]; then
    fyne package -os darwin -icon "$ICON_PATH" -id "$APP_ID" -name "$APP_NAME"
echo "✅ 빌드 완료! ${APP_NAME}.app 파일이 생성되었습니다."
else
    echo "⚠️  $ICON_PATH 파일이 없습니다. "
fi
