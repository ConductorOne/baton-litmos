# baton-litmos
Welcome to your new connector! To start out, you will want to update the dependencies.
Do this by running `make update-deps`.

```
baton-litmos

Usage:
  baton-litmos [flags]
  baton-litmos [command]

Available Commands:
  capabilities       Get connector capabilities
  completion         Generate the autocompletion script for the specified shell
  help               Help about any command

Flags:
      --api-key string            required: API Key ($BATON_API_KEY)
      --client-id string          The client ID used to authenticate with ConductorOne ($BATON_CLIENT_ID)
      --client-secret string      The client secret used to authenticate with ConductorOne ($BATON_CLIENT_SECRET)
  -f, --file string               The path to the c1z file to sync with ($BATON_FILE) (default "sync.c1z")
  -h, --help                      help for baton-litmos
      --limited-courses strings   Limit imported sources to a specific list by Course ID ($BATON_LIMITED_COURSES)
      --log-format string         The output format for logs: json, console ($BATON_LOG_FORMAT) (default "json")
      --log-level string          The log level: debug, info, warn, error ($BATON_LOG_LEVEL) (default "info")
  -p, --provisioning              This must be set in order for provisioning actions to be enabled ($BATON_PROVISIONING)
      --skip-full-sync            This must be set to skip a full sync ($BATON_SKIP_FULL_SYNC)
      --source string             required: Source ($BATON_SOURCE)
      --ticketing                 This must be set to enable ticketing support ($BATON_TICKETING)
  -v, --version                   version for baton-litmos

Use "baton-litmos [command] --help" for more information about a command.
```
