PLIST_NAME=$1
LOCATION=$2

go build -o mac-notifier
mkdir -p /opt/mac-notifier
cp -a -v ./mac-notifier /opt/mac-notifier/
cp -a -v ./config.json /opt/mac-notifier/config.json
rm -r -f ./mac-notifier

cat > ${LOCATION}/${PLIST_NAME}.plist <<-EOM
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
        <key>Label</key>
        <string>${PLIST_NAME}</string>
        <key>ServiceDescription</key>
        <string>Custom notification service</string>
        <key>ProgramArguments</key>
        <array>
                <string>/opt/mac-notifier/mac-notifier</string>
        </array>
        <key>RunAtLoad</key>
    <true/>
</dict>
</plist>
EOM

launchctl load ${LOCATION}/${PLIST_NAME}.plist
launchctl start ${PLIST_NAME}