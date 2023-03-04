<style>
                  
  *{
  margin: 0;
  padding: 0;
  }
  .tpt-bar {
    display:flex;
    /*border:1px solid #e2e2e2;*/
    /*border-radius:2px;*/
    /*background:#c6fced;*/
    /*box-shadow:0 2px 5px 0 rgba(0,0,0,.1);*/
    flex-wrap:wrap;
    /*width: 80%;*/
    width: 100%;
    margin: 0 auto;
  }
  .tpt-bar label {
    display:block;
    padding:0 20px;
    height:50px;
    line-height:50px;
    cursor:pointer; 
    order:1;
  }
  .tpt-bar .tpt-bar-con {
    /*z-index:0;*/
    display:none;
    /*padding:30px;*/
    width:100%;
    min-height:120px;
    line-height: 30px;
    /*border-top:1px solid #e2e2e2;*/
    margin-top: -1px;
    /*background:#f3f3f4;*/
    order:99;
  }
  .tpt-bar input[type=radio] {
    position:absolute;
    opacity:0;
  }
  .tpt-bar input[type=radio]:checked+label {
    /*z-index:1;*/
    /*margin-right:-1px;*/
    border-bottom: 1px solid #40a9ff;                  
    /*margin-left:-1px;*/
    /*border-right:1px solid #e2e2e2;*/
    /*border-left:1px solid #e2e2e2;*/
    /*background:#69d6e8;*/
  }
  .tpt-bar input[type=radio]:checked+label+.tpt-bar-con {
    display:block;
  }
</style>



# go-storage
分布式存储
参考：https://github.com/stuarthu/go-implement-your-object-storage

1. 准备数据

    dd if=/dev/zero of=./test bs=1m count=100
  
    openssl dgst -sha1 -binary ./test | base64 (MAC环境)
  
    openssl dgst -sha1 --binary ./test | base64 （Linux）

2. 上传

  curl -v 'http://localhost:9999/file/upload?fileName=test' -XPUT --data -binary @./test -H "Digests:" （MAC）
  
  curl -v 'http://localhost:9999/file/upload?fileName=test' -XPUT --data-binary @./test -H "Digests:" (Linux)
  
3. 下载
  
  curl -v 'http://localhost:9999/file/download?fileName=test&fileVersion=1' -o ./output
  
  curl -v 'http://localhost:9999/file/download?fileName=test&fileVersion=1' -H "Accept-Encoding: gzip" -o ./output2.gz (gzip方式下载)







<div class="tpt-bar">
  <input type="radio" name="bar" id="tab-1" checked="">
  <label for="tab-1">XML</label>
  <div class="tpt-bar-con">

```
123
```

  </div>
  <input type="radio" name="bar" id="tab-2">
  <label for="tab-2">JSON</label>
  <div class="tpt-bar-con">

   ```
   456
   
   
   
   ```

  </div>
</div>

- ldla
- lall
