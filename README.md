# Build
```bash
bee pack -be GOOS=linux GOARCH=amd64
GOOS=linux GOARCH=amd64
```

## 定时任务
### 同步群信息
*/1 * * * * /usr/bin/wget -q --delete-after http://127.0.0.1:9994/we-sns/crawler/task



