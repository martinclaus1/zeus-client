# zeus-client
This is an unofficial client for the ZEUS® time tracking tool.

## Installation

zeus-client is platform independent and can be used on Windows, Linux and macOS.

- Option 1: Just download the binaries for your platform from the [releases page](https://github.com/martinclaus1/zeus-client/releases/latest).
- Option 2: Use `go install` to install the client directly: 
  ```bash
  # Don't forget to add the $GOPATH/bin directory to your $PATH, e.g. export PATH=$PATH:$(go env GOPATH)/bin
  go install github.com/martinclaus1/zeus-client@latest
  ```

## Documentation

The documentation can be found [here](documentation/zeus-client.md).

## Error reporting

If you encounter an error, please consult the logs which can be found in `$HOME/.zeus-client/logs/` folder.

## Disclaimer

By using the zeus-client, users acknowledge and accept that they are using it at their own risk and that we will not be held liable for any damages or losses resulting from the use of the client. 
It does not guarantee that the client switches the presence state correctly. The client is provided as is and without any warranty.