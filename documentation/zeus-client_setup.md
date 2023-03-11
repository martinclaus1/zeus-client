## zeus-client setup

Sets up the username and password for the ZEUS® time tracking tool.

### Synopsis

Sets up the username and password for the ZEUS® time tracking tool. 
The username and the password will be stored in an encrypted file. 
Therefore, a folder '.zeus-client' will be created in your home directory. 
Running the setup command will overwrite the existing credentials. 
The machine id is used as the encryption key.

```
zeus-client setup [flags]
```

### Options

```
  -h, --help   help for setup
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

