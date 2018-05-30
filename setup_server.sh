#启动redis服务
redis-server ./conf/redis.conf

#启动trackerd

fdfs_trackerd  /home/itcast/go/src/ilhome/conf/tracker.conf restart
#启动storaged

fdfs_storaged  /home/itcast/go/src/ilhome/conf/storage.conf restart