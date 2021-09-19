# Scaleway deployment 


```
cat key.txt | docker login scalewayRegistry -u nologin --password-stdin
docker build . -t scalewayRegistry/bot
docker push scalewayRegistry/bot
```