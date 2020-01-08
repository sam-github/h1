Lifecycle:

- NEW: reported, initial state

- (optional) NEEDS-MORE-INFO: back-and-forth with reporter
- ... evaluate, then go to "TRIAGED" or one of the "not-fixing" states

- TRIAGED: agree it is a (non-duplicate, resolvable) issue
- ... fix it
- ... no state for an issue that is fixed and waiting on release, or
  for one that is being worked on. There is an assignee, could be used to
  track something that someone has agreed to sheppard towards resolution.
- ... there is a `resolve` tag on some H1 issues, I'm not exactly sure what it
  means.

Do all the following states count as "closed"?

Final state for a report we will fix:
- RESOLVED: fix published

Final states for a report we will not fix:
- NOT-APPLICABLE: this means we don't agree its a vulnerability?
- INFORMATIVE: this means we agree there is something problematic, but we
  can't or won't (out of support? unfixable?) publish a fix?
- DUPLICATE: obvious... but should it get disclosed once the existing report
  is disclosed?
- SPAM

Is this pseudo-state the penultimate state? Should everything get disclosed,
or just RESOLVED?
- "disclosed": ... there is no state for disclosed, but there is a "disclosed" date
