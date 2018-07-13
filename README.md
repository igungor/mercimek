![](assets/mercimek.png)

# mercimek

`mercimek` is a Telegram bot that counts [lemna
minors](https://en.wikipedia.org/wiki/Lemna_minor) from a given photo. It's a
simple wrapper around the awesome [Fiji](https://imagej.net/Fiji) image
processing toolset.

## installation

```
go get -u github.com/igungor/mercimek
```

## run

You'll need to download Fiji. See the [downloads page](https://imagej.net/Fiji/Downloads).

After you install Fiji, set the `binary-path` in mercimek.conf to the Fiji executable's full path.
Since this is a telegram bot, you'll also need a Telegram token. See
https://core.telegram.org/bots.

If everything is set, run:

`mercimek -c mercimek.conf`

[See it in action.](assets/telegram.png)

## license

MIT. See LICENSE.
