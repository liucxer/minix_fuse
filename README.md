# minix文件系统
基于golang和fuse实现用户态文件系统，用于学习文件系统的设计原理。
***把一个概念重新实现一次，是最快的掌握这个概念方法，也是掌握这个概念最深刻的方式。***

## 特性
1. 使用golang实现
2. 基于fuse实现用户态文件系统

## 安装
### 环境要求
centos 7.9最小化安装, 下载地址： [centos7.9](https://mirrors.aliyun.com/centos/7.9.2009/isos/x86_64/CentOS-7-x86_64-Minimal-2009.iso)。 
```shell
[root@liucx-centos79 ~]# cat /etc/system-release
CentOS Linux release 7.9.2009 (Core)
[root@liucx-centos79 ~]# uname -r
3.10.0-1160.el7.x86_64
```

安装步骤
```shell
yum install -y gcc
yum install -y kernel-devel
```

```shell
[root@liucx-centos7 ~]# uname -r
4.14.179-1.el7.x86_64
[root@liucx-centos7 ~]# uname -a
Linux liucx-centos7.9.novalocal 4.14.179-1.el7.x86_64 #1 SMP Tue May 12 02:22:15 EDT 2020 x86_64 x86_64 x86_64 GNU/Linux

centos7.9 安装4.14.179内核

```

## 编译minix驱动
```shell
make CONFIG_MINIX_FS=m -C /usr/src/kernels/3.10.0-1160.88.1.el7.x86_64 M=/usr/src/kernels/3.10.0-1160.88.1.el7.x86_64/fs/minix
insmod ./minix.ko # 安装驱动
```