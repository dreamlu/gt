TODO  
1.~~gt.Table()默认结构体驼峰小写值~~  
2.~~解析好的sql放入内存,加快解析效率~~  
3.~~map[string][]string转结构体(使用别名类型CMap)~~  
4.~~模型解析增加连接自身功能User->User(弃:业务极少见)~~  
5.~~添加gt:like支持(弃,用key替代)~~  
6.~~表关联通过其他字段~~  
7.~~crud.select()执行后sql clone问题~~  
8.~~subWhereSQL()没where不执行和where and多余问题等~~  
9.~~多表连接GetReflectTagMore()主题解析table[0]问题(非问题)~~  
10.~~crud事务支持~~  
11.~~驼峰解析替代json字段(弃)~~  
12.~~result返回动态添加参数~~  
13.~~多表连接多数据库初步支持~~  
14.~~优化缓存使用,提供默认set()值~~  
15.~~result add()支持结构体~~  
16.~~result mapdata和map整合(弃)~~  
17.~~批量创建支持指针数组~~  
18.~~使用CMap替代map[string][]string~~  
19.~~update()id默认不为空限制(弃)~~  
20.~~多文件上传,返回文件名数组~~(参考demo:[deercoder-gin](https://github.com/dreamlu/deercoder-gin))  
21.~~多表连接支持设置不同表名~~  
22.~~更进一步解析sql放入内存,加快解析效率~~  
23.~~字符串解析buf与移除goto,panic~~  
24.~~createMore移除id验证(可通过继承上级来去除id)~~  
25.~~事务中select()问题~~  
26.~~完善使用文档~~  
27.~~通过其他字段删除,结合select()(弃)~~  
28.~~支持mock假数据~~  
29.~~subWhereSQL同一个变量clone问题(弃)~~  
30.~~仅仅左连接LeftTable() bug~~  
31.文档介绍关于crud中常用方法关联Model()/clone等信息  
32.优化key搜索(更快搜索)  
33.~~日志大统一~~  
34.~~增加数据不存在详细内容(无法~)~~  
35.打印错误所在行数  
36.~~Table()多表创建修改支持(弃)~~  
37.~~config get内存缓存(无需)~~  
38.~~文件上传自动创建目录~~  
39.~~CDate映射sql date类型问题(无问题)~~  
40.~~crud总是复用导致的问题~~ 
41.~~在线查看日志不显示问题~~  
42.~~原生sql支持cmap参数和#{param}(or select(Struct))~~  
43.~~key,以及所有条件为空过滤~~    
44.~~日期默认值问题为空过滤(弃)~~  
45.~~keyModel不支持gt:field问题~~  
46.日志软链接docker问题  
47.~~暂无数据详细信息(弃)~~  
48.~~测试路径问题~~  
49.常用数据放缓存  
50.~~innertable()优化~~  
51.sql数据放入缓存后如何解决分页等问题  
52.~~create/update支持事务~~   
53.search搜索引擎  
54.~~集成消息中间件~~  
55.集成监控  
56.nsq挂载不包含message问题  
57.mock对于自定义类型的bug  
58.守护进程  
59.数据存储加密/解密  
60.~~跨表gt.table()问题~~  
61.错误处理统一封装，作为return返回同时还能记录日志  
62.~~事务,select后影响update问题~~  
n.xxx  