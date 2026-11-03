(() => {
    const modal = document.getElementById('cv-modal');
    const form = document.getElementById('cv-form');
    const errorCopyElement = document.getElementById('cv-error-copy');
    if (!modal || !form || !errorCopyElement) return;

    const errorCopy = JSON.parse(errorCopyElement.textContent) || {};
    const body = document.body;
    const siteShell = document.querySelector('.site-shell');
    const globalLoader = document.getElementById('global-loader');
    const captchaImage = document.getElementById('captcha-img');
    const errorMessage = document.getElementById('modal-error');
    const formFields = document.getElementById('modal-fields');
    const retryButton = document.getElementById('cv-retry');
    const submitButton = form.querySelector('.modal-btn');
    const accessCodeInput = document.getElementById('cv-access-code');
    const captchaInput = document.getElementById('cv-captcha');
    const honeypotInput = document.getElementById('cv-hp-name');
    const closeButton = modal.querySelector('.modal-close');

    let captchaId = '';
    let lastTrigger = null;

    function translatedError(slug) {
        return errorCopy[slug]?.message
            || errorCopy.error_service_unavailable?.message
            || errorCopy.error_message?.message
            || '';
    }

    function showError(slug) {
        errorMessage.textContent = translatedError(slug);
        errorMessage.style.display = 'block';
    }

    function hideError() {
        errorMessage.textContent = '';
        errorMessage.style.display = 'none';
    }

    function closeModal() {
        modal.classList.remove('open');
        siteShell.inert = false;
        body.classList.remove('no-scroll');
        accessCodeInput.value = '';
        captchaInput.value = '';
        lastTrigger?.focus();
    }

    async function calculateNonce(seed) {
        const encoder = new TextEncoder();
        const requiredPrefix = '0000';
        const maxNonce = 1_000_000;

        for (let nonce = 0; nonce <= maxNonce; nonce++) {
            const input = encoder.encode(seed + nonce);
            const digest = await crypto.subtle.digest('SHA-256', input);
            const hash = Array.from(new Uint8Array(digest), byte =>
                byte.toString(16).padStart(2, '0')
            ).join('');

            if (hash.startsWith(requiredPrefix)) {
                return String(nonce);
            }
        }

        throw new Error('error_pow_failed');
    }

    async function loadChallenge() {
        captchaId = '';
        submitButton.disabled = true;
        retryButton.hidden = true;

        try {
            const challengeResponse = await fetch('/api/pow');
            const challenge = await challengeResponse.json();
            if (!challengeResponse.ok) {
                throw new Error(challenge.error || 'error_pow_failed');
            }
            if (!challenge.seed || !challenge.signature) {
                throw new Error('error_pow_failed');
            }

            const nonce = await calculateNonce(challenge.seed);
            const captchaResponse = await fetch('/api/captcha', {
                method: 'POST',
                headers: {'Content-Type': 'application/json'},
                body: JSON.stringify({
                    seed: challenge.seed,
                    signature: challenge.signature,
                    nonce
                })
            });
            const captcha = await captchaResponse.json();
            if (!captchaResponse.ok) {
                throw new Error(captcha.error || 'error_internal');
            }
            if (!captcha.captchaId || !captcha.captchaImg) {
                throw new Error('error_captcha_server');
            }

            captchaId = captcha.captchaId;
            captchaImage.src = captcha.captchaImg;
            formFields.style.display = 'block';
            submitButton.disabled = false;
            hideError();
            return true;
        } catch (error) {
            showError(error.message);
            formFields.style.display = 'none';
            retryButton.hidden = false;
            return false;
        }
    }

    async function refreshExpiredCaptcha(slug) {
        captchaInput.value = '';
        const refreshed = await loadChallenge();
        if (refreshed) showError(slug);
    }

    async function openModal(trigger) {
        if (globalLoader.classList.contains('visible')) return;

        lastTrigger = trigger;
        globalLoader.classList.add('visible');
        accessCodeInput.value = '';
        captchaInput.value = '';
        hideError();
        formFields.style.display = 'block';

        const challengeLoaded = await loadChallenge();
        globalLoader.classList.remove('visible');
        modal.classList.add('open');
        body.classList.add('no-scroll');
        siteShell.inert = true;

        if (challengeLoaded) accessCodeInput.focus();
        else retryButton.focus();
    }

    document.addEventListener('click', event => {
        const link = event.target.closest('a.cv-link');
        if (!link) return;
        event.preventDefault();
        void openModal(link);
    });

    closeButton.addEventListener('click', closeModal);
    modal.addEventListener('click', event => {
        if (event.target === modal) closeModal();
    });

    document.addEventListener('keydown', event => {
        if (!modal.classList.contains('open')) return;

        if (event.key === 'Escape') {
            closeModal();
            return;
        }
        if (event.key !== 'Tab') return;

        const focusable = [...modal.querySelectorAll('button, input')].filter(element =>
            !element.disabled
            && !element.hidden
            && element.tabIndex >= 0
            && element.getClientRects().length > 0
        );
        if (focusable.length === 0) return;

        const first = focusable[0];
        const last = focusable[focusable.length - 1];
        if (event.shiftKey && document.activeElement === first) {
            event.preventDefault();
            last.focus();
        } else if (!event.shiftKey && document.activeElement === last) {
            event.preventDefault();
            first.focus();
        }
    });

    retryButton.addEventListener('click', async () => {
        retryButton.disabled = true;
        globalLoader.classList.add('visible');
        const challengeLoaded = await loadChallenge();
        globalLoader.classList.remove('visible');
        retryButton.disabled = false;
        if (challengeLoaded) accessCodeInput.focus();
    });

    const backToTopButton = document.querySelector('.back-to-top');
    if (backToTopButton) {
        window.addEventListener('scroll', () => {
            backToTopButton.classList.toggle('visible', window.scrollY > 300);
        });
    }

    form.addEventListener('submit', async event => {
        event.preventDefault();
        if (honeypotInput.value !== '' || !captchaId) return;

        const language = document.documentElement.lang || 'pl';
        submitButton.classList.add('loading');
        submitButton.disabled = true;
        hideError();

        try {
            const verificationResponse = await fetch('/api/captcha-verify', {
                method: 'POST',
                headers: {'Content-Type': 'application/json'},
                body: JSON.stringify({
                    captchaId,
                    captchaValue: captchaInput.value
                })
            });
            const verification = await verificationResponse.json();

            if (!verificationResponse.ok) {
                if (isExpiredCaptcha(verification.error)) {
                    await refreshExpiredCaptcha(verification.error);
                    return;
                }
                throw new Error(verification.error || 'error_captcha_invalid');
            }

            const tokenResponse = await fetch('/api/cv-token', {
                method: 'POST',
                headers: {'Content-Type': 'application/json'},
                body: JSON.stringify({
                    password: accessCodeInput.value,
                    lang: language,
                    captchaId: verification.captchaId,
                    fullName: honeypotInput.value
                })
            });
            const tokenResult = await tokenResponse.json();

            if (!tokenResponse.ok) {
                if (isExpiredCaptcha(tokenResult.error)) {
                    await refreshExpiredCaptcha(tokenResult.error);
                    return;
                }
                throw new Error(tokenResult.error || 'error_cv_server');
            }
            if (!tokenResult.token) throw new Error('error_cv_server');

            const downloadURL = new URL('/api/download/cv', window.location.origin);
            downloadURL.searchParams.set('token', tokenResult.token);
            downloadURL.searchParams.set('lang', language);
            closeModal();
            window.location.assign(downloadURL.toString());
        } catch (error) {
            showError(error.message);
        } finally {
            submitButton.classList.remove('loading');
            submitButton.disabled = !captchaId;
        }
    });

    function isExpiredCaptcha(slug) {
        return slug === 'error_captcha_expired' || slug === 'error_captcha_not_found';
    }

    window.addEventListener('pageshow', () => {
        globalLoader.classList.remove('visible');
    });
})();
