"""
Schema definition for homerun2-core-catcher deployment configuration
"""

schema CoreCatcher:
    """Configuration for homerun2-core-catcher Kubernetes deployment"""

    # Metadata
    name: str = "homerun2-core-catcher"
    namespace: str = "homerun2"

    # Image
    image: str = "ghcr.io/stuttgart-things/homerun2-core-catcher:latest"
    imagePullPolicy: str = "IfNotPresent"

    # Replicas and resources
    replicas: int = 1
    cpuRequest: str = "100m"
    cpuLimit: str = "500m"
    memoryRequest: str = "128Mi"
    memoryLimit: str = "512Mi"

    # Redis
    redisAddr: str = "redis-stack.homerun2.svc.cluster.local"
    redisPort: int | str = "6379"
    redisStream: str = "messages"

    # Consumer group
    consumerGroup: str = "homerun2-core-catcher"

    # Secrets
    redisPassword: str = ""

    # Extra environment variables
    extraEnvVars: {str:str} = {}

    # Labels and annotations
    labels: {str:str} = {}
    annotations: {str:str} = {}

    check:
        replicas >= 0, "replicas must be >= 0"
        name != "", "name must not be empty"
        namespace != "", "namespace must not be empty"
