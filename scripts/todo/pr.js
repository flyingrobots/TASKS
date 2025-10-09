#!/usr/bin/env node
const { execFileSync } = require('node:child_process');

function run(cmd, args, opts={}) {
  try { return execFileSync(cmd, args, { stdio: 'pipe', encoding: 'utf8', ...opts }).trim(); }
  catch (e) { const out = (e.stdout||'') + (e.stderr||''); throw new Error(`${cmd} ${args.join(' ')}\n${out}`); }
}

function main() {
  try {
    const branch = run('git', ['rev-parse', '--abbrev-ref', 'HEAD']);
    if (!branch) throw new Error('Unable to determine current branch');
    // Ensure we have gh
    run('gh', ['--version']);
    const out = run('gh', ['pr', 'create', '--fill', '--base', 'main', '--head', branch]);
    console.log(out);
  } catch (err) {
    console.error('Failed to create PR via gh CLI.');
    console.error(String(err.message||err));
    process.exit(1);
  }
}

main();

