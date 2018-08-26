def asciidoc(name, src, attrs={}, confs=[], data=[], tools=[], tags=None):
  args = ['--backend', 'html', '--no-header-footer']
  for key, val in attrs.items():
    if val:
      args += ["--attribute=%s=%s" % (key, val)]
    else:
      args += ["--attribute=%s!" % key]
  args += ["--conf-file=$(location %s)" % c for c in confs]
  out = name + '.html'
  native.genrule(
    name = name,
    srcs = [src] + data + confs,
    outs = [out],
    local = True,
    tags = tags,
    tools = tools,
    output_to_bindir = 1,
    cmd = '\n'.join([
        'export OUTDIR="$$PWD/$(@D)"',
        'export LOGFILE="$$(mktemp -t \"XXXXXXasciidoc\")"',
        'export PATH="$$PATH:/usr/local/bin/"',
        'trap "rm -f \"$${LOGFILE}\"" EXIT ERR INT',
        "asciidoc %s -o $(@) $(location %s) 2> \"$${LOGFILE}\"" % (' '.join(args), src),
        'cat $${LOGFILE}',
        '! grep -q -e "filter non-zero exit code" -e "no output from filter" "$${LOGFILE}"']))

