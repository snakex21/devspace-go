# MCP WebCoder Auto-Approve — Tampermonkey Script

Automatycznie klika "Zawsze zezwalaj" / "Allow" / "Połącz" dla połączeń MCP w ChatGPT.
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
| Dialog narzędzia MCP | Auto-klika "Zawsze zezwalaj" |
| Modal bez potwierdzeń | Auto-potwierdza "Zawsze zezwalaj" |
| Ekran połączenia | Auto-klika "Połącz" |
| Strona OAuth MCP WebCoder | Auto-submituje jeśli hasło już wpisane |
| Callback OAuth (pusta strona) | Auto-zamyka okno |

## Konfiguracja hasła (opcjonalne)

Jeśli chcesz, żeby skrypt automatycznie wpisywał hasło OAuth:

1. Otwórz skrypt w Tampermonkey
2. Znajdź obsługę `passwordField`
3. Przed sprawdzeniem pustego hasła dodaj:
```js
passwordField.value = 'twoje-haslo-tutaj';
```

## Wyłączenie

Kliknij Tampermonkey → przełącznik przy "MCP WebCoder Auto-Approve for ChatGPT" → OFF.
