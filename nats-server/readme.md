docker run -p 4222:4222 -p 8222:8222 nats:latest

{action}.{service}...

### Send command
.Request

> command.noolite.13.on
> command.noolite-f.13.on
< ok/fail

### Emit event
event.noolite.{channel}.{command}
