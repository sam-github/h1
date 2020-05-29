# h1 - report what issues are waiting for what actions

This attempts to report where an issue is in the Node.js security triage
[work flow][].

[work flow]: https://github.com/nodejs/TSC/blob/master/Security-Team.md#security-triage-workflow

It requires creating a HackerOne API token, see:
- https://hackerone.com/nodejs/api

Each token requires an identifier.  The identifier can be any string. For
example, if your Github ID is `spiffy-cat`, you could use `at-spiffy-cat` so
that people can tell who allocated the token when browsing the currently
allocated API tokens.

Create a `.token` file in your CWD. If you have only one API token for a single
program, such as `nodejs`, it can just be:
```text
your-identifier : your-token
```

If you have API tokens for multiple programs, you can specify which program
an identifier relates to, for example:
```text
your-nodejs-identifer@nodejs:your-token-for-nodejs
your-nodejs-ecosystem-identifer@nodejs-ecosystem:your-token-for-ecosystem
```

Despite its multi-program support, its assumptions are tied to the Node.js
[work flow][].

Usage:  `./bin/h1 -h`

Building: `make build`

Generating the report: `make day`

TODO:
- Rewrite as Node.js? :-)
