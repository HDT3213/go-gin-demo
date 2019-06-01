A web backend demo powered by go-gin

You can register, post text and follow others in this app, all your posts will be push into your followers timeline

it requires:

- mysql (required)
- redis (required)
- rabbitmq (optional, fall back to goroutine)
- zookeeper (optional, necessary for canal)

please read and edit `config.yml` to make them work properly.