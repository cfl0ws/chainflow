#!/bin/bash

TOKEN=[YOUR BOT TOKEN]
CHAT_ID=-[YOUR CHAT ID]
URL="https://api.telegram.org/bot$TOKEN/sendMessage"
MESSAGE="Solana validator not available"

solana show-validators | grep 8bRCnytB7bySmqxodNGbZuUAtncKkB8T733DD1Dm9WMb
RESULT=$?

echo -e "Grep identified as: $RESULT"

if [ $RESULT == 1 ]; then
  echo -e $MESSAGE
  # Send to Telegram
  curl -s -X POST $URL -d chat_id=$CHAT_ID -d text="$(echo -e $MESSAGE)"
fi
