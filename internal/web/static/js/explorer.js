(() => {
    const root = document.getElementById('request-explorer');
    if (!root) return;
    const find = selector => root.querySelector(selector);
    const findAll = selector => [...root.querySelectorAll(selector)];
    const copyElement = document.getElementById('explorer-copy');
    if (!copyElement) return;
    const copy = JSON.parse(copyElement.textContent);
    const scenarioText = key => {
        const [group, id] = key.split('.', 2);
        return copy.scenario_copy[group]?.[id] ?? '';
    };
    const ui = copy.ui_copy;
    let view = 'all';
    let variant = 0;
    let stepIndex = 0;
    let busy = false;
    let paused = false;
    let generation = 0;
    let selectedComponent = null;
    const common = {
        incoming: scenarioText('common.request_reaches_frontend'),
        noChange: scenarioText('common.no_state_change')
    };
    const endpoints = {
        content: {
            label: 'GET /',
            title: scenarioText('content.page_content_and_ssr'),
            description: scenarioText('content.rendering_overview'),
            nodes: ['browser', 'caddy', 'front', 'gateway', 'content', 'cache'],
            variants: [
                scenarioText('content.polish'),
                scenarioText('content.english'),
                scenarioText('content.unknown_language'),
                scenarioText('content.gateway_unavailable')
            ],
            build(variantIndex) {
                const requestedLanguage = variantIndex === 1 ? 'en' : variantIndex === 2 ? 'de' : 'pl';
                const contentLanguage = requestedLanguage === 'en' ? 'en' : 'pl';
                const steps = [{
                    label: scenarioText('content.open_the_page'),
                    path: ['browser', 'caddy', 'front'],
                    kind: 'request',
                    request: `GET /?lang=${requestedLanguage}`,
                    response: '—',
                    state: common.noChange,
                    decision: requestedLanguage === 'de'
                        ? scenarioText('content.unknown_language_uses_polish')
                        : scenarioText('content.the_front_end_starts_ssr'),
                    description: common.incoming
                }, {
                    label: scenarioText('content.fetch_content'),
                    path: ['front', 'gateway'],
                    kind: 'request',
                    request: `GET /api/v1/content?lang=${contentLanguage}`,
                    response: '—',
                    state: common.noChange,
                    decision: scenarioText('content.gateway_forwards_the_request'),
                    description: scenarioText('content.request_content_from_gateway')
                }];
                if (variantIndex === 3) return [...steps, {
                    label: scenarioText('content.return_service_error'),
                    path: ['gateway', 'front', 'caddy', 'browser'],
                    kind: 'response',
                    request: '—',
                    response: '503 Service Unavailable',
                    state: common.noChange,
                    decision: scenarioText('content.render_service_error_page'),
                    description: scenarioText('content.gateway_failure_stops_rendering')
                }];
                return [...steps, {
                    label: scenarioText('content.read_cached_content'),
                    path: ['gateway', 'content', 'cache', 'content'],
                    kind: 'request',
                    request: `gRPC ContentService.Handle(lang=${contentLanguage})`,
                    response: `content.lang=${contentLanguage}`,
                    state: scenarioText('content.json_loaded_on_startup_no_file_read'),
                    decision: requestedLanguage === 'de'
                        ? scenarioText('content.polish_fallback_forwarded')
                        : scenarioText('content.content_service_chooses_a_language'),
                    description: scenarioText('content.localized_content_response')
                }, {
                    label: scenarioText('content.return_rendered_page'),
                    path: ['content', 'gateway', 'front', 'caddy', 'browser'],
                    kind: 'response',
                    request: '—',
                    response: `200 OK · HTML lang=${contentLanguage}`,
                    state: common.noChange,
                    decision: scenarioText('content.render_html_from_content'),
                    description: scenarioText('content.rendered_html_response_path')
                }];
            }
        },
        pow: {
            label: 'GET /api/pow',
            title: scenarioText('pow.signed_pow_challenge'),
            description: scenarioText('pow.challenge_request_overview'),
            nodes: ['browser', 'caddy', 'front', 'gateway', 'captcha'],
            variants: [scenarioText('pow.fetch_challenge')],
            build() {
                return [{
                    label: scenarioText('pow.request_a_challenge'),
                    path: ['browser', 'caddy', 'front', 'gateway', 'captcha'],
                    kind: 'request',
                    request: 'GET /api/pow',
                    response: '—',
                    state: common.noChange,
                    decision: scenarioText('pow.captcha_service_signs_seed'),
                    description: scenarioText('pow.browser_requests_challenge')
                }, {
                    label: scenarioText('pow.receive_seed'),
                    path: ['captcha', 'gateway', 'front', 'caddy', 'browser'],
                    kind: 'response',
                    request: '—',
                    response: '200 {seed:"demo-seed", signature:"demo-signature"}',
                    state: common.noChange,
                    decision: scenarioText('pow.hmac_secret_remains_private'),
                    description: scenarioText('pow.seed_response_precedes_nonce_search')
                }];
            }
        },
        captcha: {
            label: 'POST /api/captcha',
            title: scenarioText('captcha.pow_and_captcha_image'),
            description: scenarioText('captcha.prior_challenge_context'),
            nodes: ['browser', 'caddy', 'front', 'gateway', 'captcha', 'redis'],
            variants: [
                scenarioText('captcha.valid_pow'),
                scenarioText('captcha.bad_signature'),
                scenarioText('captcha.expired_seed'),
                scenarioText('captcha.weak_hash'),
                scenarioText('captcha.reused_seed')
            ],
            build(variantIndex) {
                const nonce = variantIndex === 3 ? '1800' : '19228';
                const hash = variantIndex === 3 ? '9a1ed77daa464506' : '0000045a987ec2c6';
                const steps = [{
                    label: scenarioText('captcha.calculate_nonce_locally'),
                    path: [],
                    kind: 'local',
                    request: 'SHA-256("demo-seed" + nonce)',
                    response: `nonce=${nonce} · ${hash}…`,
                    state: common.noChange,
                    decision: scenarioText('captcha.nonce_calculated_in_browser'),
                    description: scenarioText('captcha.nonce_search_visualization')
                }, {
                    label: scenarioText('captcha.submit_solution'),
                    path: ['browser', 'caddy', 'front', 'gateway', 'captcha'],
                    kind: 'request',
                    request: `POST /api/captcha {seed:"demo-seed", nonce:"${nonce}", signature:"demo-signature"}`,
                    response: '—',
                    state: common.noChange,
                    decision: scenarioText('captcha.validate_pow_signature_time_hash'),
                    description: scenarioText('captcha.captcha_service_receives_request')
                }];
                const failures = [
                    null,
                    ['403 Forbidden', scenarioText('captcha.signature_mismatch_no_redis_write')],
                    ['410 Gone', scenarioText('captcha.challenge_expired_no_redis_write')],
                    ['400 Bad Request', scenarioText('captcha.weak_hash_no_redis_write')]
                ];
                if (variantIndex > 0 && variantIndex < 4) return [...steps, {
                    label: scenarioText('captcha.reject_captcha_creation'),
                    path: ['captcha', 'gateway', 'front', 'caddy', 'browser'],
                    kind: 'response',
                    request: '—',
                    response: failures[variantIndex][0],
                    state: common.noChange,
                    decision: failures[variantIndex][1],
                    description: scenarioText('captcha.invalid_solution_not_stored')
                }];
                return [...steps, {
                    label: scenarioText('captcha.mark_seed_in_redis'),
                    path: ['captcha', 'redis', 'captcha'],
                    kind: 'request',
                    request: 'SETNX pow:demo-seed',
                    response: variantIndex === 4 ? 'false' : 'true',
                    state: variantIndex === 4
                        ? scenarioText('captcha.key_already_exists')
                        : scenarioText('captcha.pow_demo_seed_used_ttl'),
                    decision: scenarioText('captcha.atomic_seed_reuse_check'),
                    description: variantIndex === 4
                        ? scenarioText('captcha.redis_rejects_reuse')
                        : scenarioText('captcha.the_seed_can_be_used_once')
                }, {
                    label: variantIndex === 4
                        ? scenarioText('captcha.return_409')
                        : scenarioText('captcha.return_captcha_image'),
                    path: ['captcha', 'gateway', 'front', 'caddy', 'browser'],
                    kind: 'response',
                    request: '—',
                    response: variantIndex === 4 ? '409 Conflict' : '200 {captchaId:"demo-id", image:"…"}',
                    state: variantIndex === 4 ? common.noChange : scenarioText('captcha.captcha_demo_id_session_ttl'),
                    decision: variantIndex === 4
                        ? scenarioText('captcha.a_reused_seed_creates_no_new_image')
                        : scenarioText('captcha.captcha_receives_a_short_ttl'),
                    description: variantIndex === 4
                        ? scenarioText('captcha.the_rejection_returns_to_the_client')
                        : scenarioText('captcha.the_image_returns_to_the_browser')
                }];
            }
        },
        verify: {
            label: 'POST /api/captcha-verify',
            title: scenarioText('verify.captcha_answer_verification'),
            description: scenarioText('verify.verification_overview'),
            nodes: ['browser', 'caddy', 'front', 'gateway', 'captcha', 'redis'],
            variants: [
                scenarioText('verify.correct_answer'),
                scenarioText('verify.wrong_answer'),
                scenarioText('verify.expired_session')
            ],
            build(variantIndex) {
                return [{
                    label: scenarioText('verify.submit_captcha_answer'),
                    path: ['browser', 'caddy', 'front', 'gateway', 'captcha'],
                    kind: 'request',
                    request: 'POST /api/captcha-verify {captchaId:"demo-id", captchaValue:"DEMO"}',
                    response: '—',
                    state: common.noChange,
                    decision: scenarioText('verify.captcha_session_lookup'),
                    description: scenarioText('verify.captcha_answer_arrives')
                }, {
                    label: scenarioText('verify.check_session_in_redis'),
                    path: ['captcha', 'redis', 'captcha'],
                    kind: 'request',
                    request: 'Lua: verify captcha:demo-id',
                    response: variantIndex === 0 ? 'solved' : variantIndex === 1 ? 'invalid_value' : 'not_found',
                    state: variantIndex === 0
                        ? scenarioText('verify.captcha_demo_id_solved_true')
                        : variantIndex === 1
                        ? scenarioText('verify.captcha_demo_id_triesleft_1')
                        : scenarioText('verify.key_expired'),
                    decision: variantIndex === 0
                        ? scenarioText('verify.solved_session_ready_for_cv')
                        : scenarioText('verify.captcha_remains_unsolved'),
                    description: scenarioText('verify.atomic_captcha_state_update')
                }, {
                    label: scenarioText('verify.receive_verification_result'),
                    path: ['captcha', 'gateway', 'front', 'caddy', 'browser'],
                    kind: 'response',
                    request: '—',
                    response: variantIndex === 0
                        ? '200 {captchaId:"demo-id"}'
                        : variantIndex === 1 ? '400 Bad Request' : '404 Not Found',
                    state: common.noChange,
                    decision: variantIndex === 0
                        ? scenarioText('verify.the_client_may_request_a_cv_token')
                        : scenarioText('verify.verification_error_returned'),
                    description: scenarioText('verify.verification_result_to_form')
                }];
            }
        },
        cvtoken: {
            label: 'POST /api/cv-token',
            title: scenarioText('cvtoken.single_use_cv_token'),
            description: scenarioText('cvtoken.cv_token_request_overview'),
            nodes: ['browser', 'caddy', 'front', 'gateway', 'rabbit', 'content', 'redis'],
            variants: [
                scenarioText('cvtoken.valid_input'),
                scenarioText('cvtoken.wrong_access_code'),
                scenarioText('cvtoken.wrong_captcha'),
                scenarioText('cvtoken.expired_captcha')
            ],
            build(variantIndex) {
                const steps = [{
                    label: scenarioText('cvtoken.request_a_token'),
                    path: ['browser', 'caddy', 'front', 'gateway'],
                    kind: 'request',
                    request: 'POST /api/cv-token {password:"••••", lang:"pl", captchaId:"demo-id"}',
                    response: '—',
                    state: common.noChange,
                    decision: scenarioText('cvtoken.gateway_sends_token_rpc'),
                    description: scenarioText('cvtoken.cv_form_submits_access_code')
                }, {
                    label: scenarioText('cvtoken.dispatch_through_rabbitmq'),
                    path: ['gateway', 'rabbit', 'content'],
                    kind: 'request',
                    request: 'RabbitMQ RPC cv-request',
                    response: '—',
                    state: common.noChange,
                    decision: scenarioText('cvtoken.token_request_reaches_content'),
                    description: scenarioText('cvtoken.token_request_delivered_by_rpc')
                }, {
                    label: scenarioText('cvtoken.validate_and_store_state'),
                    path: ['content', 'redis', 'content'],
                    kind: 'request',
                    request: 'Lua IssueToken(captcha:demo-id, cv:token:tok_demo)',
                    response: variantIndex === 0
                        ? 'token=tok_demo'
                        : variantIndex === 3 ? 'captcha expired' : 'validation failed',
                    state: variantIndex === 0
                        ? scenarioText('cvtoken.issued_token_redis_state')
                        : scenarioText('cvtoken.no_new_cv_token'),
                    decision: variantIndex === 0
                        ? scenarioText('cvtoken.captcha_and_code_are_valid')
                        : scenarioText('cvtoken.validation_denies_the_token'),
                    description: variantIndex === 0
                        ? scenarioText('cvtoken.a_short_lived_token_is_ready')
                        : scenarioText('cvtoken.invalid_input_does_not_issue_a_token')
                }];
                const status = variantIndex === 0
                    ? '200 {token:"tok_demo"}'
                    : variantIndex === 1
                    ? '401 Unauthorized'
                    : variantIndex === 3
                    ? '410 Gone'
                    : '403 Forbidden';
                return [...steps, {
                    label: scenarioText('cvtoken.receive_result'),
                    path: ['content', 'rabbit', 'gateway', 'front', 'caddy', 'browser'],
                    kind: 'response',
                    request: '—',
                    response: status,
                    state: common.noChange,
                    decision: variantIndex === 0
                        ? scenarioText('cvtoken.token_response_to_dialog')
                        : scenarioText('cvtoken.dialog_displays_generic_error'),
                    description: scenarioText('cvtoken.token_rpc_response_path')
                }];
            }
        },
        download: {
            label: 'GET /api/download/cv',
            title: scenarioText('download.cv_file_download'),
            description: scenarioText('download.download_order_overview'),
            nodes: ['browser', 'caddy', 'front', 'gateway', 'content', 'redis'],
            variants: [
                scenarioText('download.download_and_retry'),
                scenarioText('download.token_language_mismatch'),
                scenarioText('download.expired_token'),
                scenarioText('download.pdf_missing')
            ],
            build(variantIndex) {
                const requestedLanguage = variantIndex === 1 ? 'en' : 'pl';
                const steps = [{
                    label: scenarioText('download.request_the_pdf'),
                    path: ['browser', 'caddy', 'front', 'gateway', 'content'],
                    kind: 'request',
                    request: `GET /api/download/cv?token=tok_demo&lang=${requestedLanguage}`,
                    response: '—',
                    state: common.noChange,
                    decision: scenarioText('download.gateway_calls_download_service'),
                    description: scenarioText('download.download_request_arrives')
                }, {
                    label: scenarioText('download.open_the_file'),
                    path: [],
                    kind: 'local',
                    request: `open(cv_${requestedLanguage}.pdf)`,
                    response: variantIndex === 3 ? 'ENOENT' : 'file ready',
                    state: common.noChange,
                    decision: variantIndex === 3
                        ? scenarioText('download.file_failure_precedes_token_use')
                        : scenarioText('download.the_file_is_ready_to_stream'),
                    description: variantIndex === 3
                        ? scenarioText('download.missing_pdf_preserves_token')
                        : scenarioText('download.file_availability_is_checked_first')
                }];
                if (variantIndex === 3) return [...steps, {
                    label: scenarioText('download.return_file_error'),
                    path: ['content', 'gateway', 'front', 'caddy', 'browser'],
                    kind: 'response',
                    request: '—',
                    response: '500 Internal Server Error',
                    state: scenarioText('download.cv_token_tok_demo_retained'),
                    decision: scenarioText('download.no_token_consumption_was_attempted'),
                    description: scenarioText('download.the_error_returns_to_the_browser')
                }];
                steps.push({
                    label: scenarioText('download.check_token_in_redis'),
                    path: ['content', 'redis', 'content'],
                    kind: 'request',
                    request: 'Lua: validate token + language; delete on success',
                    response: variantIndex === 0
                        ? 'valid → DEL'
                        : variantIndex === 1 ? 'language mismatch' : 'missing / expired',
                    state: variantIndex === 0
                        ? scenarioText('download.cv_token_tok_demo_atomically_deleted')
                        : variantIndex === 1
                        ? scenarioText('download.cv_token_tok_demo_retained')
                        : scenarioText('download.key_missing'),
                    decision: variantIndex === 0
                        ? scenarioText('download.language_matches_token_is_consumed')
                        : variantIndex === 1
                        ? scenarioText('download.language_mismatch_token_remains')
                        : scenarioText('download.the_token_expired_before_download'),
                    description: scenarioText('download.atomic_token_validation')
                });
                steps.push({
                    label: scenarioText('download.receive_response'),
                    path: ['content', 'gateway', 'front', 'caddy', 'browser'],
                    kind: 'response',
                    request: '—',
                    response: variantIndex === 0
                        ? '200 application/pdf'
                        : variantIndex === 1 ? '400 Bad Request' : '410 Gone',
                    state: common.noChange,
                    decision: variantIndex === 0
                        ? scenarioText('download.the_pdf_returns_as_a_stream')
                        : scenarioText('download.no_cv_is_sent'),
                    description: scenarioText('download.download_response_path')
                });
                if (variantIndex === 0) {
                    steps.push({
                        label: scenarioText('download.repeat_the_same_get'),
                        path: ['browser', 'caddy', 'front', 'gateway', 'content'],
                        kind: 'request',
                        request: 'GET /api/download/cv?token=tok_demo&lang=pl',
                        response: '—',
                        state: common.noChange,
                        decision: scenarioText('download.retry_reaches_content_service'),
                        description: scenarioText('download.retry_reuses_token')
                    });
                    steps.push({
                        label: scenarioText('download.open_the_pdf_again'),
                        path: [],
                        kind: 'local',
                        request: 'open(cv_pl.pdf)',
                        response: 'file ready',
                        state: common.noChange,
                        decision: scenarioText('download.pdf_available_token_may_be_gone'),
                        description: scenarioText('download.retry_checks_pdf_before_token')
                    });
                    steps.push({
                        label: scenarioText('download.check_consumed_token'),
                        path: ['content', 'redis', 'content'],
                        kind: 'request',
                        request: 'GET cv:token:tok_demo',
                        response: 'missing',
                        state: scenarioText('download.cv_token_tok_demo_absent'),
                        decision: scenarioText('download.the_token_was_already_consumed'),
                        description: scenarioText('download.consumed_token_missing_in_redis')
                    });
                    steps.push({
                        label: scenarioText('download.receive_410_gone'),
                        path: ['content', 'gateway', 'front', 'caddy', 'browser'],
                        kind: 'response',
                        request: '—',
                        response: '410 Gone',
                        state: common.noChange,
                        decision: scenarioText('download.a_repeat_download_is_not_possible'),
                        description: scenarioText('download.the_rejection_returns_to_the_browser')
                    });
                }
                return steps;
            }
        },
        route: {
            label: 'GET /unknown',
            title: scenarioText('route.unknown_route'),
            description: scenarioText('route.unknown_route_overview'),
            nodes: ['browser', 'caddy', 'front'],
            variants: [scenarioText('route.browser_html'), scenarioText('route.api_client_json')],
            build(variantIndex) {
                return [{
                    label: scenarioText('route.open_an_unknown_route'),
                    path: ['browser', 'caddy', 'front'],
                    kind: 'request',
                    request: `GET /unknown\nAccept: ${variantIndex === 0 ? 'text/html' : 'application/json'}`,
                    response: '—',
                    state: common.noChange,
                    decision: scenarioText('route.frontend_detects_unknown_route'),
                    description: scenarioText('route.the_request_stops_at_the_front_end')
                }, {
                    label: scenarioText('route.receive_error_response'),
                    path: ['front', 'caddy', 'browser'],
                    kind: 'response',
                    request: '—',
                    response: variantIndex === 0 ? '404 text/html' : '404 application/json',
                    state: common.noChange,
                    decision: variantIndex === 0
                        ? scenarioText('route.shared_error_template')
                        : scenarioText('route.json_response_for_api'),
                    description: scenarioText('route.frontend_returns_404_locally')
                }];
            }
        }
    };

    const viewNames = {
        content: scenarioText('views.home_page'),
        pow: scenarioText('views.anti_bot_challenge'),
        captcha: scenarioText('views.create_captcha'),
        verify: scenarioText('views.verify_captcha'),
        cvtoken: scenarioText('views.cv_access_request'),
        download: scenarioText('views.download_cv'),
        route: scenarioText('views.unknown_url')
    };
    const componentCopy = {
        browser: {
            title: scenarioText('components.browser'),
            description: scenarioText('components.browser_responsibilities')
        },
        caddy: {
            title: 'Caddy',
            description: scenarioText('components.caddy_responsibilities')
        },
        front: {title: 'Front end', description: scenarioText('components.frontend_responsibilities')},
        gateway: {title: 'Gateway', description: scenarioText('components.gateway_responsibilities')},
        rabbit: {
            title: 'RabbitMQ',
            description: scenarioText('components.rabbitmq_responsibilities')
        },
        content: {
            title: 'Content service',
            description: scenarioText('components.content_service_responsibilities')
        },
        cache: {
            title: scenarioText('components.content_in_memory'),
            description: scenarioText('components.content_cache_responsibilities')
        },
        captcha: {
            title: 'CAPTCHA service',
            description: scenarioText('components.captcha_service_responsibilities')
        },
        redis: {title: 'Redis', description: scenarioText('components.redis_responsibilities')}
    };
    const variantMeta = {
        content: [
            {type: 'main', why: scenarioText('branches.render_polish_page')},
            {type: 'main', why: scenarioText('branches.render_english_page')},
            {type: 'extra', why: scenarioText('branches.fallback_for_unknown_language')},
            {type: 'extra', why: scenarioText('branches.unavailable_gateway_error')}
        ],
        pow: [{type: 'main', why: scenarioText('branches.fetch_pow_challenge')}],
        captcha: [
            {type: 'main', why: scenarioText('branches.valid_nonce_creates_captcha')},
            {type: 'extra', why: scenarioText('branches.invalid_signature_rejected')},
            {type: 'extra', why: scenarioText('branches.expired_seed_rejected')},
            {type: 'extra', why: scenarioText('branches.weak_pow_rejected')},
            {type: 'extra', why: scenarioText('branches.seed_reuse_rejected')}
        ],
        verify: [
            {type: 'main', why: scenarioText('branches.correct_answer_solves_captcha')},
            {type: 'extra', why: scenarioText('branches.wrong_answer_reduces_attempts')},
            {type: 'extra', why: scenarioText('branches.expired_captcha_rejected')}
        ],
        cvtoken: [
            {type: 'main', why: scenarioText('branches.valid_cv_access_issues_token')},
            {type: 'extra', why: scenarioText('branches.wrong_access_code_rejected')},
            {type: 'extra', why: scenarioText('branches.unsolved_captcha_rejected')},
            {type: 'extra', why: scenarioText('branches.expired_captcha_token_rejected')}
        ],
        download: [
            {type: 'main', why: scenarioText('branches.first_download_then_retry')},
            {type: 'extra', why: scenarioText('branches.language_mismatch_keeps_token')},
            {type: 'extra', why: scenarioText('branches.expired_token_rejected')},
            {type: 'extra', why: scenarioText('branches.missing_pdf_keeps_token')}
        ],
        route: [
            {type: 'main', why: scenarioText('branches.browser_html_404')},
            {type: 'main', why: scenarioText('branches.api_json_404')}
        ]
    };
    const allEdges = [
        ['browser', 'caddy'],
        ['caddy', 'front'],
        ['front', 'gateway'],
        ['gateway', 'rabbit'],
        ['rabbit', 'content'],
        ['gateway', 'content'],
        ['content', 'cache'],
        ['gateway', 'captcha'],
        ['captcha', 'redis'],
        ['content', 'redis']
    ];
    const map = find('.explorer-map');
    const svg = find('.explorer-map svg');
    const packet = find('.explorer-packet');
    const node = id => find(`.explorer-node[data-id="${id}"]`);
    let visible = [];
    let edgeElements = new Map();
    const edgeKey = (source, destination) => [source, destination].sort().join(':');

    function geometry(source, destination) {
        const sourceBounds = node(source).getBoundingClientRect();
        const destinationBounds = node(destination).getBoundingClientRect();
        const mapBounds = map.getBoundingClientRect();
        const sourceCenter = {
            x: sourceBounds.left + sourceBounds.width / 2 - mapBounds.left,
            y: sourceBounds.top + sourceBounds.height / 2 - mapBounds.top
        };
        const destinationCenter = {
            x: destinationBounds.left + destinationBounds.width / 2 - mapBounds.left,
            y: destinationBounds.top + destinationBounds.height / 2 - mapBounds.top
        };
        const distanceX = destinationCenter.x - sourceCenter.x;
        const distanceY = destinationCenter.y - sourceCenter.y;
        const borderOffset = (width, height) => Math.min(
            Math.abs(distanceX) < 0.01 ? Infinity : width / 2 / Math.abs(distanceX),
            Math.abs(distanceY) < 0.01 ? Infinity : height / 2 / Math.abs(distanceY)
        );
        const sourceOffset = borderOffset(sourceBounds.width, sourceBounds.height);
        const destinationOffset = borderOffset(destinationBounds.width, destinationBounds.height);

        return [
            {
                x: sourceCenter.x + distanceX * sourceOffset,
                y: sourceCenter.y + distanceY * sourceOffset
            },
            {
                x: destinationCenter.x - distanceX * destinationOffset,
                y: destinationCenter.y - distanceY * destinationOffset
            }
        ];
    }

    function layout() {
        visible = view === 'all'
            ? ['browser', 'caddy', 'front', 'gateway', 'rabbit', 'content', 'cache', 'captcha', 'redis']
            : endpoints[view].nodes;
        const focused = view !== 'all';
        const columns = map.clientWidth <= 560 ? 1 : map.clientWidth <= 850 ? 3 : 4;
        map.classList.toggle('is-focused', focused);
        map.style.gridTemplateColumns = focused ? `repeat(${columns},minmax(0,1fr))` : '';
        map.style.gridTemplateRows = focused ? `repeat(${Math.ceil(visible.length / columns)},minmax(88px,auto))` : '';
        findAll('.explorer-node').forEach(el => {
            const position = visible.indexOf(el.dataset.id);
            el.hidden = position < 0;
            el.disabled = focused;
            el.classList.remove('is-hit');
            if (focused && position >= 0) {
                const row = Math.floor(position / columns);
                const column = position % columns;
                el.style.gridColumn = String(row % 2 === 0 ? column + 1 : columns - column);
                el.style.gridRow = String(row + 1);
                el.style.minHeight = '88px';
            } else {
                el.style.gridColumn = '';
                el.style.gridRow = '';
                el.style.minHeight = '';
            }
        });
        requestAnimationFrame(drawEdges);
    }

    function drawEdges() {
        svg.setAttribute('viewBox', `0 0 ${map.clientWidth} ${map.clientHeight}`);
        svg.replaceChildren();
        edgeElements = new Map();
        const edges = view === 'all'
            ? allEdges
            : visible.slice(0, -1).map((id, index) => [id, visible[index + 1]]);
        edges.forEach(([source, destination]) => {
            const [start, end] = geometry(source, destination);
            const path = document.createElementNS('http://www.w3.org/2000/svg', 'path');
            path.setAttribute('d', `M ${start.x} ${start.y} L ${end.x} ${end.y}`);
            svg.append(path);
            edgeElements.set(edgeKey(source, destination), path);
        });
    }

    function renderViews() {
        const holder = find('.explorer-views');
        holder.replaceChildren();
        ['all', ...Object.keys(endpoints)].forEach(id => {
            const button = document.createElement('button');
            button.type = 'button';
            button.dataset.view = id;
            button.setAttribute('aria-pressed', String(view === id));
            const name = document.createElement('span');
            name.textContent = id === 'all' ? ui.all : viewNames[id];
            const code = document.createElement('small');
            code.textContent = id === 'all' ? ui.allSubtitle : endpoints[id].label;
            button.append(name, code);
            button.addEventListener('click', () => selectView(id));
            holder.append(button);
        });
    }

    function renderInspector() {
        const item = view === 'all' && selectedComponent ? componentCopy[selectedComponent] : null;
        find('.explorer-component-detail').hidden = view !== 'all';
        find('[data-component-label]').textContent = item ? ui.component : ui.discover;
        find('[data-component-title]').textContent = item ? item.title : ui.discoverTitle;
        find('[data-component-description]').textContent = item ? item.description : ui.discoverText;
        find('[data-hint]').hidden = view !== 'all';
        findAll('.explorer-node').forEach(button => {
            button.setAttribute('aria-pressed', String(button.dataset.id === selectedComponent));
        });
    }

    function getSteps() {
        return view === 'all' ? [] : endpoints[view].build(variant);
    }

    function renderChoice() {
        const choicePanel = find('.explorer-choice');
        choicePanel.replaceChildren();
        if (view === 'all') return;
        const meta = variantMeta[view];
        const makeGroup = (type, heading, caption) => {
            const group = document.createElement('section');
            group.className = 'explorer-choice-group' + (type === 'extra' ? ' is-extra' : '');
            const title = document.createElement('strong');
            title.textContent = heading;
            const text = document.createElement('p');
            text.textContent = caption;
            const buttons = document.createElement('div');
            buttons.className = 'explorer-choice-buttons';
            meta.forEach((entry, index) => {
                if (entry.type !== type) return;
                const button = document.createElement('button');
                button.type = 'button';
                button.textContent = endpoints[view].variants[index];
                button.setAttribute('aria-pressed', String(index === variant));
                button.disabled = stepIndex > 0 || busy;
                button.addEventListener('click', () => {
                    if (stepIndex || busy) return;
                    variant = index;
                    renderFlow();
                });
                buttons.append(button);
            });
            group.append(title, text, buttons);
            choicePanel.append(group);
        };
        makeGroup('main', view === 'route' ? ui.routeHeading : ui.mainHeading, ui.mainCaption);
        const hasExtra = meta.some(item => item.type === 'extra');
        choicePanel.classList.toggle('is-single', !hasExtra);
        if (hasExtra) makeGroup('extra', ui.extraHeading, ui.extraCaption);
        const optional = meta[variant].type === 'extra';
        const kind = find('[data-flow-kind]');
        kind.textContent = optional ? ui.extraKind : view === 'route' ? ui.routeKind : ui.mainKind;
        kind.classList.toggle('is-optional', optional);
        const note = find('[data-variant-explain]');
        note.textContent = meta[variant].why;
        note.classList.toggle('is-optional', optional);
    }

    function renderStepper() {
        const stepper = find('.explorer-stepper');
        stepper.replaceChildren();
        getSteps().forEach((step, index) => {
            const item = document.createElement('span');
            item.className = 'explorer-step';
            const number = document.createElement('b');
            number.textContent = String(index + 1).padStart(2, '0');
            item.append(number, document.createTextNode(step.label));
            if (index === stepIndex && !busy) item.classList.add('is-next');
            if (index < stepIndex) item.classList.add('is-done');
            stepper.append(item);
        });
    }

    function renderMobileControls() {
        const holder = find('[data-mobile-run]');
        holder.hidden = view === 'all';
        if (view === 'all') return;
        const steps = getSteps();
        find('[data-mobile-progress]').textContent = `${stepIndex} / ${steps.length} ${ui.step}`;
        find('[data-mobile-result]').textContent = busy
            ? ui.moving
            : stepIndex
            ? steps[stepIndex - 1].description
            : ui.ready;
        const next = find('[data-mobile-next]');
        next.textContent = busy
            ? ui.working
            : stepIndex >= steps.length
            ? ui.done
            : `${ui.nextStep}: ${steps[stepIndex].label}`;
        next.disabled = busy || stepIndex >= steps.length;
        const pause = find('[data-mobile-pause]');
        pause.textContent = paused ? ui.resume : ui.pause;
        pause.disabled = !busy;
        const reset = find('[data-mobile-reset]');
        reset.textContent = ui.reset;
        reset.disabled = stepIndex === 0 && !busy;
    }

    function initialDetail() {
        find('[data-detail-title]').textContent = ui.ready;
        find('[data-detail-description]').textContent = endpoints[view].description;
        ['request', 'response', 'state', 'decision'].forEach(name => find(`[data-code="${name}"]`).textContent = '—');
        find('.explorer-feedback').classList.remove('is-error', 'is-updated');
        find('[data-feedback]').textContent = ui.ready;
    }

    function renderFlow() {
        const all = view === 'all';
        find('.explorer-work').hidden = all;
        find('[data-flow-options]').hidden = all;
        find('[data-map-title]').textContent = all ? ui.all : endpoints[view].title;
        find('[data-map-description]').textContent = all ? ui.mapDescription : endpoints[view].description;
        find('[data-map-status]').textContent = `${all ? '9' : endpoints[view].nodes.length} / 9 ${ui.components}`;
        renderMobileControls();
        if (all) return;
        find('[data-flow-title]').textContent = viewNames[view];
        find('[data-flow-endpoint]').textContent = endpoints[view].label;
        find('[data-progress]').textContent = `${stepIndex} / ${getSteps().length} ${ui.step}`;
        renderChoice();
        renderStepper();
        if (stepIndex === 0 && !busy) {
            initialDetail();
            find('[data-nonce]').hidden = true;
        }
    }

    function selectView(id) {
        if (view === id) return;
        generation++;
        busy = false;
        paused = false;
        packet.hidden = true;
        view = id;
        variant = 0;
        stepIndex = 0;
        selectedComponent = null;
        renderViews();
        layout();
        renderInspector();
        renderFlow();
        if (id !== 'all') requestAnimationFrame(() => find('[data-map-title]').scrollIntoView({
            behavior: 'smooth',
            block: 'start'
        }));
    }

    function delay(ms, run) {
        return new Promise(resolve => {
            let elapsed = 0, last = 0;

            function frame(now) {
                if (run !== generation) {
                    resolve(false);
                    return;
                }
                if (!last) last = now;
                if (!paused) elapsed += now - last;
                last = now;
                if (elapsed >= ms) {
                    resolve(true);
                    return;
                }
                requestAnimationFrame(frame);
            }

            requestAnimationFrame(frame);
        });
    }

    async function animateNonce(run) {
        const attempts = [
            [0, 'dd4fd2dfeaa3a31f'],
            [71, 'd938284ed0b74f09'],
            [329, 'f9420618bb6a52'],
            [1800, '9a1ed77daa464506']
        ];
        if (variant !== 3) attempts.push([19228, '0000045a987ec2c6']);
        const panel = find('[data-nonce]');
        panel.hidden = false;
        panel.classList.remove('is-solved');
        for (let i = 0; i < attempts.length; i++) {
            if (run !== generation) return false;
            find('[data-nonce-try]').textContent = attempts[i][0];
            find('[data-nonce-hash]').textContent = attempts[i][1] + '…';
            find('[data-nonce-progress]').style.width = `${((i + 1) / attempts.length) * 100}%`;
            panel.classList.remove('is-ticking');
            void panel.offsetWidth;
            panel.classList.add('is-ticking');
            if (!await delay(matchMedia('(prefers-reduced-motion: reduce)').matches ? 20 : 360, run)) return false;
        }
        panel.classList.toggle('is-solved', variant !== 3);
        return true;
    }

    function animateEdge(source, destination, kind, run) {
        return new Promise(resolve => {
            if (run !== generation) {
                resolve(false);
                return;
            }
            if (map.clientWidth <= 560) {
                const target = node(destination).getBoundingClientRect();
                if (target.top < 170 || target.bottom > window.innerHeight - 40) node(destination).scrollIntoView({
                    behavior: matchMedia('(prefers-reduced-motion: reduce)').matches ? 'auto' : 'smooth',
                    block: 'center'
                });
            }
            const [from, to] = geometry(source, destination);
            const edge = edgeElements.get(edgeKey(source, destination));
            if (edge) {
                edge.classList.add('is-travel');
                edge.classList.toggle('is-response', kind === 'response');
            }
            if (matchMedia('(prefers-reduced-motion: reduce)').matches) {
                node(destination).classList.add('is-hit');
                if (edge) edge.classList.remove('is-travel', 'is-response');
                resolve(true);
                return;
            }
            packet.hidden = false;
            packet.classList.toggle('is-response', kind === 'response');
            const duration = 460;
            let elapsed = 0, last = 0;

            function frame(now) {
                if (run !== generation) {
                    packet.hidden = true;
                    resolve(false);
                    return;
                }
                if (!last) last = now;
                if (!paused) elapsed += now - last;
                last = now;
                const progress = Math.min(1, elapsed / duration);
                packet.style.left = `${from.x + (to.x - from.x) * progress}px`;
                packet.style.top = `${from.y + (to.y - from.y) * progress}px`;
                if (progress < 1) {
                    requestAnimationFrame(frame);
                } else {
                    packet.hidden = true;
                    node(destination).classList.add('is-hit');
                    if (edge) edge.classList.remove('is-travel', 'is-response');
                    resolve(true);
                }
            }

            requestAnimationFrame(frame);
        });
    }

    async function runStep(requestedStepIndex) {
        if (busy || requestedStepIndex !== stepIndex || view === 'all') return;
        const step = getSteps()[requestedStepIndex];
        const run = ++generation;
        busy = true;
        paused = false;
        renderFlow();
        find('[data-feedback]').textContent = ui.moving;
        findAll('.explorer-node').forEach(el => el.classList.remove('is-hit'));
        if (!(view === 'captcha' && requestedStepIndex === 0)) find('[data-nonce]').hidden = true;
        if (step.path.length) {
            node(step.path[0]).classList.add('is-hit');
            for (let pathIndex = 0; pathIndex < step.path.length - 1; pathIndex++) {
                const completed = await animateEdge(
                    step.path[pathIndex],
                    step.path[pathIndex + 1],
                    step.kind,
                    run
                );
                if (!completed) return;
            }
        } else {
            if (view === 'captcha' && requestedStepIndex === 0) {
                if (!await animateNonce(run)) return;
            } else if (!await delay(470, run)) return;
        }
        if (run !== generation) return;
        busy = false;
        paused = false;
        stepIndex++;
        const feedback = find('.explorer-feedback');
        feedback.classList.toggle('is-error', /^(4|5)\d\d/.test(step.response));
        feedback.classList.remove('is-updated');
        void feedback.offsetWidth;
        feedback.classList.add('is-updated');
        find('[data-feedback]').textContent = step.description;
        find('[data-detail-title]').textContent = step.label;
        find('[data-detail-description]').textContent = step.description;
        ['request', 'response', 'state', 'decision'].forEach(name => {
            find(`[data-code="${name}"]`).textContent = step[name];
        });
        find('[data-progress]').textContent = `${stepIndex} / ${getSteps().length} ${ui.step}`;
        renderChoice();
        renderStepper();
        renderMobileControls();
        if (stepIndex === getSteps().length) {
            setTimeout(() => {
                if (run === generation && stepIndex === getSteps().length) {
                    find('[data-progress]').textContent = ui.done;
                }
            }, 700);
        }
    }

    function reset() {
        generation++;
        busy = false;
        paused = false;
        stepIndex = 0;
        packet.hidden = true;
        findAll('.explorer-node').forEach(el => el.classList.remove('is-hit'));
        edgeElements.forEach(el => el.classList.remove('is-travel', 'is-response'));
        renderFlow();
    }

    renderViews();
    renderInspector();
    renderFlow();
    layout();
    findAll('.explorer-node').forEach(button => button.addEventListener('click', () => {
        if (view !== 'all') return;
        selectedComponent = button.dataset.id;
        renderInspector();
        const detail = find('.explorer-component-detail');
        detail.classList.remove('is-updated');
        void detail.offsetWidth;
        detail.classList.add('is-updated');
    }));
    find('[data-mobile-reset]').addEventListener('click', reset);
    find('[data-mobile-next]').addEventListener('click', () => runStep(stepIndex));
    const togglePause = () => {
        if (!busy) return;
        paused = !paused;
        renderMobileControls();
    };
    find('[data-mobile-pause]').addEventListener('click', togglePause);
    const refreshMap = () => {
        if (view !== 'all') layout(); else requestAnimationFrame(drawEdges);
    };
    if (typeof ResizeObserver === 'function') new ResizeObserver(refreshMap).observe(map);
    else window.addEventListener('resize', refreshMap);
    window.addEventListener('pageshow', refreshMap);
})();
