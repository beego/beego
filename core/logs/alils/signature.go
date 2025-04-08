package alils

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"
)

// GMT location
var gmtLoc = time.FixedZone("GMT", 0)

// NowRFC1123 returns now time in RFC1123 format with GMT timezone,
// eg. "Mon, 02 Jan 2006 15:04:05 GMT".
func nowRFC1123() string {
	return time.Now().In(gmtLoc).Format(time.RFC1123)
}

// signature calculates a request's signature digest
func signature(project *LogProject, method, uri string,
	headers map[string]string) (digest string, err error) {

	// 1. Extract basic header information
	contentMD5 := getHeaderSafe(headers, "Content-MD5")
	contentType := getHeaderSafe(headers, "Content-Type")

	date, ok := headers["Date"]
	if !ok {
		return "", fmt.Errorf("Can't find 'Date' header")
	}

	// 2. Calculate CanonicalizedSLSHeaders
	canoHeaders := buildCanonicalHeaders(headers)

	// 3. Calculate CanonicalizedResource
	canoResource, err := buildCanonicalResource(uri)
	if err != nil {
		return "", err
	}

	// 4. Build the signature string
	signStr := buildSignString(method, contentMD5, contentType, date, canoHeaders, canoResource)

	// 5. Calculate HMAC-SHA1 signature and encode with Base64
	// Signature = base64(hmac-sha1(UTF8-Encoding-Of(SignString)，AccessKeySecret))
	return calculateHmacSha1(signStr, project.AccessKeySecret)
}

// getHeaderSafe safely retrieves a header value, returns empty string if not exists
func getHeaderSafe(headers map[string]string, key string) string {
	if val, ok := headers[key]; ok {
		return val
	}
	return ""
}

// buildCanonicalHeaders constructs normalized SLS headers
func buildCanonicalHeaders(headers map[string]string) string {
	slsHeaders := make(map[string]string)
	var slsHeaderKeys sort.StringSlice

	// Extract headers prefixed with x-log- and x-acs-
	for k, v := range headers {
		l := strings.TrimSpace(strings.ToLower(k))
		if strings.HasPrefix(l, "x-log-") || strings.HasPrefix(l, "x-acs-") {
			slsHeaders[l] = strings.TrimSpace(v)
			slsHeaderKeys = append(slsHeaderKeys, l)
		}
	}

	// Sort headers alphabetically
	sort.Sort(slsHeaderKeys)

	// Build the canonical header string
	var result strings.Builder
	for i, k := range slsHeaderKeys {
		result.WriteString(k)
		result.WriteString(":")
		result.WriteString(slsHeaders[k])
		if i+1 < len(slsHeaderKeys) {
			result.WriteString("\n")
		}
	}

	return result.String()
}

// buildCanonicalResource constructs a normalized resource string
func buildCanonicalResource(uri string) (string, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return "", err
	}

	result := u.Path

	// Process query parameters
	if u.RawQuery != "" {
		var keys sort.StringSlice
		vals := u.Query()

		for k := range vals {
			keys = append(keys, k)
		}

		sort.Sort(keys)
		result += "?"

		for i, k := range keys {
			if i > 0 {
				result += "&"
			}

			for _, v := range vals[k] {
				result += k + "=" + v
			}
		}
	}

	return result, nil
}

// buildSignString constructs the complete signature string
func buildSignString(method, contentMD5, contentType, date, canoHeaders, canoResource string) string {
	// SignString = VERB + "\n"
	//              + CONTENT-MD5 + "\n"
	//              + CONTENT-TYPE + "\n"
	//              + DATE + "\n"
	//              + CanonicalizedSLSHeaders + "\n"
	//              + CanonicalizedResource
	return method + "\n" +
		contentMD5 + "\n" +
		contentType + "\n" +
		date + "\n" +
		canoHeaders + "\n" +
		canoResource
}

// calculateHmacSha1 calculates the HMAC-SHA1 signature and encodes it with Base64
func calculateHmacSha1(signStr, secret string) (string, error) {
	mac := hmac.New(sha1.New, []byte(secret))
	_, err := mac.Write([]byte(signStr))
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(mac.Sum(nil)), nil
}
