# minix文件系统
基于golang和fuse实现用户态文件系统，用于学习文件系统的设计原理。
***把一个概念重新实现一次，是最快的掌握这个概念方法，也是掌握这个概念最深刻的方式。***

## 特性
1. 使用golang实现
2. 基于fuse实现用户态文件系统

## 安装
### 环境要求
centos 8.5最小化安装, 下载地址： [centos8.5](https://mirrors.aliyun.com/centos/8.5.2111/isos/x86_64/CentOS-8.5.2111-x86_64-dvd1.iso)。 选择最小化安装 
```shell
[root@MiWiFi-RA72-srv minix]# uname -r
4.18.0-348.el8.x86_64
[root@MiWiFi-RA72-srv minix]# cat /etc/system-release
CentOS Linux release 8.5.2111
```

### 安装步骤

#### 编译并安装minix驱动
```shell
cd /etc/yum.repos.d/
sed -i 's/mirrorlist/#mirrorlist/g' /etc/yum.repos.d/CentOS-*
sed -i 's|#baseurl=http://mirror.centos.org|baseurl=http://vault.centos.org|g' /etc/yum.repos.d/CentOS-*
yum install -y vim
wget https://vault.centos.org/8.5.2111/BaseOS/Source/SPackages/kernel-4.18.0-348.el8.src.rpm
rpm -ivh kernel-4.18.0-348.el8.src.rpm 
cp ./rpmbuild/SOURCES/linux-4.18.0-348.el8.tar.xz .
xz -d ./linux-4.18.0-348.el8.tar.xz 
tar -xvf ./linux-4.18.0-348.el8.tar 
cd linux-4.18.0-348.el8/fs/minix/
yum install -y kernel-devel
cp -rf * /usr/src/kernels/4.18.0-348.7.1.el8_5.x86_64/fs/minix/
yum install gcc -y
yum install make -y
yum install -y elfutils-libelf-devel
make CONFIG_MINIX_FS=m -C /usr/src/kernels/4.18.0-348.7.1.el8_5.x86_64 M=/usr/src/kernels/4.18.0-348.7.1.el8_5.x86_64/fs/minix
cd /usr/src/kernels/4.18.0-348.7.1.el8_5.x86_64/fs/minix
insmod ./minix.ko
mkfs.minix /dev/sdb
mount /dev/sdb /mnt/
```

#### 安装golang 1.18版本
```shell
wget https://go.dev/dl/go1.18.10.linux-amd64.tar.gz
tar -zxvf ./go1.18.10.linux-amd64.tar.gz
echo "export GO111MODULE=on" >> ~/.bash_profile
echo "export GOPROXY=https://goproxy.cn" >> ~/.bash_profile 
echo "export GOPATH=/root/gopath/" >> ~/.bash_profile 
echo "export GOROOT=/root/go" >> ~/.bash_profile 
echo "export PATH=$PATH:$GOROOT/bin:$GOPATH/bin" >> ~/.bash_profile  
source ~/.bash_profile 
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

