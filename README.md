#petermcconnell.com

Super simple Go app designed to run on AppEngine.

- Handlers in `./main.go`
- Static assets in `./static`
- HTML templates in `./templates`

Deploy as usual for AppEngine: `appcfg.py update -A <appID> -V <version> .`

**Note** the code in this project hasn't been given much care. I wanted to put this together as a quick example of how easy it was to deploy Golang onto Appengine.
