const fs = require('fs');
const path = require('path');

const dir = path.join(__dirname, 'frontend', 'src', 'lib', 'api');
const files = fs.readdirSync(dir).filter(f => f.endsWith('.ts') && f !== 'base.ts' && f !== 'index.ts');

for (const file of files) {
  const filePath = path.join(dir, file);
  let content = fs.readFileSync(filePath, 'utf8');

  // Skip if it doesn't use fetch
  if (!content.includes('fetch(')) {
    continue;
  }

  // Replace fetch( with apiFetch(
  content = content.replace(/await fetch\(/g, 'await apiFetch(');

  // Add apiFetch to the import from './base'
  if (content.includes("from './base';")) {
    if (!content.includes('apiFetch')) {
      content = content.replace(/\{([^}]+)\}\s+from\s+'\.\/base';/, (match, p1) => {
        return `{ ${p1.trim()}, apiFetch } from './base';`;
      });
    }
  } else {
    // If not importing from base, add it
    content = `import { apiFetch } from './base';\n` + content;
  }

  fs.writeFileSync(filePath, content, 'utf8');
  console.log(`Updated ${file}`);
}
