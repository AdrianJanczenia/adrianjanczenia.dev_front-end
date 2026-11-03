(() => {
    const header = document.querySelector('[data-site-header]');
    const menuButton = document.querySelector('[data-menu-toggle]');

    if (header && menuButton) {
        const closeMenu = () => {
            menuButton.setAttribute('aria-expanded', 'false');
            header.classList.remove('menu-open');
        };

        menuButton.addEventListener('click', () => {
            const isOpen = menuButton.getAttribute('aria-expanded') !== 'true';
            menuButton.setAttribute('aria-expanded', String(isOpen));
            header.classList.toggle('menu-open', isOpen);
        });

        document.addEventListener('keydown', event => {
            if (event.key === 'Escape' && header.classList.contains('menu-open')) {
                closeMenu();
                menuButton.focus();
            }
        });

        header.querySelectorAll('.site-nav a').forEach(link => {
            link.addEventListener('click', closeMenu);
        });
    }

    const toast = document.getElementById('copy-toast');
    if (!toast) return;

    let hideToastTimer;
    let finishHidingTimer;

    async function copyEmail(email) {
        if (navigator.clipboard?.writeText) {
            try {
                await navigator.clipboard.writeText(email);
                return true;
            } catch {
                // The selection-based fallback also works outside secure contexts.
            }
        }

        const field = document.createElement('textarea');
        field.value = email;
        field.readOnly = true;
        field.style.position = 'fixed';
        field.style.opacity = '0';
        document.body.append(field);

        try {
            field.select();
            return document.execCommand('copy');
        } catch {
            return false;
        } finally {
            field.remove();
        }
    }

    function showCopiedToast() {
        clearTimeout(hideToastTimer);
        clearTimeout(finishHidingTimer);
        toast.hidden = false;
        requestAnimationFrame(() => toast.classList.add('is-visible'));

        hideToastTimer = setTimeout(() => {
            toast.classList.remove('is-visible');
            finishHidingTimer = setTimeout(() => {
                toast.hidden = true;
            }, 220);
        }, 2400);
    }

    document.querySelectorAll('[data-copy-email]').forEach(link => {
        link.addEventListener('click', async event => {
            const modifiedClick = event.button !== 0
                || event.metaKey
                || event.ctrlKey
                || event.shiftKey
                || event.altKey;
            if (event.defaultPrevented || modifiedClick) return;

            event.preventDefault();
            const copied = await copyEmail(link.dataset.copyEmail);
            if (copied) showCopiedToast();

            setTimeout(() => {
                window.location.href = link.href;
            }, copied ? 650 : 0);
        });
    });
})();
