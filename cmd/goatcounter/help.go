// Copyright © 2019 Martin Tournoij – This file is part of GoatCounter and
// published under the terms of a slightly modified EUPL v1.2 license, which can
// be found in the LICENSE file or at https://license.goatcounter.com

package main

import (
	"fmt"
	"os"
	"strings"
	"unicode/utf8"

	"zgo.at/errors"
	"zgo.at/zli"
)

func printHelp(t string) {
	fmt.Fprint(stdout, zli.Usage(zli.UsageTrim|zli.UsageHeaders, t))
}

func help() (int, error) {
	zli.WantColor = true

	if len(os.Args) == 2 {
		printHelp(usage[""])
		return 0, nil
	}

	if os.Args[2] == "all" {
		printHelp(usage[""])
		fmt.Println()
		for _, h := range []string{
			"help", "version",
			"migrate", "create", "serve",
			"reindex", "monitor",
			"db", "listen",
		} {
			head := fmt.Sprintf("─── Help for %q ", h)
			fmt.Fprintf(stdout, "%s%s\n\n",
				zli.Colorf(head, zli.Bold),
				strings.Repeat("─", 80-utf8.RuneCountInString(head)))
			printHelp(usage[h])
			fmt.Println()
		}
		return 0, nil
	}

	t, ok := usage[os.Args[2]]
	if !ok {
		return 1, errors.Errorf("no help topic for %q", os.Args[2])
	}
	printHelp(t)
	return 0, nil
}

const usageHelp = `
Show help; use "help commands" to dispay detailed help for a command, or "help
all" to display everything.
`

const helpListen = `
GoatCounter is designed to work well "out-of-the-box" for most people, but there
are some complexities surrounding the ACME/Let's Encrypt certificate creation
that can make things a bit complex.

In order for Let's Encrypt to work GoatCounter *needs* to be publicly accessible
on port 80 because Let's Encrypt must verify that you actually own the domain by
accessing http://example.com/.well-known/acme-challenge/[secret]; GoatCounter
handles all of this for you, but it does need to be reachable by Let's Encrypt's
verification server.

This is why GoatCounter listens on port 80 by default, which should work well
for most people.

listen and tls flags:

    You can change the main port GoatCounter listens on with the -listen flag.
    This works like most applications, for example:

        -listen localhost:8081     Listen on localhost:8081
        -listen :8081              Listen on :8081 for all addresses

    The -tls flag controls the TLS setup, as well as redirecting port 80 the
    -listen port with a 301 status code. Because there are a few different
    server setups GoatCounter can be used in, the flag accepts a bunch of
    different options as a comma-separated list with any combination of:

        none        Don't serve any TLS; this is the same as just leaving it
                    blank (-tls '') but more explicit/readable.

        tls         Accept TLS connections on -listen; if this isn't added it
                    will accept regular non-https connections, but may still
                    generate certificates with ACME (useful for proxy or dev).

        rdr         Redirect port 80 to the -listen port; as mentioned it's
                    required for Let's Encrypt to be available on port 80. You
                    can also use a proxy in front of GoatCounter (documented in
                    more detail below).

        acme[:dir]  Create TLS certificates with ACME.

                    This can optionally followed by a : and a cache directory
                    path (default: ./acme-secrets) where the certificates and
                    your account key will be stored. The directory will be
                    created if it doesn't exist yet. As indicated by the name,
                    the contents of this directory should be kept *secret*.

        file.pem    TLS certificate and keyfile, in one file. This can be used
                    as an alternative to Let's Encrypt if you already have a
                    certificate from your domain from a CA. This can use used
                    multiple times (e.g. "-tls tls,file1.pem,file2.pem").

                    This can also be combined with the acme option: GoatCounter
                    will try to use a certificate file for the domain first, and
                    if this doesn't exist it will try to create a certificate
                    with ACME.

    Some common examples:

        -tls tls,acme,rdr
            This is the default setting.

        -tls tls,rdr,acme:/home/gc/.acme
            The default setting, but with a different cache directory.

        -tls tls,/etc/tls/stats.example.com.pem
            Don't use ACME, but use a certificate from a CA. No port 80 redirect.

Proxy Setup:

    If you want to serve GoatCounter behind a proxy (HAproxy, Varnish, Hitch,
    nginx, Caddy, whatnot) then you'll want to use something like:

        goatcounter serve -listen localhost:8081 -tls none

    And then forward requests on port 80 and 443 for your domain to
    localhost:8081. This assumes that the proxy will take care of the TLS
    certificate story.

    You can still use GoatCounter's ACME if you want:

        goatcounter serve -listen localhost:8081 -tls tls,acme

    You will have to make the proxy reads the *.pem files from the acme cache
    directory. You may have to reload or restart the proxy for it to pick up new
    files.

    NOTE: the certificates have a short expiry of a few months and will be
    regenerated automatically. This means that the proxy will have to be
    notified of this, most accept a signal to reload the config and/or
    certificates. You probably want to do this automatically in a cron job or
    some such. Be sure to check this otherwise you'll run in to "certificate
    expired" errors a few months down the line.

    NOTE 2: this directory also contains a "acme_account+key" file; you don't
    want to read "*" but "*.pem" (some proxies ignore invalid certificates, for
    others it's a fatal error).

Using a non-standard port:

    If you make GoatCounter publicly accessibly on non-standard port (i.e. not
    80 or 443) then you must add the -port flag to tell GoatCounter which port
    to use in various redirects, messages, and emails:

        goatcounter serve -listen :9000 -port 9000

    This may seem redundant, but it's hard for GoatCounter to tell if it's
    accessible on :9000 or if there's a proxy in front of it redirecting :80 and
    :443 to :9000. Since most people will use the standard ports you need to
    explicitly tell GoatCounter to use a non-standard port.
`
