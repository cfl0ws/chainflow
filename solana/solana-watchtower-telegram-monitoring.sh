#!/bin/bash

TOKEN=[YOUR BOT TOKEN]
CHAT_ID=-[YOUR CHAT ID]
URL="https://api.telegram.org/bot$TOKEN/sendMessage"
MESSAGE="Solana validator not available"

solana-watchtower --url http://34.82.79.31 --validator-identity ~/validator-keypair.json
RESULT=$?

echo -e "Grep identified as: $RESULT"

if [ $RESULT == 1 ]; then
  echo -e $MESSAGE
  # Send to Telegram
  curl -s -X POST $URL -d chat_id=$CHAT_ID -d text="$(echo -e $MESSAGE)"
fi
