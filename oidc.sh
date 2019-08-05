### dex.coreos.com ###
ETCDCTL_API=3 etcdctl get /registry/dex.coreos.com --prefix=true --cacert=/etc/kubernetes/pki/etcd/ca.crt --cert=/etc/kubernetes/pki/apiserver-etcd-client.crt --key=/etc/kubernetes/pki/apiserver-etcd-client.key

### CRD ###
ETCDCTL_API=3 etcdctl get /registry/apiextensions.k8s.io/customresourcedefinitions/authcodes.dex.coreos.com --cacert=/etc/kubernetes/pki/etcd/ca.crt --cert=/etc/kubernetes/pki/apiserver-etcd-client.crt --key=/etc/kubernetes/pki/apiserver-etcd-client.key

ETCDCTL_API=3 etcdctl get /registry/apiextensions.k8s.io/customresourcedefinitions/authrequests.dex.coreos.com --cacert=/etc/kubernetes/pki/etcd/ca.crt --cert=/etc/kubernetes/pki/apiserver-etcd-client.crt --key=/etc/kubernetes/pki/apiserver-etcd-client.key

ETCDCTL_API=3 etcdctl get /registry/apiextensions.k8s.io/customresourcedefinitions/connectors.dex.coreos.com --cacert=/etc/kubernetes/pki/etcd/ca.crt --cert=/etc/kubernetes/pki/apiserver-etcd-client.crt --key=/etc/kubernetes/pki/apiserver-etcd-client.key

ETCDCTL_API=3 etcdctl get /registry/apiextensions.k8s.io/customresourcedefinitions/oauth2clients.dex.coreos.com --cacert=/etc/kubernetes/pki/etcd/ca.crt --cert=/etc/kubernetes/pki/apiserver-etcd-client.crt --key=/etc/kubernetes/pki/apiserver-etcd-client.key

ETCDCTL_API=3 etcdctl get /registry/apiextensions.k8s.io/customresourcedefinitions/offlinesessionses.dex.coreos.com --cacert=/etc/kubernetes/pki/etcd/ca.crt --cert=/etc/kubernetes/pki/apiserver-etcd-client.crt --key=/etc/kubernetes/pki/apiserver-etcd-client.key

ETCDCTL_API=3 etcdctl get /registry/apiextensions.k8s.io/customresourcedefinitions/passwords.dex.coreos.com --cacert=/etc/kubernetes/pki/etcd/ca.crt --cert=/etc/kubernetes/pki/apiserver-etcd-client.crt --key=/etc/kubernetes/pki/apiserver-etcd-client.key

ETCDCTL_API=3 etcdctl get /registry/apiextensions.k8s.io/customresourcedefinitions/refreshtokens.dex.coreos.com --cacert=/etc/kubernetes/pki/etcd/ca.crt --cert=/etc/kubernetes/pki/apiserver-etcd-client.crt --key=/etc/kubernetes/pki/apiserver-etcd-client.key

ETCDCTL_API=3 etcdctl get /registry/apiextensions.k8s.io/customresourcedefinitions/signingkeies.dex.coreos.com --cacert=/etc/kubernetes/pki/etcd/ca.crt --cert=/etc/kubernetes/pki/apiserver-etcd-client.crt --key=/etc/kubernetes/pki/apiserver-etcd-client.key
