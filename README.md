# 说明

LDR；

如果你的应用会涉及较多的文件读写，可以将 k8s 内存水位告警指标由 container_memory_working_set_bytes 改为 container_memory_rss。这样可以防止 page cache 占用空闲内存带来的误报警

## 场景描述

k8s 线上集群触发内存水位报警，报警显示某个 Pod 的内存使用率高达 85%。然而登陆到 pod 上发现应用实际占用的内存占用只有 50%，但是用 free 命令看到的内存占用又是符合报警水位（total 94G，free 40G）

```
$ free -h
              total        used        free      shared  buff/cache   available
Mem:            94G         31G         60G         25M        23G         40G
Swap:            0B          0B          0B
```

再细看一下 free 打印出的信息发现 buff/cache 占用有 23G，这部分内存用到哪里了？40G 的 available 是不是指这些 cache 中有 10G 是可以释放的？那不能释放的 23G 是用来做什么的？应用真实内存占用应该看哪个指标？接下来我们依次回答这些问题

## 缓存清理

Linux中的buff/cache可以被手动释放，释放缓存的代码如下：

```
echo 1 > /proc/sys/vm/drop_caches;
echo 2 > /proc/sys/vm/drop_caches;
echo 3 > /proc/sys/vm/drop_caches;
```

## 项目意义

此项目可以用于无损定时清理 buff/cache

