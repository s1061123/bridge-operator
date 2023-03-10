module github.com/s1061123/bridge-operator

go 1.18

require (
	github.com/spf13/cobra v1.6.1
	github.com/spf13/pflag v1.0.5
	github.com/vishvananda/netlink v1.1.1-0.20220112194529-e5fd1f8193de
	k8s.io/api v0.22.7
	k8s.io/apimachinery v0.22.7
	k8s.io/client-go v1.5.2
	k8s.io/code-generator v0.22.7
	k8s.io/klog v1.0.0
	k8s.io/kube-openapi v0.0.0-20211109043538-20434351676c
	k8s.io/kubernetes v1.26.0
)

require (
	github.com/PuerkitoBio/purell v1.1.1 // indirect
	github.com/PuerkitoBio/urlesc v0.0.0-20170810143723-de5bf2ad4578 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/emicklei/go-restful v2.9.5+incompatible // indirect
	github.com/evanphx/json-patch v4.11.0+incompatible // indirect
	github.com/go-logr/logr v0.4.0 // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/jsonreference v0.19.5 // indirect
	github.com/go-openapi/swag v0.19.14 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-cmp v0.5.5 // indirect
	github.com/google/gofuzz v1.1.0 // indirect
	github.com/googleapis/gnostic v0.5.5 // indirect
	github.com/imdario/mergo v0.3.5 // indirect
	github.com/inconshreveable/mousetrap v1.0.1 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.11 // indirect
	github.com/mailru/easyjson v0.7.6 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/vishvananda/netns v0.0.0-20200728191858-db3c7e526aae // indirect
	golang.org/x/mod v0.4.2 // indirect
	golang.org/x/net v0.0.0-20211209124913-491a49abca63 // indirect
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d // indirect
	golang.org/x/sys v0.0.0-20210616094352-59db8d763f22 // indirect
	golang.org/x/term v0.0.0-20210220032956-6a3ed077a48d // indirect
	golang.org/x/text v0.3.6 // indirect
	golang.org/x/time v0.0.0-20210723032227-1f47c861a9ac // indirect
	golang.org/x/tools v0.1.2 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	google.golang.org/appengine v1.6.5 // indirect
	google.golang.org/protobuf v1.26.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/apiserver v0.22.7 // indirect
	k8s.io/component-base v0.22.7 // indirect
	k8s.io/gengo v0.0.0-20201214224949-b6c5ce23f027 // indirect
	k8s.io/klog/v2 v2.9.0 // indirect
	k8s.io/utils v0.0.0-20211116205334-6203023598ed // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.2.1 // indirect
	sigs.k8s.io/yaml v1.2.0 // indirect
)

replace (
	k8s.io/api => k8s.io/api v0.22.7
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.22.7
	k8s.io/apimachinery => k8s.io/apimachinery v0.22.7
	k8s.io/apiserver => k8s.io/apiserver v0.22.7
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.22.7
	k8s.io/client-go => k8s.io/client-go v0.22.7
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.22.7
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.22.7
	k8s.io/code-generator => k8s.io/code-generator v0.22.7
	k8s.io/component-base => k8s.io/component-base v0.22.7
	k8s.io/component-helpers => k8s.io/component-helpers v0.22.7
	k8s.io/controller-manager => k8s.io/controller-manager v0.22.7
	k8s.io/cri-api => k8s.io/cri-api v0.22.7
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.22.7
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.22.7
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.22.7
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.22.7
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.22.7
	k8s.io/kubectl => k8s.io/kubectl v0.22.7
	k8s.io/kubelet => k8s.io/kubelet v0.22.7
	k8s.io/kubernetes => k8s.io/kubernetes v1.22.7
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.22.7
	k8s.io/metrics => k8s.io/metrics v0.22.7
	k8s.io/mount-utils => k8s.io/mount-utils v0.22.7
	k8s.io/pod-security-admission => k8s.io/pod-security-admission v0.22.7
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.22.7
)
