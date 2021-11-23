# ! /bin/sh

##nohup ./start.sh >/dev/null 2>&1 &

SERVICE_PATH="/Users/kennshi/Desktop/xbull"

SERVICE_NAME="xbull"

START_CMD="./$SERVICE_NAME -log -dsn root:13811237916sS@tcp(localhost:3306)/bridge“

LOG_FILE="restart.log"


cd $SERVICE_PATH

pwd

while true 

do

    procnum=`ps -ef|grep $SERVICE_NAME|grep -v grep|wc -l`

    if [ $procnum -eq 0 ]

    then

        echo "start service...................."

        echo `date +%Y-%m-%d` `date +%H:%M:%S`  $SERVICE_NAME >>$LOG_FILE

        ${START_CMD}
	
	else
		echo "wait for next monitor..............."

    fi

    sleep 1

done
