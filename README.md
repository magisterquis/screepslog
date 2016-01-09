screepslog
==========
Prints screeps console logs

Pretty much a direct ripoff of `https://github.com/TooAngel/screeps-cli`.

Installation
------------
```bash
go install github.com/kd5pbo/screepslog
```

Pre-compiled windows binaries available upon request.

Usage
-----
```bash
screepslog -u your@emailaddres.com
```

To give it a password, you can either type it when prompted or feed it a file
(hopefully with `400` permissions) on standard in:
```bash
screepslog -u your@emailaddres.com < .screepspw
```
