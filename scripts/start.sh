cd /var/tyrant-bot
nohup ./api > api.log 2>&1 &
echo $! > api.pid
nohup ./bot > bot.log 2>&1 &
echo $! > bot.pid
