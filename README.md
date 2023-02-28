# go-storage
分布式存储
参考：https://github.com/stuarthu/go-implement-your-object-storage

1. 准备数据
  dd if=/dev/zero of=./test bs=1m count=100
  openssl dgst -sha1 -binary ./test | base64 (MAC环境)
  
  openssl dgst -sha1 --binary ./test | base64 （Linux）
2. 上传
  curl -v 'http://localhost:9999/file/upload?fileName=test' -XPUT --data -binary @./test -H "Digests:"
