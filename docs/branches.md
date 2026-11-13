# Project split

`master` contains the Go terminal application, including Go audio synthesis.
The local `raylib` branch has been removed. Its history remains recoverable at
`archive/raylib-split` (commit `757ab4a8c31eb6fa5f1004903ee83cd7e2269c8b`).

The desktop spinoff is `alt-synth`, located at `/home/jm/Projects/alt-synth`.
Its Git history retains the Raylib branch tip above. Development continues
independently from the terminal application.

Remote deletion must be applied from a terminal with access to GitHub:

```sh
git push origin --delete raylib
```

Older remote experiments (`feature/raylib-gui` and `feat/c-raylib-ui`) remain
untouched. The archive tag above is local and has not been published.
