# download-artifact
Golang code to download artifacts hosted on Nexus using Lucene and nexus rest api


![Golang CI](https://github.com/sfragata/download-artifact/workflows/Golang%20CI/badge.svg)

## Usage

```
Usage of doownload-artifact
  -appName string
    	name to be used when dowloading
  -artifactId string
    	maven artifact id
  -classifier string
    	artifact classifier
  -groupId string
    	maven group id
  -nexus string
    	base nexus url
  -packaging string
    	packaging (default "war")
  -target string
    	target folder (default "/tmp")
  -version string
    	artifact version
```        