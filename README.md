# mercimek

a telegram bot that counts lemna minors from a given photo.

## installation

go get -u github.com/igungor/mercimek

## run

mercimek is a wrapper around [Fiji](https://imagej.net/Fiji), so you'll need to download Fiji. See
the [downloads page](https://imagej.net/Fiji/Downloads).

After you install Fiji, set the `binary-path` in mercimek.conf to the Fiji executable's full path.
Since this is a telegram bot, you'll also need a Telegram token. See
https://core.telegram.org/bots.

If everything is set, run:

`mercimek -c mercimek.conf`

## license

MIT. See LICENSE.
