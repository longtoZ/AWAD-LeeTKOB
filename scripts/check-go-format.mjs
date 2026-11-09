import { execFileSync } from 'node:child_process';

const files = execFileSync('git', ['ls-files', '--cached', '--others', '--exclude-standard', '*.go'], {
  encoding: 'utf8',
}).trim();

if (!files) process.exit(0);

const result = execFileSync('gofmt', ['-l', ...files.split(/\r?\n/)], {
  encoding: 'utf8',
});

if (result.trim()) {
  console.error(`Go files need formatting:\n${result}`);
  process.exit(1);
}
