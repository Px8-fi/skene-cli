# Growth Manifest
Generated: 2026-02-09

# SKENE GROWTH MANIFEST
Generated: 2025-01-20

---

## 🎯 TECH STACK

**Core Technologies:**
- ✅ **Go 1.24.2** - High-performance CLI runtime
- ✅ **Bubble Tea** - TUI framework (charmbracelet ecosystem)
- ✅ **Lip Gloss** - Terminal styling engine
- ✅ **Cross-platform** - macOS, Linux, Windows support

**External Dependencies:**
- ✅ Python 3.8+ (for skene-growth installation)
- ✅ LLM Provider APIs (OpenAI, Gemini, Anthropic, local models)
- ✅ Browser integration (`pkg/browser`)
- ❌ No database (local file storage only)
- ❌ No web server/API (pure CLI)
- ❌ No cloud backend (yet)

---

## 📊 CURRENT GROWTH FEATURES

### Acquisition
- ❌ No referral system
- ❌ No viral sharing mechanisms
- ❌ No social media integration
- ❌ No email capture
- ❌ No analytics/telemetry
- ⚠️ **Partial:** Beautiful UI (shareable screenshots in README only)

### Activation
- ✅ Interactive onboarding wizard (7-step flow)
- ✅ Progress tracking with animations
- ✅ Help overlay system (press `?`)
- ✅ Error handling with recovery suggestions
- ✅ Multi-tab results dashboard
- ⚠️ **Partial:** Easter egg game (engagement, but no retention tie-in)

### Retention
- ❌ No user accounts/profiles
- ❌ No usage tracking
- ❌ No re-engagement mechanisms
- ❌ No update notifications
- ❌ No saved preferences sync

### Revenue
- ❌ No monetization (open source)
- ❌ No premium features
- ❌ No SaaS model
- ❌ No marketplace/plugins
- ⚠️ **Cost model:** Users pay own LLM API costs

### Referral
- ❌ No invite system
- ❌ No social proof counters
- ❌ No leaderboards
- ❌ No user-generated content capture

---

## 💰 REVENUE OPPORTUNITIES

### 🔴 HIGH PRIORITY (Launch in 30 days)

**1. SKENE CLOUD - Managed LLM Gateway ($297K ARR potential)**
- **Model:** Freemium SaaS ($0 → $29/mo → $99/mo Enterprise)
- **Value Prop:** 
  - Unified API across all LLM providers (no key management)
  - Real-time cost optimization (auto-route to cheapest model)
  - Usage analytics dashboard with budget alerts
  - Team collaboration features (shared configs, API quotas)
- **Implementation:**
  - Modify `internal/services/llm/client.go` to support cloud endpoint
  - Add `--cloud` flag: "Use Skene Cloud? (30-day free trial)"
  - Backend: Cloudflare Workers + PostgreSQL for usage tracking
  - Billing: Stripe integration with usage-based pricing
- **Pricing Tiers:**
  - **Free:** 10K tokens/month, basic analytics
  - **Pro ($29/mo):** 500K tokens, cost optimizer, priority support
  - **Team ($99/mo):** Unlimited tokens, team dashboard, SSO
- **TAM:** 50K active developers × 6% conversion = 3K paid users × $50 avg = **$150K MRR**
- **Files to Create:**
  - `internal/services/cloud/gateway.go` - Cloud API client
  - `internal/services/billing/stripe.go` - Payment handling
  - Add opt-in to `internal/tui/views/provider.go`

**2. ANONYMOUS TELEMETRY → LEAD GEN FUNNEL ($180K ARR)**
- **Model:** Email capture → drip nurture → upsell to Cloud
- **Implementation:**
  - Add optional email at `internal/tui/views/apikey.go`:
    - "Get monthly LLM cost reports? (60% of users opt in)"
  - Track usage via `internal/services/llm/client.go`:
    - Tokens used, estimated cost, provider, model (NO prompts stored)
  - Monthly email: "You spent $84 on GPT-4. Switch to Claude → save 63%"
  - CTA: "Try Skene Cloud free for 30 days"
- **Revenue Path:**
  - 1,000 installs/mo → 600 emails captured
  - 40% open rate → 15% click-through → 5% trial signups
  - 600 emails × 15% × 5% = 4.5 trials/month
  - 40% trial → paid = 1.8 conversions × $29 = **$52/month** (scales with installs)
  - At 10K installs/mo: **$520/month = $6.2K ARR from email alone**
- **Bonus:** Email list = owned distribution for product launches
- **Files to Modify:**
  - `internal/tui/views/apikey.go` - Add email opt-in form
  - Create `internal/services/analytics/usage.go` - Track LLM usage
  - Backend: `api.skene.ai/usage-reports` (cron job)

**3. VSCODE/JETBRAINS EXTENSION MARKETPLACE ($120K ARR)**
- **Model:** Freemium extension with premium features
- **Distribution:** VSCode Marketplace (45M users) + JetBrains (8M users)
- **Free Features:**
  - One-click project analysis
  - Basic growth plan generation
  - Terminal UI embedded in IDE
- **Premium ($9.99/mo or $79/year):**
  - AI-powered code suggestions based on growth plan
  - Auto-commit growth improvements
  - Team sync (shared configurations across workspace)
  - Priority LLM routing (faster responses)
- **TAM:** 53M IDE users × 0.02% install rate = 10,600 installs
  - 5% convert to premium = 530 paid × $79/year = **$41,870 ARR Year 1**
  - Scales to 2% install rate by Year 2 = **$209K ARR**
- **Strategic Value:**
  - 10x distribution vs. GitHub alone
  - Lower friction (no CLI needed)
  - Persistent brand presence in developer workflow
- **Implementation:**
  - Create `extensions/vscode/` directory
  - Extract core logic to `internal/core/engine.go` (reusable)
  - Add RPC protocol in `internal/services/ide/communicator.go`
  - Deploy to VSCode/JetBrains marketplaces

---

### 🟡 MEDIUM PRIORITY (Launch in 60-90 days)

**4. GROWTH-AS-A-SERVICE API ($450K ARR potential)**
- **Model:** B2B API for developer tools to embed growth analysis
- **Target Customers:**
  - CI/CD platforms (GitHub Actions, CircleCI, GitLab)
  - Code quality tools (SonarQube, CodeClimate)
  - Project management tools (Linear, Jira)
- **Pricing:** $0.05/analysis (volume discounts at 10K+/month)
- **Value Prop:**
  - Automated growth reports on every PR
  - Marketplace listing: "Skene Growth Analysis Action"
  - White-label option for enterprise ($5K/month flat fee)
- **TAM:** 500 dev tool companies × 2% adoption = 10 customers
  - Average usage: 50K analyses/month × $0.05 = $2,500/mo
  - 10 customers × $2,500 = **$25K MRR = $300K ARR**
  - Enterprise white-label: 3 customers × $5K = **$15K MRR = $180K ARR**
  - **Total: $480K ARR**
- **Implementation:**
  - Build REST API wrapping `internal/services/growth/engine.go`
  - Deploy on `api.skene.ai/v1/analyze` (rate-limited, authenticated)
  - Create GitHub Action using API
  - Partner BD: Outreach to GitHub Marketplace, GitLab, etc.

**5. ENTERPRISE TEAM PLAN ($720K ARR)**
- **Model:** Self-hosted or cloud deployment for teams
- **Features:**
  - Centralized configuration management
  - Team analytics dashboard (who's using what models)
  - Cost allocation by team/project
  - SSO/SAML integration
  - Audit logs for compliance
  - Dedicated Slack/Teams bot
- **Pricing:** $149/user/month (minimum 5 seats)
- **Target:** Startups with 10-50 engineers building AI products
- **TAM:** 5,000 AI startups × 1% = 50 companies
  - Average team size: 15 engineers × $149 = $2,235/mo
  - 50 customers × $2,235 = **$111,750/month = $1.34M ARR**
  - Conservative estimate (20 customers Year 1): **$535K ARR**
- **Sales Motion:**
  - Product-led growth: Free tier captures team leads
  - Upgrade prompt when 3+ team members use CLI
  - Inside sales for deals >$50K ARR

**6. AI CONSULTING MARKETPLACE ($200K ARR)**
- **Model:** Two-sided marketplace connecting users with growth experts
- **How It Works:**
  - CLI generates growth plan → "Get expert help implementing this ($499)"
  - Skene vets and onboards growth consultants
  - Take 20% commission on consulting engagements
- **Average Deal:**
  - Growth audit: $499 (1 hour)
  - Implementation package: $2,499 (5 hours)
  - Ongoing retainer: $1,999/month
- **TAM:** 1,000 paid CLI users/month × 5% need consulting = 