# DevSpace Auto-Approve — Tampermonkey Script

Automatycznie klika "Zezwól" / "Allow" dla połączeń MCP w ChatGPT.
Koniec z ręcznym potwierdzaniem za każdym razem!

## Instalacja

1. Zainstaluj [Tampermonkey](https://www.tampermonkey.net/) dla swojej przeglądarki:
   - Chrome: [Tampermonkey](https://chrome.google.com/webstore/detail/tampermonkey/dhdgffkkebhmkfjojejmpbldmpobfkfo)
   - Firefox: [Tampermonkey](https://addons.mozilla.org/firefox/addon/tampermonkey/)
   - Edge: [Tampermonkey](https://microsoftedge.microsoft.com/addons/detail/tampermonkey/iikmkjmpaadaobahmlepeloendndfphd)

2. Kliknij ikonę Tampermonkey → **Create a new script...**

3. Wklej zawartość `scripts/userscripts/devspace-auto-approve.user.js`

4. **Ctrl+S** → gotowe!

## Co robi

| Akcja | Efekt |
|---|---|
| Dialog MCP "Connect to..." | Auto-klika "Allow" / "Zezwól" |
| Checkbox "I understand..." | Auto-zaznacza |
| Strona OAuth DevSpace | Auto-submituje jeśli hasło już wpisane |
| Callback OAuth (pusta strona) | Auto-zamyka okno |

## Konfiguracja hasła (opcjonalne)

Jeśli chcesz, żeby skrypt automatycznie wpisywał hasło OAuth:

1. Otwórz skrypt w Tampermonkey
2. Znajdź linię z komentarzem `// For now, just click submit`
3. Przed nią dodaj:
```js
passwordField.value = 'twoje-haslo-tutaj';
```

## Wyłączenie

Kliknij Tampermonkey → przełącznik przy "DevSpace MCP Auto-Approve" → OFF.
