# webot

一个智能的微信机器人

- [x] AI 回复（私聊 / 群组）
- [x] 自动总结聊天记录
- [ ] V2EX 热帖推送
- [ ] 每日精选新闻推送
- [ ] 每日一文 / 图

## Deploy

### docker-compose.yaml

```yaml
version: '3.0'
services:
  webot:
    container_name: webot
    image: sunls24/webot:latest
    user: 1000:1000
    environment:
      - STORAGE=data/storage.json
      - OPENAI_API_KEY=sk-xxx
      - DB_DSN=data/webot.db
    logging:
      options:
        max-size: 1m
    network_mode: host
    restart: unless-stopped
    volumes:
      - ./webot:/app/data:rw
```