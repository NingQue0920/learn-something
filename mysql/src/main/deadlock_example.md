### 更复杂的死锁场景

创建基本表结构，id是主键，其他字段均无索引。
```sql
CREATE TABLE `t_student` (
  `id` int NOT NULL,
  `no` varchar(255) DEFAULT NULL,
  `name` varchar(255) DEFAULT NULL,
  `age` int DEFAULT NULL,
  `score` int DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```
准备初始数据
```sql
INSERT INTO `t_student` VALUES (1,'001','张三',18,100),
                               (5,'002','李四',19,99),
                               (10,'003','王五',20,98),
                               (15,'004','赵六',21,97),
                               (20,'005','田七',22,96);
```
开启两个事务，模拟死锁

| 事务1                                                    | 事务2 |
|--------------------------------------------------------| --- |
| BEGIN                                                  | BEGIN |
| update t_student set score = 100 where id = 16         |  |
|                                                        | update t_student set score = 100 where id = 17 |
| insert into t_student values(16 , ‘016’,‘宁缺1’ ,21,99 ) |  |
|                                                        | insert into t_student values(17 , ‘017’,‘宁缺2’ ,21,99 ) |

逐个SQL分析，为什么会死锁：
1. id是主键，是唯一索引，所以不可能是Next-Key Lock。
2. 表中没有id=16的数据，事务1的第一条SQL，没有数据命中，无法退化为Record Lock，所以只能是Gap Lock。
3. 同理，事务2的第一条SQL也是Gap Lock。
4. Gap Lock 之间不冲突。
   > Gap locks in `InnoDB` are “purely inhibitive”, which means that their only purpose is to prevent other transactions from inserting to the gap. **Gap locks can co-exist. A gap lock taken by one transaction does not prevent another transaction from taking a gap lock on the same gap.** There is no difference between shared and exclusive gap locks. They do not conflict with each other, and they perform the same function.

5. 事务1的 insert 语句会产生插入意向锁（Insert Intention Lock），区间是(15,20)，该插入意向锁与事务2的update语句添加的Gap Lock互斥，因此被阻塞。
6. 同理，事务2的insert语句同样会与事务1的Gap Lock阻塞；导致死锁。
