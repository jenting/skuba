#!/bin/bash

ETCDCTL_API=3 etcdctl get "" --prefix=true --keys-only=true --cacert=/etc/kubernetes/pki/etcd/ca.crt --cert=/etc/kubernetes/pki/apiserver-etcd-client.crt --key=/etc/kubernetes/pki/apiserver-etcd-client.key

ETCDCTL_API=3 etcdctl get "" --prefix=true --keys-only=true --cacert=/etc/kubernetes/pki/etcd/ca.crt --cert=/etc/kubernetes/pki/cilium-etcd-client.crt --key=/etc/kubernetes/pki/cilium-etcd-client.key

ETCDCTL_API=3 etcdctl get "" --prefix=true --keys-only=true --cacert=/etc/kubernetes/pki/etcd/ca.crt --cert=/etc/kubernetes/pki/fake-etcd-client.crt --key=/etc/kubernetes/pki/fake-etcd-client.key
