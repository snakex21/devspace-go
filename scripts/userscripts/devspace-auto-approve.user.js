// ==UserScript==
// @name         DevSpace MCP Auto-Approve
// @namespace    devspace.go
// @version      1.0
// @description  Automatycznie klika "Allow" / "Zezwól" dla połączeń MCP w ChatGPT. Koniec z ręcznym potwierdzaniem!
// @author       DevSpace
// @match        https://chatgpt.com/*
// @match        https://chat.openai.com/*
// @icon         https://www.google.com/s2/favicons?sz=64&domain=chatgpt.com
// @grant        none
// @run-at       document-end
// ==/UserScript==

(function() {
    'use strict';

    const LOG_PREFIX = '[DevSpace Auto-Approve]';
    let approvedMcpUrls = new Set();
    let processedDialogs = new WeakSet();

    // Load previously approved URLs from localStorage
    try {
        const saved = localStorage.getItem('devspace_approved_mcp_urls');
        if (saved) {
            approvedMcpUrls = new Set(JSON.parse(saved));
            console.log(LOG_PREFIX, 'Loaded', approvedMcpUrls.size, 'pre-approved MCP URLs');
        }
    } catch(e) {}

    function saveApprovedUrls() {
        try {
            localStorage.setItem('devspace_approved_mcp_urls', JSON.stringify([...approvedMcpUrls]));
        } catch(e) {}
    }

    // --- Find and click buttons by text ---
    function findButton(texts) {
        const all = document.querySelectorAll('button, [role="button"], a.btn, .btn');
        for (const el of all) {
            const t = (el.textContent || '').trim().toLowerCase();
            for (const txt of texts) {
                if (t.includes(txt.toLowerCase())) {
                    return el;
                }
            }
        }
        return null;
    }

    // --- Auto-approve MCP connection dialogs ---
    function approveMcpDialog() {
        // Check for the MCP connection modal/dialog
        const dialogs = document.querySelectorAll('[role="dialog"], [role="alertdialog"], .modal, .dialog, [data-testid="mcp-connection-dialog"]');
        
        for (const dialog of dialogs) {
            if (processedDialogs.has(dialog)) continue;
            
            const text = (dialog.textContent || '').toLowerCase();
            
            // Is this an MCP-related dialog?
            const isMcpDialog = 
                text.includes('mcp') ||
                text.includes('model context protocol') ||
                text.includes('connect to') ||
                text.includes('połączenie') ||
                text.includes('serwer mcp') ||
                text.includes('custom mcp') ||
                text.includes('niestandardowe serwery mcp');
            
            if (!isMcpDialog) continue;

            console.log(LOG_PREFIX, 'Found MCP dialog, auto-approving...');
            processedDialogs.add(dialog);

            // Click "I understand and want to continue" checkbox if present
            const understand = findButton(['understand', 'rozumiem', 'i want to continue', 'chcę kontynuować']);
            if (understand) {
                console.log(LOG_PREFIX, 'Clicking "I understand"...');
                understand.click();
            }

            // Wait a tiny bit then click Approve/Allow/Connect
            setTimeout(() => {
                const approve = findButton([
                    'allow', 'connect', 'approve', 'accept', 'continue',
                    'zezwól', 'połącz', 'zatwierdź', 'akceptuj', 'kontynuuj',
                    'add', 'dodaj', 'save', 'zapisz'
                ]);
                if (approve && !approve.disabled) {
                    console.log(LOG_PREFIX, 'Clicking APPROVE!');
                    approve.click();
                    
                    // Try to extract the MCP URL being approved
                    const urlMatch = text.match(/https?:\/\/[^\s"]+/);
                    if (urlMatch) {
                        approvedMcpUrls.add(urlMatch[0]);
                        saveApprovedUrls();
                    }
                }
            }, 500);
        }
    }

    // --- Auto-handle OAuth redirect (the browser popup) ---
    function handleOAuthPopup() {
        // Check if we're on the DevSpace authorization page
        if (window.location.pathname === '/authorize' || 
            document.title.includes('DevSpace') ||
            document.title.includes('Authorization')) {
            
            console.log(LOG_PREFIX, 'Detected DevSpace OAuth page, auto-submitting...');
            
            // Find password field and submit button
            const passwordField = document.querySelector('input[type="password"]');
            const submitBtn = findButton(['authorize', 'autoryzuj', 'submit', 'wyślij', 'sign in', 'zaloguj']);
            
            if (passwordField && submitBtn) {
                // If we have a saved password, fill it (optional - user must set this)
                // For now, just click submit if the field is already filled
                setTimeout(() => {
                    if (passwordField.value.length > 0) {
                        console.log(LOG_PREFIX, 'Password detected, submitting...');
                        submitBtn.click();
                    } else {
                        console.log(LOG_PREFIX, 'Password field empty — fill it manually or save it in the script settings.');
                    }
                }, 1000);
            } else if (submitBtn) {
                // No password needed (auto-approve mode)
                setTimeout(() => {
                    console.log(LOG_PREFIX, 'No password required, submitting...');
                    submitBtn.click();
                }, 500);
            }
        }
    }

    // --- Main observer ---
    function startObserver() {
        // Watch for new dialogs/modals appearing
        const observer = new MutationObserver(() => {
            approveMcpDialog();
            handleOAuthPopup();
        });

        observer.observe(document.body, {
            childList: true,
            subtree: true,
            attributes: true,
            attributeFilter: ['class', 'style', 'aria-hidden', 'open']
        });

        // Also run immediately
        approveMcpDialog();
        handleOAuthPopup();

        console.log(LOG_PREFIX, 'Watching for MCP dialogs... ✅');
    }

    // --- OAuth redirect auto-close ---
    // When the OAuth popup redirects back, it might leave a blank page.
    // We detect this and close it.
    function handlePostOAuthBlank() {
        if (window.location.search.includes('code=') && 
            (document.body.innerText.trim() === '' || document.body.innerText.length < 100)) {
            console.log(LOG_PREFIX, 'OAuth callback detected, closing popup...');
            window.close();
            // If window.close() doesn't work (not opened by script), show a message
            setTimeout(() => {
                document.body.innerHTML = '<h2 style="text-align:center;margin-top:40px;font-family:sans-serif">✅ Authorization complete. You can close this window.</h2>';
            }, 500);
        }
    }

    // --- Start ---
    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', () => {
            startObserver();
            handlePostOAuthBlank();
        });
    } else {
        startObserver();
        handlePostOAuthBlank();
    }

})();
