## Endpoint Discover

### usage
```yaml
      - name: create endpoint discover
        uses: lwnmengjing/endpoint-discover@v0.0.1
        with:
          namespace: 'default'
          configmap-name: 'endpoint-discover'
          protocols: 'grpc,http'
          config-name: 'endpoints.yaml'
```