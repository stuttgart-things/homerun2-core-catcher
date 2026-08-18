"""
Schema definition for homerun2-core-catcher deployment configuration
"""

schema CoreCatcher:
    """Configuration for homerun2-core-catcher Kubernetes deployment"""

    # Metadata
    name: str = "homerun2-core-catcher"
    namespace: str = "homerun2"

    # Image
    # Placeholder for local renders only — the Release workflow overrides this
    # with the released version via kustomize-image-key. Points at :main rather
    # than :latest because :latest is not maintained and still resolves to
    # v0.2.0 from March.
    image: str = "ghcr.io/stuttgart-things/homerun2-core-catcher:main"
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

    # Service
    serviceType: str = "ClusterIP"
    servicePort: int = 80
    containerPort: int = 8080

    # HTTPRoute (Gateway API)
    httpRouteEnabled: bool = False
    httpRouteParentRefName: str = ""
    httpRouteParentRefNamespace: str = ""
    httpRouteHostname: str = ""
    httpRouteAnnotations: {str:str} = {}

    # Secrets
    redisPassword: str = ""

    # Catcher mode
    # Default "web" so a profile-less render still emits the Service +
    # HTTPRoute that downstream Flux/ArgoCD consumers expect. The Release
    # workflow always passes tests/kcl-deploy-profile.yaml (which also sets
    # web), so this only changes behavior when the profile is bypassed.
    catcherMode: str = "web"

    # Extra environment variables
    extraEnvVars: {str:str} = {}

    # Labels and annotations
    labels: {str:str} = {}
    annotations: {str:str} = {}

    check:
        replicas >= 0, "replicas must be >= 0"
        servicePort >= 1 and servicePort <= 65535, "servicePort must be valid"
        containerPort >= 1 and containerPort <= 65535, "containerPort must be valid"
        name != "", "name must not be empty"
        namespace != "", "namespace must not be empty"
