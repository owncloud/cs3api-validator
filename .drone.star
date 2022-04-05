config = {
    'name': 'cs3api-validator',
    'rocketchat': {
        'channel': 'builds',
        'from_secret': 'private_rocketchat'
    },
    'branches': [
        'main'
    ],
}

def main(ctx):
    before = beforePipelines(ctx)
    if not before:
        print('Errors detected. Review messages above.')
        return []
    stages = stagePipelines(ctx)
    if not stages:
        print('Errors detected. Review messages above.')
        return []
    dependsOn(before, stages)
    after = afterPipelines(ctx)
    dependsOn(stages, after)
    return before + stages + after

def beforePipelines(ctx):
    return linting(ctx)

def stagePipelines(ctx):
    testPipelines = tests(ctx)
    dockerReleasePipelines = dockerRelease(ctx, "amd64")
    dependsOn(testPipelines, dockerReleasePipelines)
    dockerAfterRelease = releaseDockerReadme(ctx) + releaseDockerManifest(ctx)
    dependsOn(dockerReleasePipelines, dockerAfterRelease)
    return testPipelines + dockerReleasePipelines + dockerAfterRelease

def afterPipelines(ctx):
    return [
        notify()
    ]

def dependsOn(earlierStages, nextStages):
    for earlierStage in earlierStages:
        for nextStage in nextStages:
            nextStage['depends_on'].append(earlierStage['name'])

def notify():
    result = {
        'kind': 'pipeline',
        'type': 'docker',
        'name': 'chat-notifications',
        'clone': {
            'disable': True
        },
        'steps': [
            {
                'name': 'notify-rocketchat',
                'image': 'plugins/slack:1',
                'pull': 'always',
                'settings': {
                    'webhook': {
                        'from_secret': config['rocketchat']['from_secret']
                    },
                    'channel': config['rocketchat']['channel']
                }
            }
        ],
        'depends_on': [],
        'trigger': {
            'ref': [
                'refs/tags/**'
            ],
            'status': [
                'success',
                'failure'
            ]
        }
    }

    for branch in config['branches']:
        result['trigger']['ref'].append('refs/heads/%s' % branch)

    return result

def linting(ctx):
    pipelines = []

    result = {
        'kind': 'pipeline',
        'type': 'docker',
        'name': 'lint',
        'steps': [
            {
                "name": "validate-go",
                "image": "golangci/golangci-lint:latest",
                "commands": [
                    "golangci-lint run -v",
                ]
            },
        ],
        'depends_on': [],
        'trigger': {
            'ref': [
                'refs/pull/**',
                'refs/tags/**'
            ]
        }
    }

    for branch in config['branches']:
        result['trigger']['ref'].append('refs/heads/%s' % branch)

    pipelines.append(result)

    return pipelines

def tests(ctx):
    pipelines = []
    result = {
        'kind': 'pipeline',
        'type': 'docker',
        'name': 'test-acceptance-cs3api',
        'steps': [
            {
                "name": "wait-for-ocis",
                "image": "owncloudci/wait-for:latest",
                "commands": [
                    "wait-for -it ocis:9200 -t 300",
                ],
            },
            {
                "name": "test",
                "image": "owncloudci/golang:1.17",
                "commands": [
                    "go test --endpoint=ocis:9142 -v",
                ],
                "volumes": [
                    {
                        "name": "gopath",
                        "path": "/go",
                    }
                ],
            },
        ],
        'services': ocisService(),
        'depends_on': [],
        'trigger': {
            'ref': [
                'refs/tags/**',
                'refs/pull/**',
            ]
        }
    }

    for branch in config['branches']:
        result['trigger']['ref'].append('refs/heads/%s' % branch)

    pipelines.append(result)

    return pipelines

def ocisService():
    return [{
        "name": "ocis",
        "image": "owncloud/ocis:latest",
        "pull": "always",
        "detach": True,
        "environment": {
            "OCIS_URL": "https://ocis:9200",
            "STORAGE_HOME_DRIVER": "ocis",
            "STORAGE_USERS_DRIVER": "ocis",
            "OCIS_LOG_LEVEL": "error",
            "STORAGE_GATEWAY_GRPC_ADDR": "0.0.0.0:9142"
        },
        "commands": [
            "ocis server",
        ],
    }]

def dockerRelease(ctx, arch):
    pipelines = []
    build_args = [
        "REVISION=%s" % (ctx.build.commit),
        "VERSION=%s" % (ctx.build.ref.replace("refs/tags/", "") if ctx.build.event == "tag" else "latest"),
    ]

    result = {
        "kind": "pipeline",
        "type": "docker",
        "name": "docker-%s" % (arch),
        "platform": {
            "os": "linux",
            "arch": arch,
        },
        "steps": [
            {
                "name": "build",
                "image": "owncloudci/golang:1.17",
                "commands": [
                    "go test -c -o cs3api-validator-linux-%s.test" % (arch),
                ],
            },
            {
                "name": "dryrun",
                "image": "plugins/docker:latest",
                "settings": {
                    "dry_run": True,
                    "tags": "linux-%s" % (arch),
                    "dockerfile": "docker/Dockerfile.linux.%s" % (arch),
                    "repo": ctx.repo.slug,
                    "build_args": build_args,
                },
                "when": {
                    "ref": {
                        "include": [
                            "refs/pull/**",
                        ],
                    },
                },
            },
            {
                "name": "docker",
                "image": "plugins/docker:latest",
                "settings": {
                    "username": {
                        "from_secret": "docker_username",
                    },
                    "password": {
                        "from_secret": "docker_password",
                    },
                    "auto_tag": True,
                    "auto_tag_suffix": "linux-%s" % (arch),
                    "dockerfile": "docker/Dockerfile.linux.%s" % (arch),
                    "repo": ctx.repo.slug,
                    "build_args": build_args,
                },
                "when": {
                    "ref": {
                        "exclude": [
                            "refs/pull/**",
                        ],
                    },
                },
            },
        ],
        "depends_on": [],
        "trigger": {
            "ref": [
                "refs/tags/v*",
                "refs/pull/**",
            ],
        },
        "volumes": [
            {
                "name": "gopath",
                "temp": {},
            }
        ],
    }

    for branch in config['branches']:
        result['trigger']['ref'].append('refs/heads/%s' % branch)

    pipelines.append(result)
    return pipelines

def releaseDockerManifest(ctx):
    pipelines = []
    result = {
        "kind": "pipeline",
        "type": "docker",
        "name": "manifest",
        "platform": {
            "os": "linux",
            "arch": "amd64",
        },
        "steps": [
            {
                "name": "execute",
                "image": "plugins/manifest:1",
                "settings": {
                    "username": {
                        "from_secret": "docker_username",
                    },
                    "password": {
                        "from_secret": "docker_password",
                    },
                    "spec": "docker/manifest.tmpl",
                    "auto_tag": True,
                    "ignore_missing": True,
                },
            },
        ],
        "depends_on": [],
        "trigger": {
            "ref": [
                "refs/tags/v*",
            ],
        },
    }
    for branch in config['branches']:
        result['trigger']['ref'].append('refs/heads/%s' % branch)

    pipelines.append(result)
    return pipelines

def releaseDockerReadme(ctx):
    pipelines = []
    result = {
        "kind": "pipeline",
        "type": "docker",
        "name": "readme",
        "platform": {
            "os": "linux",
            "arch": "amd64",
        },
        "steps": [
            {
                "name": "execute",
                "image": "chko/docker-pushrm:1",
                "environment": {
                    "DOCKER_USER": {
                        "from_secret": "docker_username",
                    },
                    "DOCKER_PASS": {
                        "from_secret": "docker_password",
                    },
                    "PUSHRM_TARGET": "owncloud/${DRONE_REPO_NAME}",
                    "PUSHRM_SHORT": "Docker images for %s" % (ctx.repo.name),
                    "PUSHRM_FILE": "README.md",
                },
            },
        ],
        "depends_on": [],
        "trigger": {
            "ref": [
                "refs/tags/v*",
            ],
        },
    }
    for branch in config['branches']:
        result['trigger']['ref'].append('refs/heads/%s' % branch)

    pipelines.append(result)
    return pipelines

def getPipelineNames(pipelines = []):
    """getPipelineNames returns names of pipelines as a string array

    Args:
      pipelines: array of drone pipelines

    Returns:
      names of the given pipelines as string array
    """
    names = []
    for pipeline in pipelines:
        names.append(pipeline["name"])
    return names
