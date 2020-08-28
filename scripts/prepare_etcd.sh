#!/bin/bash

etcdctl put current.float 1.23
etcdctl put current.bool true
etcdctl put current.int 11
etcdctl put current.string hello
etcdctl put current.serialize.name test
etcdctl put sub.sub.key1 sub.sub.key