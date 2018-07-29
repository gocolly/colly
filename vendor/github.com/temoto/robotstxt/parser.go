package robotstxt

// Comments explaining the logic are taken from either the google's spec:
// https://developers.google.com/webmasters/control-crawl-index/docs/robots_txt
//
// or the Wikipedia's entry on robots.txt:
// http://en.wikipedia.org/wiki/Robots.txt

import (
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type lineType uint

const (
	lIgnore lineType = iota
	lUnknown
	lUserAgent
	lAllow
	lDisallow
	lCrawlDelay
	lSitemap
	lHost
)

type parser struct {
	tokens []string
	pos    int
}

type lineInfo struct {
	t  lineType       // Type of line key
	k  string         // String representation of the type of key
	vs string         // String value of the key
	vf float64        // Float value of the key
	vr *regexp.Regexp // Regexp value of the key
}

func newParser(tokens []string) *parser {
	return &parser{tokens: tokens}
}

func parseGroupMap(groups map[string]*Group, agents []string, fun func(*Group)) {
	var g *Group
	for _, a := range agents {
		if g = groups[a]; g == nil {
			g = new(Group)
			groups[a] = g
		}
		fun(g)
	}
}

func (p *parser) parseAll() (groups map[string]*Group, host string, sitemaps []string, errs []error) {
	groups = make(map[string]*Group, 16)
	agents := make([]string, 0, 4)
	isEmptyGroup := true

	// Reset internal fields, tokens are assigned at creation time, never change
	p.pos = 0

	for {
		if li, err := p.parseLine(); err != nil {
			if err == io.EOF {
				break
			}
			errs = append(errs, err)
		} else {
			switch li.t {
			case lUserAgent:
				// Two successive user-agent lines are part of the same group.
				if !isEmptyGroup {
					// End previous group
					agents = make([]string, 0, 4)
				}
				if len(agents) == 0 {
					isEmptyGroup = true
				}
				agents = append(agents, li.vs)

			case lDisallow:
				// Error if no current group
				if len(agents) == 0 {
					errs = append(errs, fmt.Errorf("Disallow before User-agent at token #%d.", p.pos))
				} else {
					isEmptyGroup = false
					var r *rule
					if li.vr != nil {
						r = &rule{"", false, li.vr}
					} else {
						r = &rule{li.vs, false, nil}
					}
					parseGroupMap(groups, agents, func(g *Group) { g.rules = append(g.rules, r) })
				}

			case lAllow:
				// Error if no current group
				if len(agents) == 0 {
					errs = append(errs, fmt.Errorf("Allow before User-agent at token #%d.", p.pos))
				} else {
					isEmptyGroup = false
					var r *rule
					if li.vr != nil {
						r = &rule{"", true, li.vr}
					} else {
						r = &rule{li.vs, true, nil}
					}
					parseGroupMap(groups, agents, func(g *Group) { g.rules = append(g.rules, r) })
				}

			case lHost:
				host = li.vs

			case lSitemap:
				sitemaps = append(sitemaps, li.vs)

			case lCrawlDelay:
				if len(agents) == 0 {
					errs = append(errs, fmt.Errorf("Crawl-delay before User-agent at token #%d.", p.pos))
				} else {
					isEmptyGroup = false
					delay := time.Duration(li.vf * float64(time.Second))
					parseGroupMap(groups, agents, func(g *Group) { g.CrawlDelay = delay })
				}
			}
		}
	}
	return
}

func (p *parser) parseLine() (li *lineInfo, err error) {
	t1, ok1 := p.popToken()
	if !ok1 {
		// proper EOF
		return nil, io.EOF
	}

	t2, ok2 := p.peekToken()
	if !ok2 {
		// EOF, no value associated with the token, so ignore token and return
		return nil, io.EOF
	}

	// Helper closure for all string-based tokens, common behaviour:
	// - Consume t2 token
	// - If empty, return unkown line info
	// - Otherwise return the specified line info
	returnStringVal := func(t lineType) (*lineInfo, error) {
		p.popToken()
		if t2 != "" {
			return &lineInfo{t: t, k: t1, vs: t2}, nil
		}
		return &lineInfo{t: lIgnore}, nil
	}

	// Helper closure for all path tokens (allow/disallow), common behaviour:
	// - Consume t2 token
	// - If empty, return unkown line info
	// - Otherwise, normalize the path (add leading "/" if missing, remove trailing "*")
	// - Detect if wildcards are present, if so, compile into a regexp
	// - Return the specified line info
	returnPathVal := func(t lineType) (*lineInfo, error) {
		p.popToken()
		if t2 != "" {
			if !strings.HasPrefix(t2, "*") && !strings.HasPrefix(t2, "/") {
				t2 = "/" + t2
			}
			if strings.HasSuffix(t2, "*") {
				t2 = strings.TrimRight(t2, "*")
			}
			// From google's spec:
			// Google, Bing, Yahoo, and Ask support a limited form of
			// "wildcards" for path values. These are:
			//   * designates 0 or more instances of any valid character
			//   $ designates the end of the URL
			if strings.ContainsAny(t2, "*$") {
				// Must compile a regexp, this is a pattern.
				// Escape string before compile.
				t2 = regexp.QuoteMeta(t2)
				t2 = strings.Replace(t2, `\*`, `.*`, -1)
				t2 = strings.Replace(t2, `\$`, `$`, -1)
				if r, e := regexp.Compile(t2); e != nil {
					return nil, e
				} else {
					return &lineInfo{t: t, k: t1, vr: r}, nil
				}
			} else {
				// Simple string path
				return &lineInfo{t: t, k: t1, vs: t2}, nil
			}
		}
		return &lineInfo{t: lIgnore}, nil
	}

	switch strings.ToLower(t1) {
	case "\n":
		// Don't consume t2 and continue parsing
		return &lineInfo{t: lIgnore}, nil

	case "user-agent", "useragent":
		// From google's spec:
		// Handling of <field> elements with simple errors / typos (eg "useragent"
		// instead of "user-agent") is undefined and may be interpreted as correct
		// directives by some user-agents.
		// The user-agent is non-case-sensitive.
		t2 = strings.ToLower(t2)
		return returnStringVal(lUserAgent)

	case "disallow":
		// From google's spec:
		// When no path is specified, the directive is ignored (so an empty Disallow
		// CAN be an allow, since allow is the default. The actual result depends
		// on the other rules in the group).
		return returnPathVal(lDisallow)

	case "allow":
		// From google's spec:
		// When no path is specified, the directive is ignored.
		return returnPathVal(lAllow)

	case "host":
		// Host directive to specify main site mirror
		// Read more: https://help.yandex.com/webmaster/controlling-robot/robots-txt.xml#host
		return returnStringVal(lHost)

	case "sitemap":
		// Non-group field, applies to the host as a whole, not to a specific user-agent
		return returnStringVal(lSitemap)

	case "crawl-delay", "crawldelay":
		// From http://en.wikipedia.org/wiki/Robots_exclusion_standard#Nonstandard_extensions
		// Several major crawlers support a Crawl-delay parameter, set to the
		// number of seconds to wait between successive requests to the same server.
		p.popToken()
		if cd, e := strconv.ParseFloat(t2, 64); e != nil {
			return nil, e
		} else {
			return &lineInfo{t: lCrawlDelay, k: t1, vf: cd}, nil
		}
	}

	// Consume t2 token
	p.popToken()
	return &lineInfo{t: lUnknown, k: t1}, nil
}

func (p *parser) popToken() (tok string, ok bool) {
	tok, ok = p.peekToken()
	if !ok {
		return
	}
	p.pos++
	return tok, true
}

func (p *parser) peekToken() (tok string, ok bool) {
	if p.pos >= len(p.tokens) {
		return "", false
	}
	return p.tokens[p.pos], true
}
