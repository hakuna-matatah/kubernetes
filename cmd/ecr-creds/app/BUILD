load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["ecrcredsprovider.go"],
    importpath = "k8s.io/kubernetes/cmd/ecr-creds/app",
    visibility = ["//visibility:public"],
    deps = [
        "//pkg/credentialprovider/aws:go_default_library",
        "//pkg/credentialprovider/external:go_default_library",
        "//staging/src/k8s.io/component-base/version/verflag:go_default_library",
        "//vendor/github.com/spf13/cobra:go_default_library",
        "//vendor/k8s.io/klog:go_default_library",
    ],
)

filegroup(
    name = "package-srcs",
    srcs = glob(["**"]),
    tags = ["automanaged"],
    visibility = ["//visibility:private"],
)

filegroup(
    name = "all-srcs",
    srcs = [":package-srcs"],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
