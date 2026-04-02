const fs = require('fs');
const path = require('path');

function walkDir(dir, callback) {
  fs.readdirSync(dir).forEach(f => {
    let dirPath = path.join(dir, f);
    let isDirectory = fs.statSync(dirPath).isDirectory();
    isDirectory ? walkDir(dirPath, callback) : callback(dirPath);
  });
}

walkDir(path.join(__dirname, 'frontend', 'src', 'components'), function(filePath) {
  if (!filePath.endsWith('.tsx') && !filePath.endsWith('.ts')) return;
  
  let content = fs.readFileSync(filePath, 'utf8');

  if (content.includes('fetch(')) {
    // We only replace fetch calls that are likely API calls (await fetch)
    if (content.includes('await fetch(')) {
        content = content.replace(/await fetch\(/g, 'await apiFetch(');
        
        if (!content.includes('apiFetch')) {
            // Need to add import
            // Assuming it's in a subfolder, we can just use @/lib/api/base
            if (content.includes("from '@/lib/api/base'")) {
                 content = content.replace(/\{([^}]+)\}\s+from\s+'@\/lib\/api\/base'/, (match, p1) => `{ ${p1.trim()}, apiFetch } from '@/lib/api/base'`);
            } else {
                 content = `import { apiFetch } from '@/lib/api/base';\n` + content;
            }
        } else if (!content.includes("from '@/lib/api/base'") && !content.includes("from '../../lib/api/base'") && !content.includes("from '../lib/api/base'")) {
             content = `import { apiFetch } from '@/lib/api/base';\n` + content;
        }

        fs.writeFileSync(filePath, content, 'utf8');
        console.log(`Updated ${filePath}`);
    }
  }
});
