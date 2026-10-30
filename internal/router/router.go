package router

import "strings"

// RewriteURI mirrors the CloudFront Function rewrite logic for Next.js static exports.
func RewriteURI(uri string) string {
	if strings.HasSuffix(uri, "/") {
		return uri + "index.html"
	}
	if !strings.Contains(uri, ".") {
		return uri + ".html"
	}
	return uri
}
