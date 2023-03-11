# zeus-client
This is an unofficial client for the ZEUS® time tracking system. It's a private project which is not implemented during working hours.

## Installation


## Usage

Run `zeus-client --help` to get a list of all available commands.

Example output:

```bash
zeus-client is a CLI for the ZEUS time tracking tool. 
It allows you to toggle your presence state and get information about your current presence state, for example:
zeus-client toggle
zeus-client get status

Usage:
  zeus-client [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  get         Gets information from ZEUS time tracking tool
  help        Help about any command
  toggle      Toggles the presence state in ZEUS time tracking tool.

Flags:
  -d, --debug             Enable debug mode
      --headless          Runs the playwright script in headless mode (default true)
  -h, --help              help for zeus-client
  -p, --password string   Password for the ZEUS time tracking tool. If not provided, the script will prompt for the password
  -u, --username string   Username for the ZEUS time tracking tool

Use "zeus-client [command] --help" for more information about a command.
```

## Error reporting

If you encounter an error, please consult the logs which can be found on the same directory as the binary.

## Disclaimer

By using the zeus-client, users acknowledge and accept that they are using it at their own risk and that we will not be held liable for any damages or losses resulting from the use of the client. It does not guarantee that the client switches the presence state correctly. The client is provided as is and without any warranty.