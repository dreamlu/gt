- 我需要做什么  
守护进程  

- 如何做  
1.chan/goroutine  
2.设计守护进程模式  

- 我想要的守护进程  
1.指定任意时刻执行程序  
2.通过任务队列来执行  
3.新task来临时,队列更新(指针即可),总体维持固定携程数量进行执行  