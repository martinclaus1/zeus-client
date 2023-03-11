## zeus-client toggle

Toggles the presence state in ZEUS® time tracking tool.

### Synopsis

Toggles the presence state in ZEUS® time tracking tool.
If the current state is "present", it will be set to "absent" and vice versa. 
In the end, the current state will be printed to the console.

```
zeus-client toggle [flags]
```

### Options

```
  -h, --help   help for toggle
```

### Options inherited from parent commands

```
  -d, --debug             Enable debug mode
      --headless          Runs the playwright script in headless mode (default true)
  -p, --password string   Password for the ZEUS time tracking tool. If not provided, the script will prompt for the password
  -u, --username string   Username for the ZEUS time tracking tool
```

### SEE ALSO

* [zeus-client](zeus-client.md)	 - zeus-client is a CLI for the ZEUS® time tracking tool.

