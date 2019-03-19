[![Current Release](https://img.shields.io/github/release/daylioti/docker-commander.svg)](https://github.com/daylioti/docker-commander/releases/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/daylioti/docker-commander)](https://goreportcard.com/report/github.com/daylioti/docker-commander)

<code>docker-commander</code> is a cross-platform, customizable, execution commands in docker containers dashboard based on <a href="https://github.com/gizak/termui">termui</a>

<img src="./_examples/example.gif" ></img>

## Installation
Please use `1.0.3` version, master branch have a lot of new features and it not stable.

```bash
sudo wget https://github.com/daylioti/docker-commander/releases/download/1.0.3/docker-commander_linux_amd64 -O /usr/local/bin/docker-commander
sudo chmod +x /usr/local/bin/docker-commander
```

`docker-commander` is also avaliable for Arch in the <a href="https://aur.archlinux.org/packages/docker-commander">AUR</a>

### Options

Option | Description
--- | ---
-api-host| docker api host, f.e tcp://127.0.0.1:2376
-api-v | docker api version
-h	| display help dialog
-c  | path to yml config file
-v	| output version information and exit

### Keybindings

Key | Action
--- | ---
\<enter\> | Execute command
\<Left\>, \<Right\>, \<Up\>, \<Down\>, H, J, K, L  | Menu/terminal controls 
\<Tab\> | Switch between terminal and menu
q | Quit docker-commander

## Usage

`docker-commander` requires config file to build menu, default config path - ./config.yml.
 You can also use `-c` param to specify path to yml file.
### Docker api
By default `docker-commander` tries to find local docker api client or you can specify it with 
`-api-host` param
 
### Config file
 ```yaml
  config:
  - name: "menu item name"
    config:
    - name: "another child menu item"
      config:
      - name: "item with command"
        exec: # it just example.
          workdir: "/usr/src/site"
          from_image: "flatland-site-reactjs-public"
          cmd: "npm run-script start"
      - name: "item"
        config:
        - name: "item2"
          exec:
            from_image: "ubuntu"
            workdir: "/var"
            cmd: "ls -lah"      
  - name: "menu item 2"
    config:
    - name: "another child menu item 2"
      exec:
        from_image: "ubuntu"
        cmd: "ls -lah /var"
  ```
  Configuration of menu can be with any depth, yml support anchors for some optimizations
