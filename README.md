## GoBin
Just a simple HTTP Pastebin using GoLang

There is no API but you can POST whatever you want to `/`

The server will accept your file if its under 10MB

It compresses pastes using using Cloudflare's Pako version of GZip on the browser side.

Files get decompressed when they get accessed

Stores pastes using random 3 letter file names in current directory

### """API"""

Access UI: `GET http://localhost`

Read paste: `GET http://localhost/000` 

Create paste: `POST http://localhost/` with your GZip content as Body
