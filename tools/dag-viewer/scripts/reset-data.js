#!/usr/bin/env node

import { unlinkSync, existsSync } from 'fs';
import { dirname, join } from 'path';
import { fileURLToPath } from 'url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

// Files to delete for reset
const STATE_FILES = [
  'dag-state-full.json',
  'dag-state.json'
];

console.log('🔄 Resetting DAG viewer data...\n');

// Delete state files
STATE_FILES.forEach(file => {
  const fullPath = join(dirname(__dirname), file);
  if (existsSync(fullPath)) {
    try {
      unlinkSync(fullPath);
      console.log(`✅ Deleted: ${file}`);
    } catch (error) {
      console.error(`❌ Failed to delete ${file}:`, error.message);
    }
  } else {
    console.log(`⏭️  Skipped (not found): ${file}`);
  }
});

// Reset by sending empty state to server if it's running
try {
  const response = await fetch('http://localhost:8080/reset', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ confirm: true })
  }).catch(() => null);
  
  if (response && response.ok) {
    console.log('\n✅ Server state reset successfully');
  } else {
    console.log('\n⚠️  Server not running or /reset endpoint not available');
  }
} catch (error) {
  console.log('\n⚠️  Could not connect to server (this is normal if server is not running)');
}

console.log('\n🎉 Data reset complete!');
console.log('→ Start the server with: npm run server:full');
console.log('→ Start the React app with: npm run dev:react');