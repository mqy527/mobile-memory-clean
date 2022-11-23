#!/system/bin/sh
process_id=`ps|grep mobile-memory-clean|sed -r 's/ +/ /g'|cut -f2 -d" "`
if [[ -n ${process_id} ]];then
echo -e "begin stopping the process: $process_id"
`kill -9 ${process_id}`
echo -e "process stopped successfully"
else
echo -e "process stopped successfully"
fi
