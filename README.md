# alertlogger

Provides a webhook for prometheus alertmanager and simple logs to stderr.



Install:
```
> kubectl apply -f https://raw.githubusercontent.com/mwennrich/alertlogger/master/samples/alertloggerStatefulSet.yaml
> kubectl apply -f https://raw.githubusercontent.com/mwennrich/alertlogger/master/samples/alertloggerService.yaml
```

Add to alertmanager config:
```
(...)
      routes:
        - match_re:
            severity: 'critical'
          receiver: 'alertlogger'
          continue: true
        - match_re:
            severity: 'warning'
          receiver: 'alertlogger'
          continue: true
(...)
    receivers:
    - name: 'null'
    - name: 'alertlogger'
      webhook_configs:
        - url: http://alertlogger:5001
```

Watch logs:
``` 
> kubectl logs -n monitoring alertlogger-0 -f
``` 

alerts can now be ingested to your preferred log-shipper solution like fluentd or promtail/loki.
