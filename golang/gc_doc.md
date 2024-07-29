Go的一个GC周期分为4个阶段：
1. **Sweep Termination（清除终止）**：对应的函数是`runtime.gcWaitOnMark`，该函数会等待上一轮的GC工作完成。
   全局变量`gcphase`的值不为`_GCmark`，表示GC工作已经完成。
    ```
   n = work.cycles.Load()
   // 自旋
   for{
        nMarks = work.cycles.Load()
        if gcphase != _GCmark {
            nMarks++
        }
        if nMarks > n {
            return
        }
   }
   ```
2. **Mark（标记）**：对应函数`runtime.gcMarkRootPrepare`(标记Root对象，STW),`runtime.gcMarkTinyAllocs`,`runtime.gcBgMarkStartWorkers`（启动后台标记Goroutine）,`runtime.gcMarkWorker / runtime.gcDrain`（标记工作，与工作Goroutine并发）。
3. **MarkTermination（标记终止）**：对应函数`runtime.gcMarkTermination`(STW)，修改GC状态为_GCTermination，关闭GC线程和协助线程。
4. **Sweep（清除）**：对应函数`runtime.gcSweep`，清除不再使用的对象，释放内存。