package handler

const indexHTML = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>Service Registry</title>
<style>
  :root {
    --bg: #f5f5f5; --card: #fff; --border: #e0e0e0;
    --text: #1a1a1a; --muted: #666; --accent: #2563eb;
    --green: #16a34a; --red: #dc2626; --yellow: #ca8a04;
    --green-bg: #dcfce7; --red-bg: #fee2e2; --yellow-bg: #fef9c3;
  }
  * { box-sizing: border-box; margin: 0; padding: 0; }
  body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; background: var(--bg); color: var(--text); padding: 2rem; }
  h1 { font-size: 1.5rem; font-weight: 600; margin-bottom: .25rem; }
  .subtitle { color: var(--muted); font-size: .875rem; margin-bottom: 1.5rem; }
  .toolbar { display: flex; gap: .75rem; margin-bottom: 1.5rem; align-items: center; flex-wrap: wrap; }
  .toolbar button { padding: .5rem 1rem; border: 1px solid var(--border); border-radius: .375rem; background: var(--card); cursor: pointer; font-size: .875rem; transition: all .15s; }
  .toolbar button:hover { border-color: var(--accent); color: var(--accent); }
  .toolbar button.primary { background: var(--accent); color: #fff; border-color: var(--accent); }
  .toolbar button.primary:hover { opacity: .85; }
  .toolbar button.copied { background: var(--green); color: #fff; border-color: var(--green); }
  .stats { margin-left: auto; font-size: .8rem; color: var(--muted); }
  table { width: 100%; border-collapse: collapse; background: var(--card); border-radius: .5rem; overflow: hidden; box-shadow: 0 1px 3px rgba(0,0,0,.06); }
  th { text-align: left; padding: .75rem 1rem; font-size: .75rem; text-transform: uppercase; letter-spacing: .05em; color: var(--muted); background: #fafafa; border-bottom: 1px solid var(--border); }
  td { padding: .75rem 1rem; border-bottom: 1px solid var(--border); font-size: .875rem; vertical-align: middle; }
  tr:last-child td { border-bottom: none; }
  tr:hover td { background: #f9fafb; }
  .badge { display: inline-block; padding: .125rem .5rem; border-radius: 9999px; font-size: .75rem; font-weight: 500; }
  .badge-healthy { background: var(--green-bg); color: var(--green); }
  .badge-unhealthy { background: var(--red-bg); color: var(--red); }
  .badge-unknown { background: var(--yellow-bg); color: var(--yellow); }
  .svc-name { font-weight: 600; }
  .svc-desc { color: var(--muted); font-size: .8rem; margin-top: .125rem; }
  .svc-url a { color: var(--accent); text-decoration: none; font-size: .8rem; word-break: break-all; }
  .svc-url a:hover { text-decoration: underline; }
  .svc-ext { margin-top: .25rem; }
  .svc-ext a { color: var(--green); font-size: .75rem; }
  .time { font-size: .8rem; color: var(--muted); }
  .btn-del { background: none; border: none; cursor: pointer; color: var(--muted); font-size: .8rem; padding: .25rem .5rem; border-radius: .25rem; }
  .btn-del:hover { background: var(--red-bg); color: var(--red); }
  .empty { text-align: center; padding: 3rem 1rem; color: var(--muted); }
  .empty p { margin-bottom: .5rem; }

  /* Modal */
  .modal-overlay { display: none; position: fixed; inset: 0; background: rgba(0,0,0,.35); z-index: 100; align-items: center; justify-content: center; }
  .modal-overlay.active { display: flex; }
  .modal { background: var(--card); border-radius: .5rem; padding: 1.5rem; width: 100%; max-width: 440px; box-shadow: 0 8px 30px rgba(0,0,0,.15); }
  .modal h2 { font-size: 1.1rem; margin-bottom: 1rem; }
  .modal label { display: block; font-size: .8rem; font-weight: 500; margin-bottom: .25rem; color: var(--muted); }
  .modal input, .modal textarea { width: 100%; padding: .5rem .75rem; border: 1px solid var(--border); border-radius: .375rem; font-size: .875rem; margin-bottom: .75rem; font-family: inherit; }
  .modal textarea { resize: vertical; min-height: 60px; }
  .modal input:focus, .modal textarea:focus { outline: none; border-color: var(--accent); box-shadow: 0 0 0 2px rgba(37,99,235,.15); }
  .modal-actions { display: flex; justify-content: flex-end; gap: .5rem; margin-top: .5rem; }
  .modal-actions button { padding: .5rem 1rem; border-radius: .375rem; border: 1px solid var(--border); cursor: pointer; font-size: .875rem; background: var(--card); }
  .modal-actions .btn-save { background: var(--accent); color: #fff; border-color: var(--accent); }
</style>
</head>
<body>

<h1>Service Registry</h1>
<p class="subtitle">Local service discovery & health monitoring</p>

<div class="toolbar">
  <button class="primary" onclick="openModal()">+ Register Service</button>
  <button onclick="loadServices()">Refresh</button>
  <button id="btn-prompt" onclick="copyPrompt()">Copy AI Prompt</button>
  <span class="stats" id="stats"></span>
</div>

<table>
  <thead>
    <tr>
      <th>Service</th>
      <th>URL</th>
      <th>Status</th>
      <th>Last Checked</th>
      <th></th>
    </tr>
  </thead>
  <tbody id="tbody"></tbody>
</table>

<div class="modal-overlay" id="modal">
  <div class="modal">
    <h2>Register a Service</h2>
    <label for="f-name">Name</label>
    <input id="f-name" placeholder="e.g. my-api" autofocus>
    <label for="f-url">URL</label>
    <input id="f-url" placeholder="e.g. http://localhost:3000/health">
    <label for="f-desc">Description (optional)</label>
    <textarea id="f-desc" placeholder="What does this service do?"></textarea>
    <div class="modal-actions">
      <button onclick="closeModal()">Cancel</button>
      <button class="btn-save" onclick="register()">Register</button>
    </div>
  </div>
</div>

<script>
const API = '/services';
const $ = id => document.getElementById(id);

async function loadServices() {
  const res = await fetch(API);
  const svcs = await res.json();
  const tbody = $('tbody');

  if (!svcs.length) {
    tbody.innerHTML = '<tr><td colspan="5" class="empty"><p>No services registered</p><p style="font-size:.8rem">Use the button above or POST to /services to register one.</p></td></tr>';
    $('stats').textContent = '';
    return;
  }

  const healthy = svcs.filter(s => s.status === 'healthy').length;
  const unhealthy = svcs.filter(s => s.status === 'unhealthy').length;
  $('stats').textContent = svcs.length + ' services | ' + healthy + ' healthy | ' + unhealthy + ' unhealthy';

  tbody.innerHTML = svcs.map(s => {
    const desc = s.description ? '<div class="svc-desc">' + esc(s.description) + '</div>' : '';
    const badge = 'badge badge-' + s.status;
    const checked = s.last_checked_at ? timeAgo(new Date(s.last_checked_at)) : '-';
    const displayUrl = s.display_url || s.url;
    let urlCell = '<a href="' + esc(displayUrl) + '" target="_blank" rel="noopener noreferrer">' + esc(displayUrl) + '</a>';
    if (s.external_url) {
      urlCell += '<div class="svc-ext"><a href="' + esc(s.external_url) + '" target="_blank" rel="noopener noreferrer">External: ' + esc(s.external_url) + '</a></div>';
    }
    return '<tr>' +
      '<td><div class="svc-name">' + esc(s.name) + '</div>' + desc + '</td>' +
      '<td class="svc-url">' + urlCell + '</td>' +
      '<td><span class="' + badge + '">' + s.status + '</span></td>' +
      '<td class="time">' + checked + '</td>' +
      '<td><button class="btn-del" onclick="del(\'' + s.id + '\',\'' + esc(s.name) + '\')">Remove</button></td>' +
      '</tr>';
  }).join('');
}

async function register() {
  const name = $('f-name').value.trim();
  const url = $('f-url').value.trim();
  const desc = $('f-desc').value.trim();
  if (!name || !url) { alert('Name and URL are required'); return; }
  const res = await fetch(API, {
    method: 'POST',
    headers: {'Content-Type': 'application/json'},
    body: JSON.stringify({name, url, description: desc})
  });
  if (!res.ok) {
    const e = await res.json();
    alert('Error: ' + (e.error || res.statusText));
    return;
  }
  closeModal();
  loadServices();
}

async function del(id, name) {
  if (!confirm('Remove "' + name + '" from registry?')) return;
  await fetch(API + '/' + id, {method: 'DELETE'});
  loadServices();
}

function openModal() {
  $('f-name').value = ''; $('f-url').value = ''; $('f-desc').value = '';
  $('modal').classList.add('active');
  $('f-name').focus();
}
function closeModal() { $('modal').classList.remove('active'); }
$('modal').addEventListener('click', e => { if (e.target === $('modal')) closeModal(); });
document.addEventListener('keydown', e => { if (e.key === 'Escape') closeModal(); });

function esc(s) { const d = document.createElement('div'); d.textContent = s; return d.innerHTML; }

function timeAgo(date) {
  const s = Math.floor((Date.now() - date.getTime()) / 1000);
  if (s < 10) return 'just now';
  if (s < 60) return s + 's ago';
  const m = Math.floor(s / 60);
  if (m < 60) return m + 'm ago';
  const h = Math.floor(m / 60);
  if (h < 24) return h + 'h ago';
  return date.toLocaleString();
}

loadServices();
setInterval(loadServices, 15000);

const AI_PROMPT = ` + "`" + `Please add a self-registration feature to the current project's web service. Requirements:

1. On service startup, asynchronously send a POST request to http://localhost:1234/services to register itself.
2. The request body is JSON: {"name": "<service name>", "url": "<root URL>", "description": "<brief description>"}.
   - "name" should be the project name or a reasonable identifier.
   - "url" must be the root URL of this service (e.g. http://localhost:<port>), NOT an API endpoint or health check path. The registry uses this URL both as the homepage link and for health checks (GET that returns 2xx).
   - "description" should briefly describe what this service does.
3. The registration must be fire-and-forget: run it in a background goroutine / async task / non-blocking thread. It must NOT block the main startup process.
4. ALL errors (connection refused, timeout, non-2xx response, etc.) must be silently ignored — at most print a debug-level log. The service must start and run normally regardless of whether registration succeeds.
5. Do not add any new dependencies if possible; use the standard library / built-in HTTP client.
6. Fit the implementation style and conventions of the existing codebase.
` + "`" + `;

function copyPrompt() {
  navigator.clipboard.writeText(AI_PROMPT.trim()).then(() => {
    const btn = $('btn-prompt');
    btn.textContent = 'Copied!';
    btn.classList.add('copied');
    setTimeout(() => { btn.textContent = 'Copy AI Prompt'; btn.classList.remove('copied'); }, 2000);
  });
}
</script>
</body>
</html>`
