# Rule R276: Runbook Requirement

## Rule Statement
Every Software Factory project MUST include a comprehensive RUNBOOK.md that enables operators to deploy, run, troubleshoot, and maintain the software in production.

## Criticality Level
**MANDATORY** - Required for operational readiness
Violation = -25% grade penalty

## Core Principle
**"If operators can't run it, it's not production-ready"**

## Required Runbook Sections

### Minimum Required Template

```markdown
# RUNBOOK - [Project Name]

## Overview
Brief description of what this software does and its purpose.

## Prerequisites
### System Requirements
- OS: Linux/MacOS/Windows
- Runtime: Go 1.21+ / Python 3.9+ / Node 18+
- Memory: Minimum 2GB RAM
- Disk: 500MB free space
- Network: Outbound HTTPS (443)

### Required Tools
- [ ] git
- [ ] docker (if containerized)
- [ ] kubectl (if Kubernetes)
- [ ] make

### Access Requirements
- [ ] Repository access
- [ ] Container registry credentials (if applicable)
- [ ] Cloud provider credentials (if applicable)

## Installation

### From Source
\`\`\`bash
git clone <repository>
cd <project>
make build
make install
\`\`\`

### From Binary
\`\`\`bash
curl -L https://releases.../latest/binary -o binary
chmod +x binary
sudo mv binary /usr/local/bin/
\`\`\`

### From Container
\`\`\`bash
docker pull registry/image:latest
docker run registry/image:latest
\`\`\`

## Configuration

### Environment Variables
\`\`\`bash
export APP_PORT=8080
export APP_ENV=production
export APP_LOG_LEVEL=info
export DATABASE_URL=postgres://...
\`\`\`

### Configuration File
\`\`\`yaml
# config.yaml
server:
  port: 8080
  host: 0.0.0.0
database:
  url: postgres://...
\`\`\`

### Secrets Management
- Store secrets in environment variables or secret manager
- Never commit secrets to repository
- Rotate secrets regularly

## Deployment

### Local Development
\`\`\`bash
make run-local
# or
docker-compose up
\`\`\`

### Production Deployment
\`\`\`bash
# Kubernetes
kubectl apply -f manifests/
kubectl rollout status deployment/app

# Docker
docker run -d --name app -p 8080:8080 registry/image:latest

# Systemd
sudo systemctl start app
sudo systemctl enable app
\`\`\`

## Operations

### Starting the Service
\`\`\`bash
systemctl start app
# or
docker start app
# or
kubectl scale deployment/app --replicas=3
\`\`\`

### Stopping the Service
\`\`\`bash
systemctl stop app
# or
docker stop app
# or
kubectl scale deployment/app --replicas=0
\`\`\`

### Restarting the Service
\`\`\`bash
systemctl restart app
# or
docker restart app
# or
kubectl rollout restart deployment/app
\`\`\`

### Health Checks
\`\`\`bash
curl http://localhost:8080/health
# Expected response: {"status":"healthy"}
\`\`\`

## Monitoring

### Key Metrics
- CPU usage: <80%
- Memory usage: <1GB
- Response time: <200ms
- Error rate: <1%

### Log Locations
- Application logs: /var/log/app/app.log
- Error logs: /var/log/app/error.log
- Access logs: /var/log/app/access.log

### Log Queries
\`\`\`bash
# View recent logs
tail -f /var/log/app/app.log

# Search for errors
grep ERROR /var/log/app/app.log

# Docker logs
docker logs -f app

# Kubernetes logs
kubectl logs -f deployment/app
\`\`\`

## Troubleshooting

### Service Won't Start
1. Check port availability: `lsof -i :8080`
2. Verify configuration: `app validate-config`
3. Check permissions: `ls -la /var/lib/app`
4. Review logs: `tail -100 /var/log/app/error.log`

### High Memory Usage
1. Check for memory leaks: `pprof`
2. Review cache settings
3. Increase memory limits
4. Enable memory profiling

### Slow Performance
1. Check database queries
2. Review API response times
3. Check network latency
4. Enable performance profiling

### Common Error Messages
| Error | Cause | Solution |
|-------|-------|----------|
| "port already in use" | Another process on port | Change port or kill process |
| "database connection failed" | DB unreachable | Check DB URL and network |
| "permission denied" | Insufficient permissions | Check file ownership |

## Maintenance

### Backup Procedures
\`\`\`bash
# Database backup
pg_dump database > backup.sql

# File backup
tar -czf backup.tar.gz /var/lib/app
\`\`\`

### Update Procedures
\`\`\`bash
# Rolling update
kubectl set image deployment/app app=registry/image:v2

# Blue-green deployment
kubectl apply -f manifests/v2/
kubectl patch service app -p '{"spec":{"selector":{"version":"v2"}}}'
\`\`\`

### Rollback Procedures
\`\`\`bash
# Kubernetes rollback
kubectl rollout undo deployment/app

# Docker rollback
docker stop app
docker run -d --name app registry/image:previous

# Git rollback
git revert HEAD
make build && make deploy
\`\`\`

## Security

### Security Checklist
- [ ] TLS enabled for all endpoints
- [ ] Authentication required
- [ ] Rate limiting configured
- [ ] Secrets properly managed
- [ ] Regular security updates

### Incident Response
1. Isolate affected systems
2. Preserve logs for analysis
3. Apply security patches
4. Document incident
5. Post-mortem review

## Disaster Recovery

### RTO/RPO Targets
- Recovery Time Objective: 1 hour
- Recovery Point Objective: 1 hour

### Backup Strategy
- Daily automated backups
- Weekly full backups
- Monthly archives
- Offsite backup storage

### Recovery Procedures
1. Provision new infrastructure
2. Restore from latest backup
3. Verify data integrity
4. Test application functionality
5. Switch traffic to recovered system

## Support

### Contact Information
- Team: Platform Engineering
- Email: platform@company.com
- Slack: #platform-support
- On-call: PagerDuty

### Escalation Path
1. L1: Application logs and restart
2. L2: Configuration and deployment issues
3. L3: Code changes required
4. L4: Architecture changes needed

### SLA
- Response time: 15 minutes
- Resolution time: 4 hours
- Availability target: 99.9%

## Appendix

### Architecture Diagram
[Include system architecture diagram]

### API Documentation
[Link to API docs or swagger]

### Related Documentation
- [Design Document](docs/design.md)
- [API Reference](docs/api.md)
- [Development Guide](docs/development.md)
```

## Validation Script

```bash
#!/bin/bash
# validate-runbook.sh

echo "🔍 Validating RUNBOOK.md..."

if [ ! -f "RUNBOOK.md" ]; then
    echo "❌ RUNBOOK.md not found!"
    exit 1
fi

MISSING_SECTIONS=0

check_section() {
    local section="$1"
    if ! grep -q "$section" RUNBOOK.md; then
        echo "❌ Missing section: $section"
        ((MISSING_SECTIONS++))
    else
        echo "✅ Found section: $section"
    fi
}

# Check required sections
check_section "## Prerequisites"
check_section "## Installation"
check_section "## Configuration"
check_section "## Deployment"
check_section "## Operations"
check_section "## Monitoring"
check_section "## Troubleshooting"
check_section "## Maintenance"

if [ $MISSING_SECTIONS -gt 0 ]; then
    echo "❌ Runbook incomplete: $MISSING_SECTIONS sections missing"
    exit 1
else
    echo "✅ Runbook validation passed"
    exit 0
fi
```

## Runbook Quality Criteria

### Excellent Runbook (A Grade)
- All required sections present
- Step-by-step procedures
- Troubleshooting decision trees
- Performance tuning guide
- Disaster recovery tested

### Good Runbook (B Grade)
- Most sections complete
- Clear deployment steps
- Basic troubleshooting
- Monitoring explained

### Minimal Runbook (C Grade)
- Installation steps
- Configuration basics
- How to start/stop
- Where to find logs

### Failed Runbook (F Grade)
- Missing or empty
- No deployment instructions
- No troubleshooting help
- Operators cannot run software

## Integration with Other Rules

### Prerequisites
- Software must be deployable (R275)
- Production ready (R274)

### Enables
- External user validation (R278)
- PR plan execution (R279)

## Grading Impact

- No RUNBOOK.md: -25%
- Missing critical sections: -5% each
- No troubleshooting guide: -10%
- No deployment instructions: -15%
- Poor quality/unclear: -10%

## Summary

R276 ensures every Software Factory project includes operational documentation that enables real-world deployment and maintenance. A good runbook is the difference between software that works in development and software that runs in production.