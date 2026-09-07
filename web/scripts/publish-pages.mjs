import { execFileSync } from 'node:child_process';
import { mkdtempSync, readFileSync, rmSync, writeFileSync } from 'node:fs';
import { tmpdir } from 'node:os';
import { dirname, join, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';

const webRoot = resolve(dirname(fileURLToPath(import.meta.url)), '..');
const repoRoot = resolve(webRoot, '..');
const htmlPath = resolve(repoRoot, 'docs/index.html');

function gh(args, options = {}) {
	return execFileSync('gh', args, {
		cwd: repoRoot,
		encoding: 'utf8',
		...options
	});
}

const html = readFileSync(htmlPath);
const payload = {
	message: 'Publish quiz UI to GitHub Pages',
	content: html.toString('base64'),
	branch: 'main'
};

try {
	payload.sha = gh([
		'api',
		'repos/{owner}/{repo}/contents/docs/index.html',
		'--jq',
		'.sha'
	]).trim();
} catch {
	// docs/index.html がまだ無い初回は sha なしで作成する
}

const dir = mkdtempSync(join(tmpdir(), 'publish-pages-'));
const payloadPath = join(dir, 'payload.json');
writeFileSync(payloadPath, JSON.stringify(payload));

try {
	const result = JSON.parse(
		gh([
			'api',
			'--method',
			'PUT',
			'repos/{owner}/{repo}/contents/docs/index.html',
			'--input',
			payloadPath
		])
	);
	console.log(result.commit?.html_url ?? 'updated docs/index.html');
	console.log('https://haruyama480.github.io/learning-engineer-english/');
} finally {
	rmSync(dir, { recursive: true, force: true });
}
