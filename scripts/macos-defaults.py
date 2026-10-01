#!/usr/bin/env python3
"""Keep macOS `defaults` domains in this repository as readable YAML.

macOS preference plists cannot be symlinked - cfprefsd rewrites them by atomic
replace, which breaks the link - so each managed domain is mirrored into
data/macos/defaults/<domain>.yaml instead. YAML rather than the plist XML so
that a setting changed in a GUI shows up as a diff you can actually read.

Whole domains are mirrored, not a hand-picked set of keys: the point is to
notice changes made through a GUI, and a hand-picked list only ever notices
what someone already thought of.

Excluded keys are the exception. They are skipped on export and left alone on
import, so volatile junk (timestamps, window coordinates, analytics stamps)
and machine-specific values (display UUIDs) neither pollute the diff nor get
copied onto another machine.

    scripts/macos-defaults.py export   # live defaults -> YAML
    scripts/macos-defaults.py import   # YAML -> live defaults
    scripts/macos-defaults.py diff     # what changed since the last export

Round-tripping is lossless: plist data becomes YAML `!!binary`, plist dates
become YAML timestamps, and keys are sorted so a rewrite that only reorders a
domain is not a diff.
"""

import argparse
import difflib
import fnmatch
import os
import plistlib
import subprocess
import sys
import tempfile

try:
    import yaml
except ImportError:
    sys.exit("macos-defaults: PyYAML is required (pip install pyyaml)")

REPO = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
CONFIG = os.path.join(REPO, "data", "macos", "config.yaml")
OUT_DIR = os.path.join(REPO, "data", "macos", "defaults")


def load_config():
    with open(CONFIG, encoding="utf-8") as handle:
        config = yaml.safe_load(handle) or {}
    domains = config.get("domains") or []
    if not domains:
        sys.exit("macos-defaults: no domains configured in %s" % CONFIG)
    return domains, config.get("exclude") or {}


def excluded_patterns(exclude, domain):
    """Patterns that apply to this domain: its own, plus the global ones."""
    return list(exclude.get("*") or []) + list(exclude.get(domain) or [])


def is_excluded(key, patterns):
    return any(fnmatch.fnmatchcase(key, pattern) for pattern in patterns)


def read_domain(domain):
    """The live defaults for a domain, or None when it does not exist."""
    result = subprocess.run(
        ["defaults", "export", domain, "-"],
        capture_output=True,
        check=False,
    )
    if result.returncode != 0 or not result.stdout.strip():
        return None
    return plistlib.loads(result.stdout)


def write_domain(domain, data, drop=()):
    """Make a domain's defaults match `data`, through cfprefsd.

    `defaults import` merges rather than replaces, so a key that disappeared
    from the repository would survive it. Those are deleted by name instead,
    which is also why this never simply wipes the domain first: that would
    take the excluded keys with it.
    """
    for key in drop:
        subprocess.run(["defaults", "delete", domain, key], check=True)
    handle = tempfile.NamedTemporaryFile(suffix=".plist", delete=False)
    try:
        plistlib.dump(data, handle, sort_keys=True)
        handle.close()
        subprocess.run(["defaults", "import", domain, handle.name], check=True)
    finally:
        os.unlink(handle.name)


def to_yaml(data):
    return yaml.safe_dump(
        data,
        allow_unicode=True,
        default_flow_style=False,
        sort_keys=True,
        width=100,
    )


def yaml_path(domain):
    return os.path.join(OUT_DIR, domain + ".yaml")


def export_domain(domain, patterns):
    """Return the YAML a domain exports to, or None when it does not exist."""
    live = read_domain(domain)
    if live is None:
        return None
    return to_yaml({k: v for k, v in live.items() if not is_excluded(k, patterns)})


def cmd_export(domains, exclude):
    os.makedirs(OUT_DIR, exist_ok=True)
    changed = 0
    for domain in domains:
        patterns = excluded_patterns(exclude, domain)
        text = export_domain(domain, patterns)
        if text is None:
            print("  skip     %s (no such domain on this machine)" % domain)
            continue
        path = yaml_path(domain)
        before = None
        if os.path.exists(path):
            with open(path, encoding="utf-8") as handle:
                before = handle.read()
        if before == text:
            print("  ok       %s" % domain)
            continue
        with open(path, "w", encoding="utf-8") as handle:
            handle.write(text)
        print("  %-8s %s" % ("updated" if before is not None else "new", domain))
        changed += 1
    print("\n%d domain(s) changed." % changed)
    return 0


def cmd_import(domains, exclude):
    for domain in domains:
        path = yaml_path(domain)
        if not os.path.exists(path):
            print("  skip     %s (not exported yet)" % domain)
            continue
        with open(path, encoding="utf-8") as handle:
            wanted = yaml.safe_load(handle) or {}
        patterns = excluded_patterns(exclude, domain)
        live = read_domain(domain) or {}
        # Excluded keys are this machine's business: keep whatever is there.
        for key, value in live.items():
            if is_excluded(key, patterns):
                wanted[key] = value
        drop = [key for key in live if key not in wanted]
        if live == wanted:
            print("  ok       %s" % domain)
            continue
        write_domain(domain, wanted, drop)
        print("  applied  %s" % domain)
    print("\nRestart the affected apps (or log out) for every change to take effect.")
    return 0


def cmd_diff(domains, exclude):
    dirty = []
    for domain in domains:
        patterns = excluded_patterns(exclude, domain)
        live = export_domain(domain, patterns)
        if live is None:
            continue
        path = yaml_path(domain)
        stored = ""
        if os.path.exists(path):
            with open(path, encoding="utf-8") as handle:
                stored = handle.read()
        if stored == live:
            continue
        dirty.append(domain)
        sys.stdout.writelines(
            difflib.unified_diff(
                stored.splitlines(keepends=True),
                live.splitlines(keepends=True),
                fromfile="repository/%s.yaml" % domain,
                tofile="live/%s" % domain,
            )
        )
    if not dirty:
        print("Every managed domain matches the repository.")
        return 0
    print("\n%d domain(s) differ: %s" % (len(dirty), ", ".join(dirty)))
    print("Run `scripts/macos-defaults.py export` to record them.")
    return 1


def main():
    parser = argparse.ArgumentParser(description=__doc__.split("\n")[0])
    parser.add_argument("command", choices=["export", "import", "diff"])
    args = parser.parse_args()

    if sys.platform != "darwin":
        sys.exit("macos-defaults: macOS only")

    domains, exclude = load_config()
    return {
        "export": cmd_export,
        "import": cmd_import,
        "diff": cmd_diff,
    }[args.command](domains, exclude)


if __name__ == "__main__":
    sys.exit(main())
