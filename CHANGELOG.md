# 2.0.0 - 2019.11.28

 - Breaking change: Change Collector.RedirectHandler member to Collector.SetRedirectHandler function
 - Go module support
 - Collector.HasVisited method added to be able to check if an url has been visited
 - Collector.SetClient method introduced
 - HTMLElement.ChildTexts method added
 - New user agents
 - Multiple bugfixes

# 1.2.0 - 2019.02.13

 - Compatibility with the latest htmlquery package
 - New request shortcut for HEAD requests
 - Check URL availibility before visiting
 - Fix proxy URL value
 - Request counter fix
 - Minor fixes in examples

# 1.1.0 - 2018.08.13

 - Appengine integration takes context.Context instead of http.Request (API change)
 - Added "Accept" http header by default to every request
 - Support slices of pointers in unmarshal
 - Fixed a race condition in queues
 - ForEachWithBreak method added to HTMLElement
 - Added a local file example
 - Support gzip decompression of response bodies
 - Don't share waitgroup when cloning a collector
 - Fixed instagram example


# 1.0.0 - 2018.05.13
