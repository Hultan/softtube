# SoftTube

SoftTube is a YouTube client that downloads the videos and thumbnails using a youtube-dl fork (yt-dlp), and keeps track of subscriptions using RSS feeds.

For now, it requires a specific MySQL-database, that cannot be created from the application, so the application is of no use for anyone else than me.

## TODO

### Playlist download:

Command : **yt-dlp --get-id --flat-playlist PLB6x_4-4tcYP_aCh66KkYdO1mH_XdSIVD**

A playlist that starts with **RD** is autogenerated, and cannot be downloaded.

A real playlist starts with **PL** and should be downloadable.