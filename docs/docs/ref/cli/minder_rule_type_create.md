---
title: minder rule type create
---
## minder rule_type create

Create a rule type within a minder control plane

### Synopsis

The minder rule type create subcommand lets you create new rule types for a project
within a minder control plane.

```
minder rule_type create [flags]
```

### Options

```
  -f, --file stringArray   Path to the YAML defining the rule type (or - for stdin). Can be specified multiple times. Can be a directory.
  -h, --help               help for create
```

### Options inherited from parent commands

```
      --config string            Config file (default is $PWD/config.yaml)
      --grpc-host string         Server host (default "api.stacklok.com")
      --grpc-insecure            Allow establishing insecure connections
      --grpc-port int            Server port (default 443)
      --identity-client string   Identity server client ID (default "minder-cli")
      --identity-realm string    Identity server realm (default "stacklok")
      --identity-url string      Identity server issuer URL (default "https://auth.stacklok.com")
```

### SEE ALSO

* [minder rule_type](minder_rule_type.md)	 - Manage rule types within a minder control plane

