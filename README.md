# Turnik bot

You need horizontal bar (turnik) in your office.  
And your need subscribe on bot. Its possible  
to launch your own bot. Or use mine.

## Deployment

See [deployment.yml](deployment.yml), or [docker-compose.yml](docker-compose.yml).

### Envars

```yaml
# Mysql
DB_TYPE: mysql
DB_URI: "root:123@tcp(db)/turnik?charset=utf8mb4&parseTime=True"
TELEGRAM_TOKEN: YOUTOKETFROMBOTFATHER
TELEGRAM_URL: https://api.telegram.org #or your proxy http://IP:30012 @see https://github.com/fullpipe/tele-proxy

# sqlite
DB_TYPE: sqlite3 # mount volume at /app/data
TELEGRAM_TOKEN: YOUTOKETFROMBOTFATHER
TELEGRAM_URL: https://api.telegram.org #or your proxy http://IP:30012 @see https://github.com/fullpipe/tele-proxy
```

## Local development

```sh
DB_TYPE=sqlite3 DB_URI=data/test.db TELEGRAM_TOKEN=YOUTOKETFROMBOTFATHER TELEGRAM_URL=https://api.telegram.org justrun -c 'go run .' .
```

## TODO

- time zone handling
- add translations
- move motivations to extenal json|yaml file
- add more randomness
- add more motivations
