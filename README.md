# go-omxremote

Control raspberry pi omxplayer from the browser (including mobile browsers). It has absolutely zero dependencies. To install just [download the latest release](https://github.com/dplesca/go-omxremote/releases) and run it. For help run it with the `-h` flag. Example usage (you can of course add it in your path and run it as a systemd service, unit file example below):

`./go-omxremote -bind :some-port -media path/to/video/files`

To play youtube videos you need to install the python lib YouTube-DL [download here](https://youtube-dl.org/)

Command flags:

```
-bind string
    Address to bind on. If this value has a colon, as in ":8000" or
            "127.0.0.1:9001", it will be treated as a TCP address.
            (default ":31415")
-media string
    path to look for videos in (default ".")

-omx string
    omx options (default "-o hdmi")
-yt string
    youtube-dl options (default "-g")
```

The project is geared towards mobile usage, the interface has been tested on both Android and iOS devices. 

### Exposed HTTP API

This application exposes a simple, minimal HTTP API for interfacing with the player. The exposed endpoints are:
 - GET /files - get all video files as an array of objects with two fields: `file`, `hash`; example response:
 ```json
 [{"file":"test.avi","hash":"dGVzdC5hdmk="},{"file":"test.mp4","hash":"dGVzdC5tcDQ="}]
```
 - POST /start/:hash - starts playback for a video file
 - POST /player/pause - toggles pause/play for video playback
 - POST /player/stop - stops playback
 - POST /player/forward - seeks forward in video
 - POST /player/backward - seeks backward in video
 - POST /player/subs - switches subs stream

 One important note is that the `start` call does not return until the playback is stopped, so your request will probably time out.

### Example systemd unit file

A minimal systemd unit file. It makes `go-omxremote` a service so you can control it between sessions and/or don't need a tmux/screen to keep it open.

```
[Unit]
Description="omxremote go web service"

[Service]
User=pi
Group=pi
Restart=on-failure
ExecStart=/usr/local/bin/go-omxremote -media /path/to/media/files
WorkingDirectory=/home/pi

[Install]
WantedBy=multi-user.target
```

### Modify it

 - Clone repo
 - `npm install`
 - `npm run build` (after some changes have been made to front-end files)
 - regenerate assets file using [esc](https://github.com/mjibson/esc): `esc -o assets.go -prefix="dist" dist views`
 - build again: `go build`

### Credits

It's written in go, uses [httprouter](https://github.com/julienschmidt/httprouter) as a router and [esc](https://github.com/mjibson/esc) to generate and embed assets in go source files. The front-end is written in [svelte](https://svelte.dev/), the style uses [tailwind](http://https://tailwindcss.com/).

### Screenshot

![Android](https://i.imgur.com/ZRff2I2.png)