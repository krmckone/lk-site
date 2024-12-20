# lk-site
LK-Site, is a custom, basic, static site generator. It takes a combination of markdown files and raw HTML and creates a static site as output. There's a lot of hard-coded implementation that needs to get cleaned up and less reliance on writing raw HTML.

The raw files for the site are located under `assets`. The Go tool that pulls everything together is located under `cmd`. Packages containing custom logic are implemented under `internal`.

To run the site locally under `localhost:8080` with auto-reloading:
```shell
air -c .air.toml
```
