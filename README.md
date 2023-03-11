# zeus-client
This is an unofficial client for the ZEUS® time tracking system. It's a private project which is not implemented during working hours.

## Installation


## Usage

Run `zeus-client -help` to get a list of all available commands.

Example output:

```bash
zeus-client is a script to toggle the presence status in Zeus time tracking tool.

options:
  -debug
        Enables debug mode (default true)
  -dry-run
        Does a dry run without toggling the presence state
  -password string
        Password for the zeus time tracking tool. If not provided, the script will prompt for the password.
  -silent
        Runs the selenium script in headless mode (default true)
  -user string
        Username for the zeus time tracking tool

example usage:
  zeus-client -user <username>
```

## Error reporting

If you encounter an error, please consult the logs which can be found on the same directory as the binary.

## Disclaimer

By using the zeus-client, users acknowledge and accept that they are using it at their own risk and that we will not be held liable for any damages or losses resulting from the use of the client. It does not guarantee that the client switches the presence state correctly. The client is provided as is and without any warranty.