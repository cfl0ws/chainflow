# This script will send a message to a Telegram chat if the specified validator vote account pubkey IS NOT in the current voting set.

#!/bin/bash

TOKEN=[YOUR TELEGRAM BOT TOKEN]
CHAT_ID=-[YOUR TELEGRAM CHAT ID]
URL="https://api.telegram.org/bot$TOKEN/sendMessage"
MESSAGE="Solana Validator NOT Voting"

[YOUR PATH TO]/solana show-validators | grep [YOUR VALIDATOR VOTE ACCOUNT PUBKEY]
RESULT=$?

echo -e "Grep identified as: $RESULT"

if [ $RESULT == 1 ]; then
  echo -e $MESSAGE
  # Send to Telegram
  curl -s -X POST $URL -d chat_id=$CHAT_ID -d text="$(echo -e $MESSAGE)"
fi
