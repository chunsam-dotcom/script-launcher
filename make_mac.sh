#!/bin/bash

APP_NAME="ScriptLauncher"
APP_ID="com.chunsam.scriptlauncher"
ICON_PATH="app.png"

echo "🚀 빌드 시작: $APP_NAME..."

go mod tidy

if [ -f "$ICON_PATH" ]; then
    fyne package -os darwin -icon "$ICON_PATH" -id "$APP_ID" -name "$APP_NAME"
echo "✅ 빌드 완료! ${APP_NAME}.app 파일이 생성되었습니다."
else
    echo "⚠️  $ICON_PATH 파일이 없습니다. "
fi
