#!/usr/bin/env python3
"""Exercise the installer offline with a fake release download service."""
import hashlib
import os
from pathlib import Path
import subprocess
import tempfile

installer = Path(__file__).resolve().parents[1] / 'install.sh'
with tempfile.TemporaryDirectory() as directory:
    root = Path(directory)
    tools = root / 'tools'
    tools.mkdir()
    payload = root / 'binary'
    payload.write_text('#!/bin/sh\necho fixture\n')
    manifest = root / 'checksums.txt'
    digest = hashlib.sha256(payload.read_bytes()).hexdigest()
    manifest.write_text(''.join(f'{digest}  alt-tab_{os_name}_{arch}\n'
                              for os_name in ('linux', 'darwin')
                              for arch in ('amd64', 'arm64')))
    (tools / 'curl').write_text('''#!/bin/sh
while [ "$#" -gt 0 ]; do
 case "$1" in
  https://*) url=$1; shift ;;
  -o) dest=$2; shift 2 ;;
  --retry|--proto) shift 2 ;;
  *) shift ;;
 esac
done
case "$url" in
 */checksums.txt) cp "$FIXTURE/checksums.txt" "$dest" ;;
 *) cp "$FIXTURE/binary" "$dest" ;;
esac
''')
    (tools / 'uname').write_text('''#!/bin/sh
case "$1" in -s) echo "$TEST_OS" ;; -m) echo "$TEST_ARCH" ;; esac
''')
    for tool in tools.iterdir():
        tool.chmod(0o755)
    env = dict(os.environ, PATH=f'{tools}:{os.environ["PATH"]}', FIXTURE=str(root),
               HOME=str(root), TEST_OS='Linux', TEST_ARCH='x86_64')

    def run(*args, success=True):
        result = subprocess.run(['sh', str(installer), *args], env=env,
                                capture_output=True, text=True)
        assert (result.returncode == 0) == success, result.stdout + result.stderr
        return result

    run('--help')
    run('--bogus', success=False)
    run('--dir', success=False)
    run('--version', '../bad', success=False)
    for system, arch in [('Linux', 'x86_64'), ('Linux', 'aarch64'),
                         ('Darwin', 'x86_64'), ('Darwin', 'arm64')]:
        env.update(TEST_OS=system, TEST_ARCH=arch)
        target = root / f'{system} {arch}'
        run('--dir', str(target), '--version', 'v0.1.0')
        assert (target / 'alt-tab').read_bytes() == payload.read_bytes()
        assert os.access(target / 'alt-tab', os.X_OK)
    target = root / '.local/bin'
    run()
    installed = (target / 'alt-tab').read_bytes()
    payload.write_text('tampered')
    run(success=False)
    assert (target / 'alt-tab').read_bytes() == installed
    assert not list(target.glob('.alt-tab.*'))
    manifest.write_text('')
    run(success=False)
    env['TEST_OS'] = 'Windows'
    run(success=False)
    env.update(TEST_OS='Linux', TEST_ARCH='riscv64')
    run(success=False)
print('Installer checks passed: four platforms, options, checksums, atomic replacement, unsupported platforms.')
