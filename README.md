# pp-amazon — Amazon repurchase CLI + Claude Code skill

A Go CLI + Playwright helper that lets a Claude Code agent reorder previously-purchased Amazon items, with explicit confirmation gates before any money moves.

**Repurchase-only**: refuses to add anything the user hasn't bought before. **Full checkout via headless Chromium** so Amazon doesn't see static-POST automation. **Multi-account** via profiles. **Real money — no sandbox.** Read [SKILL.md](SKILL.md) before you let an agent loose with this.

## What you get

- `amazon-pp-cli` — the Go binary
- `amazon-checkout.mjs` — Playwright helper for cart-show / add-to-cart / checkout / history-sync
- `SKILL.md` — the Claude Code skill file (drop into `~/.claude/skills/pp-amazon/`)

## Requirements

- macOS (tested on Apple Silicon) or Linux x86_64/arm64
- Go 1.26.3 or newer ([go.dev/dl](https://go.dev/dl/))
- Node 22+ and npm
- Chromium (Playwright will download one on first run; or supply your own via `PLAYWRIGHT_CHROMIUM_EXECUTABLE_PATH`)
- A logged-in amazon.com session in Safari/Chrome (for copying the Cookie header)

## Install (5 minutes)

```bash
# 1. Run the installer — builds the binary, installs the helper, drops SKILL.md.
./install.sh

# 2. Create a profile and paste your Amazon Cookie header.
amazon-pp-cli profiles add personal --label "Personal account"
amazon-pp-cli --profile personal auth paste
# (paste your Cookie header from DevTools, press Ctrl+D)

# 3. Verify the session reaches Amazon.
amazon-pp-cli --profile personal doctor
# Expected: amazon_reached: true

# 4. Sync your order history.
amazon-pp-cli --profile personal history sync
# Walks /your-orders for current year + previous 2 years.

# 5. Tell the CLI your default card last-4 (Amazon's cart page doesn't reliably expose it).
amazon-pp-cli --profile personal defaults set --card-last4 NNNN --card-label "Visa"

# 6. Try a dry-run add (no cart write).
amazon-pp-cli --profile personal add 'paper towels' --dry-run --json
```

The skill is now installed at `~/.claude/skills/pp-amazon/`. Any Claude Code agent running with that skill loaded can reorder via natural language.

## Getting the Cookie header

1. Open amazon.com in your browser (must be logged in)
2. Open DevTools (Cmd+Option+I on macOS)
3. Network tab → click any request → Headers → Request Headers → copy the entire `Cookie:` value
4. Paste into `amazon-pp-cli --profile <name> auth paste`

Cookies expire / get challenged periodically (Amazon enforces a 15-minute "max auth age" on order history and checkout). If you start seeing exit 9 (`manual_required`) errors, re-paste your Cookie header.

## Multi-account setup

```bash
amazon-pp-cli profiles add personal --label "Personal"
amazon-pp-cli profiles add work --label "Work (Amazon Business)"

amazon-pp-cli --profile personal auth paste   # paste personal cookies
amazon-pp-cli --profile work auth paste       # paste work cookies (sign out + sign in first!)

amazon-pp-cli --profile personal history sync
amazon-pp-cli --profile work history sync
```

**Important**: Amazon's `session-id` cookie is browser-bound, not account-bound. If you paste cookies from a Safari session that was recently in a different account, the dumper may return the wrong account's orders. After switching accounts, refresh the page so all cookies update, then paste.

## Skill flow (what the agent does)

When you say "order more paper towels":

1. Agent asks which profile (if unclear)
2. Agent runs `add 'paper towels' --dry-run` → shows you the matched item + last purchase date
3. You say "yep"
4. Agent runs `add 'paper towels'` → drives Playwright to the product page, clicks Add to Cart, verifies the item landed in the active cart by qty delta
5. Agent runs `cart show` → lists EVERY line item (warns about pre-existing items)
6. Agent asks for placement confirmation, including the card last-4
7. You say "yes, place it"
8. Agent runs `checkout --yes` → drives Playwright through the full checkout flow, returns the order ID

If any step hits a CAPTCHA or sign-in gate, the helper exits 9 with a deeplink. The agent hands you the deeplink — you tap it in Safari, finish the action there.

## Honest caveats

- **No upstream API.** This scrapes amazon.com via Playwright. Amazon changes their DOM regularly. When selectors break, the helper now exits 7 with the actual stderr reason (e.g. "no proceed-to-checkout button found"). Fix the selector and rebuild.
- **CAPTCHAs are unavoidable.** Amazon's anti-bot is aggressive on the checkout endpoint. The helper's stealth shims help but don't eliminate. The graceful-fallback is the exit-9 deeplink path.
- **Loose matches are a real risk.** `match_quality=loose` means only some query tokens hit the title. The CLI requires `--allow-loose` to commit, and the skill says "ask the user to confirm the title back" before passing it. Don't let an agent rush past this.
- **History is a SQLite snapshot.** Run `history sync` regularly so new orders are visible to repurchase resolution.

## Repo layout

```
pp-amazon-skill/
├── SKILL.md                 # The Claude Code skill (drops into ~/.claude/skills/pp-amazon/)
├── amazon-checkout.mjs      # Playwright helper (cart-show, add-to-cart, checkout, history-sync)
├── README.md                # This file
├── install.sh               # Build + install script
└── package.json             # npm metadata so install.sh can `npm install playwright`
```

The Go CLI source itself is at `github.com/<your-fork>/printing-press-library/library/commerce/amazon/`. The installer clones / pulls it, builds with `go install`, and drops the binary at `$GOPATH/bin/amazon-pp-cli`.

## Provenance

Built on top of the [printing-press](https://github.com/mvanhorn/cli-printing-press) library conventions, hand-coded (not generated) because Amazon shopping has no public API spec. Structure mirrors `library/commerce/instacart/` from the same library.

## License

Personal use, no warranty. If Amazon's lawyers come knocking, that's between you and them.
