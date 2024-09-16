# mitmproxy
A pipeline-based stream editor for remote web pages that provides text filters (like ```sed```, ```awk``` and ```nkf```) and HTML filters.

## Usage

### case 1
```http://mitmproxy-server/s/asahi/akahi/https://www.akahi.com/```

This fetches https://www.asahi.com/ after the request URI is rewritten.

**Request headers and body:** Every occurrence of ```akahi``` is replaced by ```asahi```, hiding the altered string ```akahi``` from the upstream.

**Response headers and body:** Every occurrence of ```asahi``` is replaced by ```akahi```.

### case 2
```http://mitmproxy-server/proxy:pass/http://www.akahi.com/https://www.asahi.com/```

This fetches https://www.asahi.com/ and ensures hyperlinks are rewritten.

**Request headers and body:** Every occurrence of ```http://www.akahi.com``` (origin) and ```akahi.com``` (domain) is respectively replaced by ```https://www.asahi.com``` and ```asahi.com```.

**Response headers and body:** Every occurrence of ```https://www.asahi.com``` (origin) and ```asahi.com``` (domain) is respectively replaced by ```http://www.akahi.com``` and ```akahi.com```.

### case 3

```http://mitmproxy-server/remove:script/nkf/https://www.itmedia.co.jp/```

This fetches https://www.itmedia.co.jp/, a ```Shift_JIS``` page.

**Request headers and body:** Not modified.

**Response headers and body:** A decoded ```UTF-8``` page after the string ```charset=Shift_JIS``` occurs in the page itself.

**Response body:** Every occurrence of the tag ```<script>``` is removed.

## List of filters

### /charset=\<encoding\>/\<url\>

**Request headers and body:** Encoded from ```UTF-8``` to ```<encoding>```.

**Response headers and body:** Decoded from ```<encoding>``` to ```UTF-8```.

### /nkf/\<url\>

**Request headers and body:** Not modified.

**Response headers and body:** If the string ```charset=<encoding>``` occurs in the response headers or body, it is decoded from ```<encoding>``` to ```UTF-8```.

### /s/\<original\>/\<altered\>/\<url\>

**Request headers and body:** Every occurrence of ```<altered>``` is replaced by ```<original>```.

**Response headers and body:** Every occurrence of ```<original>``` is replaced by ```<altered>```.

### /proxy:pass/\<scheme\>://\<host\>/\<url\>

This is useful for creating a mirror site and ensures hyperlinks are rewritten.

**Request headers and body:** Every occurrence of ```http://www.mirror.site``` (origin) and ```mirror.site``` (domain) is respectively replaced by ```https://www.original.site``` and ```original.site```, hiding the origin of the mirror site from the upstream.

**Response headers and body:** Every occurrence of ```https://www.original.site``` (origin) and ```original.site``` (domain) is respectively replaced by ```http://www.mirror.site``` and ```mirror.site```.

### /proxy:\<filter\>/\<scheme\>://\<host\>/\<url\>

Equivalent to ```/proxy:pass/<scheme>://<host>/<url>``` following another filter.

### /awk=\<program\>/\<url\>

**Request headers and body:** Not modified.

**Response body:** Rewrited by an AWK program.

### /awk:csv=\<program\>/\<url\>

**Request headers and body:** Not modified.

**Response body:** Rewrited by an AWK program run with CSV mode.

### /awk:tsv=\<program\>/\<url\>

**Request headers and body:** Not modified.

**Response body:** Rewrited by an AWK program run with TSV mode.

### /minify/\<url\>

**Request headers and body:** Not modified.

**Response body:** A minified HTML, CSS, JS, etc. Some special characters and HTML closing tags may be omitted.

### /minify:keep/\<url\>

**Request headers and body:** Not modified.

**Response body:** A minified HTML, CSS, JS, etc, keeping special characters and HTML closing tags.

### /remove=\<selector\>/\<url\>

**Request headers and body:** Not modified.

**Response body:** Every element that matches the CSS selector is removed.

### /remove:keep=\<CSS selector\>/\<url\>

**Request headers and body:** Not modified.

**Response body:** Every tag that matches the CSS selector is removed, keeping its inner content.

### /remove:comment/\<url\>

**Request headers and body:** Not modified.

**Response body:** Every comment tag is removed.

### /remove:script/\<url\>

**Request headers and body:** Not modified.

**Response body:** Equivalent to ```/remove=script, link[as=script]``` followed by ```/remove:keep=noscript/<url>```.

### /geo=\<latitude\>:\<longitude\>

**Response body:** Japanese address of the location.

### /youtube/\<id\>

**Response body:** YouTube raw video/audio URLs.

### /youtube:type=\<mime\>/\<id\>

**Response body:** YouTube raw video/audio URLs that matches ```<mime>```.
