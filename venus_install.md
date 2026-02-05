# Venus deploy


## 部署类型 

- 高可用部署：  
  - 管理节点：  
    1. 数量：3 个
    2. 配置：>=8 核 16G内存，磁盘越大越好

  - 工作节点：  
    1. 数量：N个（根据实际应用确认）
    2. 配置：>=8 核 16G内存

- 单机部署：  
  - 管理节点：  
    1. 数量：1 个
    2. 配置：>=16 核 32G内存

## 安装目标  

单节点部署、安装master节点

## 环境准备  

虚拟机一台，操作系统为：ubuntu 20.04 x86_64  
配置：10核32G内存  
安装包：http://10.0.3.23:49173/Venus/releases/product/master/venus-amd64-master.tar.gz  

前置条件：  

- 需要开启root登录：  

```sh

# 修改 /etc/ssh/sshd_config
PermitRootLogin yes

```

- 重启ssh服务：  

```sh
sudo systemctl restart ssh
```
