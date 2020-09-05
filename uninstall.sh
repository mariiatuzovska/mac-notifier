PLIST_NAME=$1
LOCATION=$2

launchctl unload ${LOCATION}/${PLIST_NAME}.plist
rm -r -f ${LOCATION}/${PLIST_NAME}.plist
rm -r -f /opt/mac-notifier