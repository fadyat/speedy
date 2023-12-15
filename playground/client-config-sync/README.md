### Client config sync

Current example shows how is config synchronized between client and server.

```shell
go run main.go
```

`main.go` launches three servers + one client.

- All servers have the same configuration.
- Client have a deprecated configuration, which is needed to be updated.
- In some point of the time, servers are updated with a new configuration.

After some time client configuration is updated to the latest one.

#### Files description:

- `f0client.yaml` - client initial configuration.
- `f1server.yaml` - server configuration.
- `f2server.yaml` - server new configuration.
- `server-copy.yaml` - temporary file for sharing server configuration.