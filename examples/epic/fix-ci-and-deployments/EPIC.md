# Epic: Fix CI and Deployments

## Executive Summary

This epic implements a Docker-based local CI simulation system to reduce GitHub Actions budget burn by 70%+ while providing developers with sub-30-second feedback on CI failures. The solution combines intelligent change detection, exact environment replication, and automated testing to catch 95% of CI failures locally before pushing to GitHub, dramatically improving developer velocity and reducing infrastructure costs.

## Problem Statement

### Current Pain Points

**Budget & Resource Waste:**
- Excessive GitHub Actions usage burning through budget due to environment drift failures
- Failed CI runs consuming resources without delivering value
- Estimated $200-400/month in unnecessary Actions costs

**Developer Experience Issues:**
- CI failures only discovered after push, causing workflow interruption and context switching
- Slow feedback loop: waiting 3-5 minutes for CI results delays development velocity
- Developers losing confidence in CI reliability due to inconsistent failures

**Environment Inconsistencies:**
- Different Node versions between local (varies) and CI (18.17.0) environments
- Missing system dependencies causing build failures in CI but not locally
- Build cache differences leading to inconsistent test results
- Resource constraint differences masking performance issues

## Proposed Solution Overview

### Core Architecture

**Docker-First Infrastructure:** Complete environment parity using Docker Compose with exact CI resource constraints (2 CPU, 4GB RAM) and Node 18.17.0 to eliminate environment drift entirely.

**Intelligent Change Detection:** Smart file-based change detection to run only affected tests:
- `public-website/*` changes → Run public website tests only (15s)
- `admin-webtool/*` changes → Run admin tool tests only (15s)
- `shared/*` or `supabase/*` changes → Run both test suites (30s)
- Infrastructure changes → Full rebuild and test (60s)

**Multi-Layer Caching Strategy:** 
- Docker layer caching for dependencies
- Persistent volumes for `node_modules` and build cache
- Hash-based cache invalidation for dependency changes

**Git Hook Integration:** Pre-push hook with escape hatch (`--no-verify`) for emergencies, providing seamless developer workflow integration.

## Success Criteria (Measurable)

### Primary KPIs
- **95% CI failure detection rate locally** - Catch failures before GitHub push
- **<30 second feedback time** for typical single-application changes
- **80% reduction in GitHub Actions minutes** usage within 4 weeks
- **>85% developer adoption rate** after rollout

### Secondary Metrics
- **Cold start time <75 seconds** for complete environment setup
- **Warm cache time <20 seconds** for subsequent runs
- **Zero false positives** - if local passes, CI should pass
- **$200-400/month cost savings** in GitHub Actions usage

## Implementation Phases

### Phase 1: Core Infrastructure (Week 1)
**Deliverables:**
- `docker-compose.ci.yml` with exact resource constraints matching GitHub Actions
- Multi-stage Dockerfile for environment replication (Node 18.17.0, system packages)
- `scripts/ci-local.sh` orchestration script with basic change detection
- Root-level npm scripts (`pnpm ci:local`, `pnpm ci:setup`) for developer interface

**Acceptance Criteria:**
- Docker containers start successfully with resource limits
- Basic test execution works in containerized environment
- Manual script execution provides clear success/failure feedback

### Phase 2: Smart Detection & Testing (Week 2)
**Deliverables:**
- Intelligent file change detection algorithm implementation
- Database isolation with per-run PostgreSQL schema creation
- Vitest configuration matching CI exactly (threads, coverage, timeouts)
- Pre-push git hook integration with Husky

**Acceptance Criteria:**
- Change detection correctly identifies affected applications
- Tests run with identical configuration to GitHub Actions
- Pre-push hook blocks pushes on test failure with clear messaging

### Phase 3: Optimization & Monitoring (Week 3)
**Deliverables:**
- Fine-tuned Docker layer caching strategy for maximum performance
- Performance monitoring with detailed timing reports
- Flaky test detection and automatic retry logic
- Comprehensive developer documentation and troubleshooting guide

**Acceptance Criteria:**
- Cached runs complete in <30 seconds for single-app changes
- Performance reports show timing breakdown by phase
- Documentation enables self-service developer onboarding

### Phase 4: Validation & Rollout (Week 4)
**Deliverables:**
- Correlation analysis proving 95%+ local/CI result matching
- Performance validation meeting <30 second targets
- Team training materials and adoption tracking system
- GitHub Actions usage monitoring dashboard

**Acceptance Criteria:**
- Statistical validation of CI prediction accuracy
- All performance targets met consistently
- Team adoption >85% with positive feedback
- Measurable reduction in GitHub Actions usage

## Risk Assessment

### High Risk
**Docker Resource Requirements**
- *Risk:* Developer machines lack sufficient Docker memory allocation (need 8GB+)
- *Mitigation:* Pre-implementation environment audit, clear system requirements
- *Contingency:* Fallback to lighter containers or selective feature disabling

**Developer Adoption Resistance** 
- *Risk:* Developers bypass pre-push hooks due to perceived slowness
- *Mitigation:* Ensure <30s performance target met, provide clear escape hatch
- *Contingency:* Make hooks advisory-only initially, track voluntary usage

### Medium Risk
**Cache Invalidation Complexity**
- *Risk:* Incorrect cache invalidation leads to false positives/negatives
- *Mitigation:* Conservative invalidation strategy, extensive testing
- *Contingency:* Manual cache clearing commands, gradual cache sophistication

**Environment Parity Drift**
- *Risk:* GitHub Actions environment changes break local simulation
- *Mitigation:* Automated environment monitoring, CI/CD pipeline notifications
- *Contingency:* Quick rollback to previous known-good configuration

### Low Risk
**Database Test Isolation**
- *Risk:* Test contamination despite schema isolation
- *Mitigation:* Unique schema naming with UUID, cleanup verification
- *Contingency:* Fall back to full database reset between test runs

## Resource Requirements

### Development Time
- **Lead Developer:** 32 hours over 4 weeks (8h/week)
- **DevOps Support:** 8 hours for infrastructure review and validation
- **Team Testing:** 16 hours total (2h per developer) for validation and feedback

### Infrastructure
- **Docker Desktop:** All developer machines with 8GB+ RAM allocation
- **Disk Space:** ~2GB per developer for Docker images and cache
- **Network:** Stable internet for initial Docker image pulls

### Ongoing Maintenance
- **Monthly:** 2 hours for environment updates and performance monitoring
- **Quarterly:** 4 hours for GitHub Actions parity validation and optimization
- **Annual:** 8 hours for major version upgrades (Node.js, dependencies)

### Budget Impact
- **Investment:** ~80 development hours ($8,000 at $100/hour)
- **Monthly Savings:** $200-400 in GitHub Actions costs
- **ROI Timeline:** 2-3 weeks break-even, 12-month ROI >300%

---

**Epic Owner:** DevOps Lead  
**Stakeholders:** Development Team, Technical Leadership, Budget Management  
**Dependencies:** Docker Desktop rollout, pnpm workspace migration completion  
**Success Review:** Weekly progress reviews, final retrospective after Phase 4