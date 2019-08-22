# This script queries https://harmony.one/balances.json and https://harmony.one/1h.json for a Harmony Node and sends the output to a Telegram chat.
# Provided without warranty by https://chainflow.io/staking. Use at your own risk and may it be of benefit to you and the Harmony network.

# Replace [YOUR TELEGRAM BOT TOKEN], [YOUR TELEGRAM CHAT ID] and [YOUR HARMONY NODE ADDRESS] with your information. Then install this scrip as a cron job to run at your preferred interval.

#!/bin/bash

TOKEN=[YOUR TELEGRAM BOT TOKEN]
CHAT_ID=-[YOUR TELEGRAM CHAT ID]
URL="https://api.telegram.org/bot$TOKEN/sendMessage"
MESSAGE="Harmony Monitoring Bot\n"

# 1 - Process the log file.

# 2
ONEsPerHour=$(curl -s  https://harmony.one/1h.json | jq -r '.onlineNodes[] | select(.address=="[YOUR HARMONY NODE ADDRESS]") | .ONEsPerHour')

if [ $(echo "${ONEsPerHour} > 0" | bc) ]
then
  MESSAGE+="ONEsPerHour: ${ONEsPerHour}\n"
else
  MESSAGE+="No ONEsPerHour, may be node is offline\n"
fi

# 3
BALANCE=$(curl -s https://harmony.one/balances.json | jq -r '.onlineNodes[], .offlineNodes[] | select(.address=="[YOUR HARMONY NODE ADDRESS]") | .totalBalance')

if [ $(echo "${BALANCE} > 0" | bc) ]
then
  MESSAGE+="BALANCE: ${BALANCE}\n"
else
  MESSAGE+="No BALANCE, node may be offline\n"
fi

echo $ONEsPerHour
echo $BALANCE
echo -e $MESSAGE

# Send to Telegram
curl -s -X POST $URL -d chat_id=$CHAT_ID -d text="$(echo -e $MESSAGE)"
