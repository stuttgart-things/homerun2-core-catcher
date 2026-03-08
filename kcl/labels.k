"""
Common labels for homerun2-core-catcher resources
"""

import schema

# Get command-line options
_image = option("config.image")
_namespace = option("config.namespace")
_redisAddr = option("config.redisAddr")
_redisPort = option("config.redisPort")
_redisStream = option("config.redisStream")
_consumerGroup = option("config.consumerGroup")
_redisPassword = option("config.redisPassword")
_extraEnvVars = option("config.extraEnvVars") or {}

# Config instance with command-line overrides
config: schema.CoreCatcher = schema.CoreCatcher {
    if _image:
        image = _image
    if _namespace:
        namespace = _namespace
    if _redisAddr:
        redisAddr = _redisAddr
    if _redisPort:
        redisPort = _redisPort
    if _redisStream:
        redisStream = _redisStream
    if _consumerGroup:
        consumerGroup = _consumerGroup
    if _redisPassword:
        redisPassword = _redisPassword
    if _extraEnvVars:
        extraEnvVars = _extraEnvVars
}

# Common labels applied to all resources
commonLabels = {
    "app.kubernetes.io/name": config.name
    "app.kubernetes.io/instance": config.name
    "app.kubernetes.io/component": "consumer"
    "app.kubernetes.io/part-of": "homerun2"
    **config.labels
}

# Selector labels for pods
selectorLabels = {
    "app.kubernetes.io/name": config.name
    "app.kubernetes.io/instance": config.name
}
