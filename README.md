# mitmproxy
A pipeline-based stream editor for remote web pages that provides text filters (like ```sed```, ```awk``` and ```nkf```) and HTML filters.

## Usage

### case 1
```http://mitmproxy-server/s/asahi/akahi/https://www.akahi.com/```

This fetches https://www.asahi.com/, but every occurrence of ```akahi``` in the request (like the request URI ```www.akahi.com```)  and ```asahi``` in the response is respectively replaced by ```asahi``` and ```akahi```.

### case 2
```http://mitmproxy-server/proxy:pass/http://www.akahi.com/https://www.asahi.com/```

This fetches https://www.asahi.com/, ensuring hyperlinks are rewritten. Every occurrence of ```https://www.asahi.com``` (the base URL) and ```asahi.com``` (the domain without ```www```) is respectively replaced by ```http://www.akahi.com``` and ```akahi.com``` in the response, while the replacement is inverse in the request.

### case 3

```http://mitmproxy-server/remove:script/nkf/https://www.itmedia.co.jp/```

This fetches https://www.itmedia.co.jp/, a ```Shift_JIS``` page, but you'd get a decoded ```UTF-8``` page after the occurrence of ```charset=Shift_JIS``` is found and every occurrence of ```<script>``` is removed.

### case 4

```http://mitmproxy-server/proxy:nkf/https://www.atmarkit.jp/https://atmarkit.itmedia.co.jp/```

This fetches https://atmarkit.itmedia.co.jp/, ensuring hyperlinks are rewritten and this Shift_JIS page is decoded to UTF-8.

## List of filters
Under construction
