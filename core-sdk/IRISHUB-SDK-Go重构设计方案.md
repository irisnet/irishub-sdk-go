
#irishub-sdk-go 重构设计方案
* 现有sdk的整理设计

##重构目标
- 将现有sdk的功能模块条理化、清晰化；
- 将冗余的接口整理清晰；
- 设计新的测试模块；
- 增加部分代码注释；

##具体模块的重构

###一、 sdk-core-go 
该sdk-core包含以下核心功能
- auth 模块；
- bank 模块；
  单位转换
- keys 模块；
- client.go

###二、将以下模块单个独立出来
sdk-module-go

- coinswap 模块；
  - --go.mod
    
- htlc 模块；
    - --go.mod
    
- gov 模块；
    - --go.mod
    
- nft 模块；
    - --go.mod
    
- oracle 模块；
    - --go.mod
    
- random 模块；
    - --go.mod
    
- record 模块；
    - --go.mod
    
- service 模块；
    - --go.mod
    
- staking 模块；
    - --go.mod
    
- token 模块；
    - --go.mod


###三、设计新的测试模块；
- sdk重构后，需要进一步测试，新的测试模块用于重构后的sdk测试

###四、对于代码中必要的注释可以适当的增加

###五、sdk的使用说明文档



###六、 core-sdk重构设计具体项
该sdk-core包含以下核心功能
- auth 模块；
- bank 模块；
  单位转换
- keys 模块；
  - keys 需要接口设计，支持SM2（用于联盟链 irita）
  - 目前keys主要是入口，实际逻辑在crypto里面
- client.go





