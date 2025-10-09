#!/usr/bin/env node
const { execFileSync } = require('node:child_process');

function run(cmd, args, opts={}) {
  try { return execFileSync(cmd, args, { stdio: 'inherit', ...opts }); }
  catch (e) { process.exit(e.status || 1); }
}

function main() {
  const [feature, taskId] = process.argv.slice(2);
  if (!feature || !taskId) {
    console.error('Usage: node scripts/todo/branch.js <feature-slug> <TASK_ID>');
    console.error('Example: node scripts/todo/branch.js validators T070');
    process.exit(1);
  }
  const branch = `feat/${feature}-task-${taskId}`;
  run('git', ['fetch', 'origin']);
  // refuse to overwrite existing branch
  try {
    run('git', ['rev-parse', '--verify', branch], { stdio: 'pipe' });
    console.error(`Branch ${branch} already exists. Aborting.`);
    process.exit(1);
  } catch (_) {
    // rev-parse failed -> branch does not exist; proceed
  }
  run('git', ['checkout', '-b', branch, 'origin/main']);
  console.log(`Switched to branch ${branch}`);
}

main();
