"""
Main entry point for homerun2-core-catcher Kubernetes deployment
Imports all resources and exports manifests
"""

import serviceaccount
import configmap
import secret
import deploy
import service
import httproute

# Export all manifests (filter out None values).
# Namespace is NOT emitted from the KCL profile — the consuming Application
# uses `CreateNamespace=true`, and PR-preview namespaces are pre-created by
# Kyverno. Co-tenanting multiple Applications in one ns (the homerun2 stack
# pattern) requires that no single child claim ownership of the Namespace.
# See stuttgart-things/homerun2-omni-pitcher#132.
manifests = [
    m for m in [
        serviceaccount.serviceAccount
        configmap.configMap
        secret.secretRedis
        deploy.deployment
        service.service
        httproute.httproute
    ] if m
]
