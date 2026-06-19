// ==UserScript==
// @name         MCP WebCoder Auto-Approve for ChatGPT
// @namespace    mcp-webcoder
// @version      1.1.0
// @description  Automatycznie klika "Zawsze zezwalaj" / "Połącz" dla MCP WebCoder w ChatGPT.
// @author       MCP WebCoder
// @match        https://chatgpt.com/*
// @match        https://chat.openai.com/*
// @icon         https://www.google.com/s2/favicons?sz=64&domain=chatgpt.com
// @grant        none
// @run-at       document-end
// ==/UserScript==

(function () {
    'use strict';

    const LOG_PREFIX = '[MCP WebCoder Auto-Approve]';
    const CLICK_DELAY_MS = 180;
    const SCAN_INTERVAL_MS = 350;

    const clickedButtons = new WeakSet();

    function log(...args) {
        console.log(LOG_PREFIX, ...args);
    }

    function norm(s) {
        return (s || '')
            .replace(/\s+/g, ' ')
            .trim()
            .toLowerCase();
    }

    function isVisible(el) {
        if (!el || !(el instanceof Element)) return false;
        const style = window.getComputedStyle(el);
        if (style.display === 'none' || style.visibility === 'hidden' || style.opacity === '0') return false;
        const rect = el.getBoundingClientRect();
        return rect.width > 0 && rect.height > 0;
    }

    function buttonText(btn) {
        return norm(btn.textContent || btn.getAttribute('aria-label') || btn.getAttribute('title') || '');
    }

    function buttonScore(btn) {
        let score = 0;
        const cls = btn.className || '';
        if (String(cls).includes('btn-primary')) score += 100;
        if (String(cls).includes('btn-large')) score += 20;
        if (btn.matches('button')) score += 10;
        return score;
    }

    function allButtons(scope = document) {
        return [...scope.querySelectorAll('button, [role="button"], a.btn, .btn')]
            .filter((el) => isVisible(el) && !el.disabled && el.getAttribute('aria-disabled') !== 'true');
    }

    function findButton(scope, labels, opts = {}) {
        const normalizedLabels = labels.map(norm);
        const candidates = [];

        for (const btn of allButtons(scope)) {
            const text = buttonText(btn);
            if (!text) continue;

            for (const label of normalizedLabels) {
                const exact = text === label;
                const contains = text.includes(label);
                if ((opts.exactOnly && exact) || (!opts.exactOnly && (exact || contains))) {
                    candidates.push({ btn, exact, score: buttonScore(btn) + (exact ? 1000 : 0) });
                    break;
                }
            }
        }

        candidates.sort((a, b) => b.score - a.score);
        return candidates[0]?.btn || null;
    }

    function safeClick(btn, reason) {
        if (!btn || clickedButtons.has(btn)) return false;
        clickedButtons.add(btn);
        log('click:', reason, '=>', (btn.textContent || '').trim());
        setTimeout(() => {
            try {
                btn.click();
            } catch (e) {
                log('click failed:', e);
            }
        }, CLICK_DELAY_MS);
        return true;
    }

    function dialogText(el) {
        return norm(el?.textContent || '');
    }

    function isMcpPermissionText(text) {
        return (
            text.includes('mcp') ||
            text.includes('model context protocol') ||
            text.includes('mcp webcoder') ||
            text.includes('devspace') ||
            text.includes('zezwolić chatgpt na użycie aplikacji') ||
            text.includes('let chatgpt use') ||
            text.includes('without confirmation') ||
            text.includes('bez potwierdzenia') ||
            text.includes('workspaceidentifiers') ||
            text.includes('udostępniane dane') ||
            text.includes('shared data') ||
            text.includes('custom mcp') ||
            text.includes('niestandardowe serwery mcp')
        );
    }

    function isConnectText(text) {
        return (
            text.includes('connect') ||
            text.includes('połącz') ||
            text.includes('authorization') ||
            text.includes('autoryz') ||
            text.includes('mcp webcoder') ||
            text.includes('devspace')
        );
    }

    function findRelevantScopes() {
        const selectors = [
            '[role="dialog"]',
            '[role="alertdialog"]',
            '[data-testid="tool-action-buttons"]',
            '.modal',
            '.dialog',
            'main',
            'body',
        ];

        const scopes = [];
        for (const selector of selectors) {
            for (const el of document.querySelectorAll(selector)) {
                if (isVisible(el)) scopes.push(el);
            }
        }
        return scopes;
    }

    function approveToolPermission(scope, text) {
        if (!isMcpPermissionText(text)) return false;

        // First ChatGPT tool-use card: prefer persistent permission.
        const alwaysSecondary = findButton(scope, [
            'zawsze zezwalaj',
            'always allow',
            'always approve',
            'allow always',
        ]);
        if (safeClick(alwaysSecondary, 'always allow / persistent permission')) return true;

        // Fallback if ChatGPT only shows one-time allow.
        const allowOnce = findButton(scope, [
            'zezwól tym razem',
            'allow this time',
            'allow once',
            'zezwól',
            'allow',
            'approve',
        ]);
        if (safeClick(allowOnce, 'allow this time fallback')) return true;

        return false;
    }

    function approveAlwaysConfirmation(scope, text) {
        const isConfirm =
            text.includes('without confirmation') ||
            text.includes('bez potwierdzenia') ||
            text.includes("won't ask before") ||
            text.includes('nie będzie pytać') ||
            text.includes('elevated risk') ||
            text.includes('podwyższonym ryzykiem');

        if (!isConfirm) return false;

        const confirm = findButton(scope, [
            'zawsze zezwalaj',
            'always allow',
            'allow always',
        ]);
        return safeClick(confirm, 'confirm always allow modal');
    }

    function clickConnect(scope, text) {
        if (!isConnectText(text) && !document.querySelector('button.btn-primary')) return false;

        const connect = findButton(scope, [
            'połącz',
            'connect',
            'continue',
            'kontynuuj',
            'authorize',
            'autoryzuj',
        ]);
        return safeClick(connect, 'connect / authorize');
    }

    function scan() {
        // Most-specific scopes first so modal buttons win over background cards.
        const scopes = findRelevantScopes();
        for (const scope of scopes) {
            const text = dialogText(scope);
            if (!text) continue;

            if (approveAlwaysConfirmation(scope, text)) return;
            if (approveToolPermission(scope, text)) return;
            if (clickConnect(scope, text)) return;
        }
    }

    function handleOAuthPopup() {
        const title = norm(document.title);
        const text = norm(document.body?.textContent || '');
        if (!window.location.pathname.includes('/authorize') && !title.includes('authorization') && !text.includes('authorization')) {
            return;
        }

        const passwordField = document.querySelector('input[type="password"]');
        const submitBtn = findButton(document, ['authorize', 'autoryzuj', 'submit', 'wyślij', 'sign in', 'zaloguj', 'connect', 'połącz']);

        if (!submitBtn) return;
        if (passwordField && passwordField.value.length === 0) return;

        safeClick(submitBtn, 'oauth authorize');
    }

    function handlePostOAuthBlank() {
        if (window.location.search.includes('code=') &&
            (document.body.innerText.trim() === '' || document.body.innerText.length < 100)) {
            log('OAuth callback detected, closing popup...');
            window.close();
            setTimeout(() => {
                document.body.innerHTML = '<h2 style="text-align:center;margin-top:40px;font-family:sans-serif">✅ Authorization complete. You can close this window.</h2>';
            }, 500);
        }
    }

    function run() {
        scan();
        handleOAuthPopup();
        handlePostOAuthBlank();
    }

    function startObserver() {
        const observer = new MutationObserver(run);
        observer.observe(document.body, {
            childList: true,
            subtree: true,
            attributes: true,
            attributeFilter: ['class', 'style', 'aria-hidden', 'open', 'disabled', 'aria-expanded'],
        });

        setInterval(run, SCAN_INTERVAL_MS);
        run();
        log('watching for MCP permission dialogs ✅');
    }

    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', startObserver);
    } else {
        startObserver();
    }
})();
