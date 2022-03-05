## Endpoint Discover

### usage
```yaml
      - name: create endpoint discover
        uses: lwnmengjing/endpoint-discover@v0.0.2
        with:
          cluster-url: ${{ steps.kubeconfig.outputs.cluster_url }}
          token: ${{ steps.kubeconfig.outputs.token }}
          configmap-name: endpoint-discover
          namespace: beta
          protocols: 'grpc,http'
          config-name: 'endpoints.yml'
```