## Usage

```yaml
environments:
    DB_TYPE: mysql
    DB_URI: "root:123@tcp(db)/turnik?charset=utf8mb4&parseTime=True"
    TELEGRAM_TOKEN: YOUTOKETFROMBOTFATHER
    TELEGRAM_URL: https://api.telegram.org #or your proxy http://IP:30012 @see https://github.com/fullpipe/tele-proxy
```

## Local 

```sh
DB_TYPE=sqlite3 DB_URI=var/test.db TELEGRAM_TOKEN=YOUTOKETFROMBOTFATHER TELEGRAM_URL=https://api.telegram.org justrun -c 'go run .' .
```
