# Growth Plan
Generated: 2026-02-10

## EXECUTIVE SUMMARY

Skene CLI is a developer onboarding tool disguised as a codebase analyzer—a beautiful terminal wizard that introduces developers to Skene's AI-powered growth analysis service. Currently operating as a single-use conversion funnel with exceptional UX but zero viral mechanics, it's positioned perfectly for a product-led transformation. The CLI's terminal aesthetics, local-first approach, and multi-provider flexibility create trust with technical audiences, but lack of analytics, sharing features, and repeat-use patterns leave massive growth opportunities untapped.

## GROWTH LOOPS IDENTIFIED

### Loop 1: Developer Community Amplification Loop
**Priority: HIGH | Estimated Impact: 3-5x GitHub stars in 60 days, 40%+ referral traffic**

**Current State:** Beautiful UI generates organic screenshots/tweets but no structured sharing mechanism. Easter egg game is viral-ready but invisible (must wait through installation).

**Implementation:**
- Add `--demo` flag that launches directly into space shooter game (skips installation) - makes easter egg discoverable and shareable
- Inject share prompts at high-emotion moments: after game high score, after stunning analysis results render, when error recovery succeeds
- Built-in `skene share [analysis-id]` command generates shareable link to public results dashboard (anonymized, opt-in)
- Embed "Analyzed with Skene CLI" footer in growth-plan.md output with clickable GitHub link
- Add `--theme` flag (cyberpunk, retro, matrix, etc.) to encourage "rice" culture sharing in r/unixporn
- Track share-to-install conversion with UTM parameters in generated links (append ?ref=share-twitter to GitHub URLs)

**Why This Works:** Developers share tools that make them look skilled. Terminal aesthetics + game + beautiful output = triple shareability trigger. Recorded demos in `internal/tui/components/ascii_motion_placeholder.go` (4.9MB!) suggest video content potential.

**Activation Metric:** % of completed analyses that generate shares | Target: 15% share rate

---

### Loop 2: Codebase Analysis Network Effect Loop
**Priority: HIGH | Estimated Impact: 2x MAU, 60% improved retention**

**Current State:** Tool is single-use—analyze once, never return. No memory, no comparison, no collaboration. `internal/services/analyzer/analyzer.go` has rich project detection but results vanish after viewing.

**Implementation:**
- Create `.skene/` directory in analyzed projects with persistent analysis history (JSON manifests of past runs)
- Add `skene diff` command comparing current codebase to previous analyses—"Your growth score improved from 73 to 81 since last month"
- Build `skene team` multiplayer mode: share project analysis with teammates, see who's viewed results, async comments on growth recommendations
- Introduce "Skene Score" (0-100) calculated from codebase health signals in `analyzer.go` (HasVenv, ConfigPath, PythonVersion, etc.) + AI analysis sentiment
- Generate dynamic score badge for README: `![Skene Score](https://skene.ai/badge/{project-id})` 
- Send weekly Slack/Discord digest: "3 teammates ran analyses, average score +5 this week" (requires opt-in webhook)
- Reference check in `internal/services/config/manager.go`: enable team workspace config storage

**Why This Works:** Transforms one-off tool into habit-forming dashboard. Competitive score + team visibility = weekly check-in behavior. Badges in READMEs = passive advertising to repo visitors.

**Activation Metric:** % of users who run 2nd analysis within 30 days | Target: 35% retention

---

### Loop 3: AI Provider Expansion → Content Marketing Loop
**Priority: MEDIUM | Estimated Impact: 10-15% top-of-funnel growth, SEO authority boost**

**Current State:** Multi-provider support (`internal/services/llm/client.go` supports 7+ providers) is competitive advantage but hidden value. No content leveraging this flexibility.

**Implementation:**
- Add `--benchmark` mode that runs same analysis across multiple AI providers, generates comparison report: "Claude found 12 growth opportunities, GPT-4 found 8, Gemini found 15"
- Auto-generate blog posts from benchmark results: "We analyzed 1,000 repos—here's which AI is best for SaaS vs dev tools"
- Create `skene/benchmarks` GitHub repo with weekly automated runs against popular open source projects, published as markdown tables
- Build provider-specific landing pages: "Best AI for analyzing Python projects" linking to CLI installation
- Add `--model-playground` flag for A/B testing prompts—"Test your custom system prompt against 5 models instantly"
- Submit benchmark data to AI leaderboards (Hugging Face, Papers with Code) with attribution links

**Why This Works:** Solves real decision paralysis ("Which AI should I use?") while generating SEO content. Each benchmark = 7 provider mentions = 7 backlink opportunities. Positions Skene as neutral evaluator (builds trust despite Skene provider preference).

**Activation Metric:** Benchmark report shares + organic search traffic | Target: 500 monthly benchmark runs, 20% traffic from SEO

---

### Loop 4: Magic Link Auth → Data Flywheel Loop
**Priority: MEDIUM | Estimated Impact: 25%+ conversion to Skene accounts, powers AI improvement**

**Current State:** `internal/tui/views/auth.go` implements magic link but no post-auth hooks. No analytics pipeline. Skene provider positioned first in list but no data capture justifying that preference.

**Implementation:**
- Add privacy-respecting telemetry (opt-in during welcome screen): anonymized flow completion %, error rates, provider choices, project types
- Inject A/B testing framework in `internal/services/growth/engine.go`—randomly vary wizard order, button copy, provider positioning
- Build internal dashboard showing conversion funnels: "73% complete system check, 45% choose Skene provider, 89% of Skene users complete analysis vs 62% of OpenAI users"
- Use data to train better defaults: auto-select provider based on detected project type (Python → recommend Claude, Go → recommend GPT-4)
- Generate monthly "State of Developer Tools" reports from aggregate analysis data (anonymized): "42% of analyzed projects lack virtual environments"
- Feed growth plan quality ratings back to AI prompts—"Users rate Claude growth plans 4.2/5, GPT-4 3.8/5" → adjust system prompts

**Why This Works:** Currently flying blind (no analytics). Data enables iterative optimization + creates unique content assets. Privacy-first approach maintains developer trust while capturing critical insights.

**Activation Metric:** Opt-in telemetry rate + iteration velocity | Target: 60% opt-in, 2-week experiment cycles

---

### Loop 5: Results Dashboard → Implementation Services Loop  
**Priority: LOW | Estimated Impact: 15-20% revenue conversion (assuming paid services exist)**

**Current State:** `internal/tui/views/results.go` renders growth plans but ends with vague "Next Steps"—no clear path from analysis to implementation. Beautiful output but unclear monetization bridge.

**Implementation:**
- Add "Get Help Implementing" CTA in results dashboard linking to Skene consultation booking
- Generate implementation difficulty scores in growth-plan.md: "🟢 Easy: 8 items (2-week effort) | 🟡 Medium: 5 items (1-month) | 🔴 Hard: 3 items (partner needed)"
- Introduce `skene track [recommendation-id]` command for marking growth items as complete—builds project timeline
- After 3 tracked completions, prompt: "You've implemented 3 growth loops! Want help with the complex ones?" → lead form
- Create "Skene Certified Partner" directory—CLI shows local agencies who can implement growth plans
- Build self-serve marketplace: "Hire a PLG expert for $X to implement viral loop #2" (25% take rate)

**Why This Works:** Bridges gap from insights to execution. CLI generates demand for services; services validate CLI's value. Implementation tracking = engagement hook + upsell trigger.

**Activation Metric:** % of analyses that lead to paid engagement | Target: 8-12% conversion

---

## IMPLEMENTATION ROADMAP

**Week 1-2: Instrumentation & Quick Wins**
- Implement opt-in telemetry system with privacy controls
- Add `--demo` flag for instant game access (viral Easter egg unlock)
- Inject "Share your results" prompt after analysis completion
- Ship "Analyzed with Skene" footer in growth-plan.md outputs
- Deploy basic funnel analytics (system check → installation → analysis completion)

**Week 3-4: Retention Foundation**
- Build `.skene/` persistent history storage in analyzed projects
- Create `skene diff` command for longitudinal analysis comparison
- Implement Skene Score calculation (0-100) from analyzer.go signals
- Add dynamic badge generation endpoint for README embeds
- A/B test provider ordering in wizard (Skene first vs alphabetical)

**Week 5-6: Network Effects**
- Launch `skene team` multiplayer analysis sharing
- Build webhook integrations (Slack, Discord) for team digests
- Create public analysis dashboard with shareable links (opt-in)
- Add `--theme` flag with 3 visual variants (cyberpunk, retro, matrix)
- Implement UTM tracking for share attribution

**Week 7-8: Content & SEO Flywheel**
- Ship `--benchmark` mode for multi-provider comparison
- Set up automated weekly benchmark runs against top 100 GitHub repos
- Build provider-specific landing pages with SEO optimization
- Submit first benchmark report to AI leaderboards
- Create `skene/benchmarks` public repo with historical data

---

## SUCCESS METRICS

**Activation (30-day targets):**
- Install-to-complete-analysis rate: **65%** (currently ~45% estimated from code flow)
- Second analysis within 30 days: **35%** (currently ~5% single-use)
- Telemetry opt-in rate: **60%** (new metric)

**Viral Growth:**
- Share rate (analyses → shares): **15%** (currently 0% structured)
- Share-to-install conversion: **8%** (new funnel)
- GitHub stars growth: **+300/month** (currently organic only)

**Engagement:**
- Weekly active analysts: **2x** (from repeat-use features)
- Skene Score checks per user: **4x/month** (new behavior)
- Team workspaces created: **25** (new feature)

**Monetization (if services exist):**
- Free-to-paid conversion: **10%** (new path)
- Average implementation deal size: **$5,000** (assumption)

**Strategic:**
- SEO traffic from benchmarks: **500 visits/month** (currently 0)
- Organic backlinks: **+50/quarter** (from reports, badges)
- Average session duration: **+40%** (from game, diff features)

---

**Critical Dependencies:**
1. Skene backend must support analysis ID tracking for sharing/history
2. Badge generation requires simple image rendering service
3. Team features need lightweight auth/storage layer (Firebase/Supabase)
4. Benchmark automation needs CI/CD pipeline (GitHub Actions)

**Risk Mitigation:**
- All loops designed privacy-first (opt-in, anonymized)
- Features degrade gracefully if backend unavailable (local-only mode)
- No forced account creation—preserve low-friction UVX path
- Open source remains fully functional without telemetry

This strategy transforms Skene CLI from a beautiful one-time installer into a habit-forming growth platform that developers return to weekly, share proactively, and eventually pay to accelerate. The terminal UI remains the moat—competitors can copy features, but not the craft in `internal/tui/styles/styles.go`.