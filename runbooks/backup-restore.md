# Backup and Restore Procedures

## PostgreSQL Backup

```bash
# Create backup
kubectl exec -it postgres-pod -- pg_dump -U postgres shopmicro > backup.sql

# Restore backup
kubectl exec -it postgres-pod -- psql -U postgres shopmicro < backup.sql
```

## Redis Backup

```bash
# Create backup
kubectl exec -it redis-pod -- redis-cli SAVE

# Copy backup file
kubectl cp redis-pod:/data/dump.rdb ./redis-backup.rdb

# Restore
kubectl cp ./redis-backup.rdb redis-pod:/data/dump.rdb
kubectl exec -it redis-pod -- redis-cli SHUTDOWN NOSAVE
kubectl delete pod redis-pod  # Let it restart
```

## Automated Backup (CronJob)

```yaml
apiVersion: batch/v1
kind: CronJob
metadata:
  name: postgres-backup
spec:
  schedule: "0 2 * * *" # Daily at 2 AM
  jobTemplate:
    spec:
      template:
        spec:
          containers:
            - name: backup
              image: postgres:15
              command:
                ["pg_dump", "-U", "postgres", "-h", "postgres", "shopmicro"]
              env:
                - name: PGPASSWORD
                  valueFrom:
                    secretKeyRef:
                      name: postgres-secret
                      key: password
          restartPolicy: OnFailure
```
