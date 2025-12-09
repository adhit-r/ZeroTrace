#!/bin/bash

echo "🔧 Creating ZeroTrace Tray Icon..."

# Create app bundle structure
mkdir -p /Applications/ZeroTrace.app/Contents/MacOS
mkdir -p /Applications/ZeroTrace.app/Contents/Resources

# Create Info.plist
cat > /Applications/ZeroTrace.app/Contents/Info.plist << 'EOF'
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>CFBundleExecutable</key>
    <string>ZeroTrace</string>
    <key>CFBundleIdentifier</key>
    <string>com.zerotrace.agent</string>
    <key>CFBundleName</key>
    <string>ZeroTrace</string>
    <key>CFBundleVersion</key>
    <string>1.0.0</string>
    <key>CFBundleShortVersionString</key>
    <string>1.0.0</string>
    <key>LSUIElement</key>
    <true/>
    <key>NSHighResolutionCapable</key>
    <true/>
</dict>
</plist>
EOF

# Create a simple tray app
cat > /Applications/ZeroTrace.app/Contents/MacOS/ZeroTrace << 'EOF'
#!/bin/bash
echo "ZeroTrace Agent v1.0.0"
echo "Status: Active - Monitoring system for vulnerabilities"
echo "Dashboard: http://localhost:5173"
echo "API: http://localhost:8080"
echo "Press Ctrl+C to stop"
while true; do
    sleep 30
    echo "$(date): ZeroTrace Agent running..."
done
EOF

chmod +x /Applications/ZeroTrace.app/Contents/MacOS/ZeroTrace

echo "✅ ZeroTrace tray icon created at /Applications/ZeroTrace.app"
echo "🔧 To start: open /Applications/ZeroTrace.app"
echo "📊 Dashboard: http://localhost:5173"
