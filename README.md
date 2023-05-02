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

切换4.14.179内核
```

### 安装步骤

#### 编译并安装minix驱动
```shell
make CONFIG_MINIX_FS=m -C /usr/src/kernels/3.10.0-1160.88.1.el7.x86_64 M=/usr/src/kernels/3.10.0-1160.88.1.el7.x86_64/fs/minix
insmod ./minix.ko # 安装驱动
```

```shell
yum install -y gcc
yum install -y wget
yum install -y vim 
yum install -y fuse
```

#### 安装golang 1.18版本
```shell
wget https://go.dev/dl/go1.18.10.linux-amd64.tar.gz
[root@linux-centos79 ~]# cat ~/.bash_profile 
# .bash_profile

# Get the aliases and functions
if [ -f ~/.bashrc ]; then
        . ~/.bashrc
fi

# User specific environment and startup programs

PATH=$PATH:$HOME/bin

export PATH
export GO111MODULE=on
export GOPROXY=https://goproxy.cn
export GOPATH=/root/gopath/
export GOROOT=/root/go
export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
```

#### 编译
```shell
liuchangxi@5c1bf473e937 minix_fuse % CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build .
liuchangxi@5c1bf473e937 minix_fuse % ls
README.md       fuse            go.mod          go.sum          main.go         minix           minix_decoder   minix_fuse
```

#### 运行
```shell
dd if=/dev/zero of=/dev/vdb bs=1M count=1024
mkfs.minix /dev/vdb
mount /dev/vdb /mnt
touch /mnt/111
echo "222" >> /mnt/111
umount /mnt/
./minix_fuse /mnt /dev/vdb
```

