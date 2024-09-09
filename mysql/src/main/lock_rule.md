### 加锁的规则
### 通用规则

- **锁定读（SELECT … FOR UPDATE , SELECT … LOCK IN SHARE MODE）, DELETE , UPDATE 语句会在命中的的每个索引记录上加锁，这些锁通常是Next-Key Lock。**
- **如果在搜索的过程中使用到了二级索引，且二级索引上要设置排他锁，则InnoDB会检索相应的聚簇索引，并在其上加锁。**
- **如果SQL语句没有触发索引，则MySQL需要扫描整个表，此时表中的每一行，以及行与行之间的间隙都会被锁定。即锁全表。**

### 具体细节

- 普通的`SELECT`语句不会加锁（快照读），除非在`串行化(Serializable)`隔离级别下，才会对涉及的所有行加`共享锁`
- `Locking Read` ，`UPDATE`，`DELETE` 语句，加锁规则取决于使用的索引类型：
    - 唯一索引：仅锁定命中的记录，不锁间隙。
    - 无索引/非唯一索引：锁定扫描的索引范围，可能加`Gap Lock` 或 `Next-Key Lock`。
- `SELECT … FOR UPDATE` 与 `SELECT … FOR SHARE` 是冲突的。
- 一致性读（Consistent Read 快照读）是特殊的，它不会被锁定，因为它基于事务开始时的快照，所以不会受到其他事务中锁的干扰。
- UPDATE修改聚簇索引时，会对受到影响的二级索引数据加隐式锁。
- `INSERT`语句会在插入的行上添加`Record Lock`，而不是Next-Key Lock。如果发生duplicate-key error，则会在该索引记录上添加`共享锁`。

### 锁的退化
一句话概括：如果仅使用Record Lock 或 Gap Lock 就能解决幻读，则Next-Key Lock会退化。

详细一点：
1. 索引上的等值查询，给唯一索引加锁时，如果命中，Next-Key Lock 退化为Record Lock。
2. 索引上的等值查询，向右遍历且最后一个值不满足等值条件时，Next-Key Lock 退化为Gap Lock。