#go-omxremote

Control raspberry pi omxplayer from the browser (including mobile browsers). To install just [download the executable](https://github.com/dplesca/go-omxremote/releases/download/v2.0/go-omxremote) and run it. For help run it with the `-h` flag. Example usage (you can of course add in your path):

`./go-omxremote -bind :some-port -media path/to/video/files`

Command flags:

```
-bind string
    Address to bind on. If this value has a colon, as in ":8000" or
            "127.0.0.1:9001", it will be treated as a TCP address. If it
            begins with a "/" or a ".", it will be treated as a path to a
            UNIX socket. If it begins with the string "fd@", as in "fd@3",
            it will be treated as a file descriptor (useful for use with
            systemd, for instance). If it begins with the string "einhorn@",
            as in "einhorn@0", the corresponding einhorn socket will be
            used. If an option is not explicitly passed, the implementation
            will automatically select among "einhorn@0" (Einhorn), "fd@3"
            (systemd), and ":8000" (fallback) based on its environment. 
            (default ":8000")
-media string
    path to look for videos in (default ".")
```

The project is geared towards mobile usage, it has been tested on both Android and iOS devices.

###Modify it

Generate react components file with babel:  

`babel --presets react components/ --minified -o assets/all.js`

Generate assets file using esc:  

`esc -o assets.go -prefix="assets" assets views`

Build again:

`go build`

###Credits

It's written in go, uses [httprouter](https://github.com/julienschmidt/httprouter) as a router, [color](https://github.com/fatih/color) for colorized output and [esc](https://github.com/mjibson/esc) to generate and embed assets in go source files. The front-end is written in [react](http://facebook.github.io/react/), the style uses [skeleton](http://getskeleton.com/).

###Screenshot

![Android](http://s10.postimg.org/6susaybqh/screen_p.png)
