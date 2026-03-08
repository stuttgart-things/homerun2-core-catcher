"""
Main entry point for homerun2-core-catcher Kubernetes deployment
Imports all resources and exports manifests
"""

import namespace
import serviceaccount
import configmap
import secret
import deploy

# Export all manifests (filter out None values)
manifests = [
    m for m in [
        namespace.namespace
        serviceaccount.serviceAccount
        configmap.configMap
        secret.secretRedis
        deploy.deployment
    ] if m
]
