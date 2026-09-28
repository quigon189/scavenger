document.body.addEventListener('toast', (e) => {
    const { text, type = 'info' } = e.detail || {};
    const el = document.createElement('div');
    el.className = `toast align-items-center text-bg-${type} border-0`;
    el.innerHTML = `
        <div class="d-flex">
            <div class="toast-body">${text}</div>
            <button class="btn-close btn-close-white me-2 m-auto" data-bs-dismiss="toast"></button>
        </div>`;
    const host = document.getElementById('toasts');
    if (!host) return;
    host.appendChild(el);
    new bootstrap.Toast(el).show();
    el.addEventListener('hidden.bs.toast', () => el.remove());
});

document.body.addEventListener('htmx:responseError', (e) => {
    console.error('htmx error', e.detail);
});
