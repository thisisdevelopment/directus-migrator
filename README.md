# directus-migrator
Quick commandline tool to migrate directus schemas

This tool has a few commands 

* **status** - checks the status between two environements and shows the diff
* **apply** - checks the status between two environements and applies the diff, this will automatically create a backup of your target environment schema
* **backup** - backup a schema of environment 
* **restore** - restores a schema to an environment (currently unsupport, check command output)

## Usage
First create a config file in the ./config directory. under **envs**: you can specify your environments (check **config/config.example.yaml** for reference)

### Commandline
to check the status between dev and test environment
```
# ./directus-migrator status -s dev -t test
```

> to apply changes from test to your dev environment
```
# ./directus-migrator apply -s test -t dev
```

#### reference
```
Usage:
  directus-migrator [command]

Available Commands:
  apply       apply diff
  backup      backup schema of given env
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  restore     Restores a backup file made by apply command to given env
  status      check status of env vs other env

Flags:
      --config string   config file (default is ./config/config.yaml)
  -f, --force           forces proceed on version mismatches, use at own risk
  -h, --help            help for directus-migrator
  -s, --source string   source environment (default "dev")
  -t, --target string   target environment

Use "directus-migrator [command] --help" for more information about a command.
```

# Contributing 
You can help to deliver a better fanunmarshaller, check out how you can do things [CONTRIBUTING.md](CONTRIBUTING.md)

# License 
© [This is Development BV](https://www.thisisdevelopment.nl), 2023
~time.Now()
Released under the [Apache License](https://github.com/thisisdevelopment/fanunmarshal/blob/master/LICENSE)
