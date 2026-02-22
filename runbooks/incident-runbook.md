# Incident Response Runbook: ShopMicro Backend Outage

## Overview

This runbook covers response to a simulated backend service outage in the ShopMicro platform.

## Detection

- Alert: Backend health check fails
- Monitoring: Grafana dashboard shows backend pod restarts or high error rates
- Logs: Loki shows repeated connection failures

## Initial Assessment (5 minutes)

1. Check Grafana dashboard for backend metrics
2. Review recent deployments: `kubectl rollout history deployment/backend`
3. Check pod status: `kubectl get pods -l app=backend`
4. Review logs: `kubectl logs -l app=backend --tail=100`

## Containment (10 minutes)

1. If pods are crashing, scale down to 0: `kubectl scale deployment backend --replicas=0`
2. Check database connectivity: `kubectl exec -it postgres-pod -- psql -U postgres -d shopmicro -c "SELECT 1"`
3. Check Redis connectivity: `kubectl exec -it redis-pod -- redis-cli ping`

## Recovery (15 minutes)

1. Rollback deployment: `kubectl rollout undo deployment/backend`
2. Scale back up: `kubectl scale deployment backend --replicas=2`
3. Monitor recovery: `kubectl rollout status deployment/backend`

## Post-Incident (30 minutes)

1. Root cause analysis in Grafana Tempo traces
2. Update monitoring alerts if needed
3. Document findings in incident report

## Prevention

- Implement circuit breakers
- Add more comprehensive health checks
- Regular dependency updates
