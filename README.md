# download-artifact
Golang code to download artifacts hosted on Nexus using Lucene and nexus rest api


![Golang CI](https://github.com/sfragata/download-artifact/workflows/Golang%20CI/badge.svg)

## Installation

### Mac

```
brew install sfragata/tap/download-artifact
```

### Linux and Windows

get latest release [here](https://github.com/sfragata/download-artifact/releases)

## Usage

```
download-artifact - Utility to download artifacts hosted on Nexus using Lucene and nexus rest api

  Flags:
       --version            Displays the program version string.
    -h --help               Displays help with available flag, subcommand, and positional value parameters.
    -a --artifact-id        Maven artifact id
    -g --group-id           Maven group id
    -v --artifact-version   Artifact version
    -p --packaging          Type of packaging (ex. pom, jar, war etc)
    -n --appName            Name to be used as filename when dowloading
    -c --classifier         Artifact classifier
    -t --target             Target folder (default: .)
    -H --host               Base nexus url
    -r --repository         Nexus repository id
    -V --verbose            Verbose mode
    -nv --nexus-version     Nexus version (default: 3)
```        