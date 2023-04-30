# minix_fuse

## 环境
```shell
[root@liucx-centos7 ~]# uname -r
4.14.179-1.el7.x86_64
[root@liucx-centos7 ~]# uname -a
Linux liucx-centos7.9.novalocal 4.14.179-1.el7.x86_64 #1 SMP Tue May 12 02:22:15 EDT 2020 x86_64 x86_64 x86_64 GNU/Linux

centos7.9 安装4.14.179内核

```

## 编译minix驱动
```shell
make CONFIG_MINIX_FS=m -C /usr/src/kernels/4.14.179-1.el7.x86_64 M=/usr/src/kernels/4.14.179-1.el7.x86_64/fs/minix
insmod ./minix.ko # 安装驱动
```